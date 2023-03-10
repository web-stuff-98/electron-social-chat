package attachmentserver

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
AttachmentServer. This is cleaner than the version in my last project.
Chunks are 4mb each.
*/

type AttachmentServer struct {
	Uploaders Uploaders

	ChunkChan  chan InChunk
	DeleteChan chan Delete
}

type Uploaders struct {
	// Outer map key is UID, inner map key is MsgId
	data  map[primitive.ObjectID]map[primitive.ObjectID]Upload
	mutex sync.Mutex
}

type Upload struct {
	Index      uint16
	TotalBytes uint32
	NextId     primitive.ObjectID
	LastChunk  time.Time
	MsgId      primitive.ObjectID
	// if timed out, the last chunk was received too long ago. upload has failed
	TimedOut bool
}

type InChunk struct {
	Uid           primitive.ObjectID
	MsgId         primitive.ObjectID
	SendUpdatesTo map[primitive.ObjectID]struct{}
	Data          []byte
	RecvChan      chan<- bool
}

type Delete struct {
	MsgId primitive.ObjectID
	Uid   primitive.ObjectID
}

func Init(ss *socketserver.SocketServer, colls *db.Collections) *AttachmentServer {
	as := &AttachmentServer{
		Uploaders: Uploaders{
			data: make(map[primitive.ObjectID]map[primitive.ObjectID]Upload),
		},

		ChunkChan:  make(chan InChunk),
		DeleteChan: make(chan Delete),
	}
	runServer(as, ss, colls)
	return as
}

func runServer(as *AttachmentServer, ss *socketserver.SocketServer, colls *db.Collections) {
	/* ------- Chunk loop ------- */
	go chunkLoop(as, ss, colls)
	/* ------- Delete loop ------- */
	go deleteLoop(as, ss, colls)
	/* ------- Watch for disconnects from the socketServer to clear uploaders & delete incomplete attachments ------- */
	go disconnectLoop(as, ss, colls)

	/* ------- Attachments fail when chunks haven't been received for a while. Keeps memory clear of stale uploads. ------- */
	cleanUpTicker := time.NewTicker(time.Second * 15)
	go func() {
		for {
			select {
			case <-cleanUpTicker.C:
				as.Uploaders.mutex.Lock()
				timedOut := make(map[primitive.ObjectID][]primitive.ObjectID)
				for uid, v := range as.Uploaders.data {
					for uploadId, u := range v {
						if u.LastChunk.Before(time.Now().Add(-time.Second * 15)) {
							as.Uploaders.data[uid][uploadId] = Upload{
								TimedOut:   true,
								Index:      u.Index,
								TotalBytes: u.TotalBytes,
								NextId:     u.NextId,
								LastChunk:  u.LastChunk,
							}
							timedOut[uid] = append(timedOut[uid], u.MsgId)
						}
					}
				}
				as.Uploaders.mutex.Unlock()
				for uid, uploads := range timedOut {
					for _, oi := range uploads {
						// don't use the delete channel because it also deletes the attachment metadata document
						// only the chunks and Upload struct should be removed
						deleteAttachmentChunks(oi, uid, oi, as, colls)
						colls.AttachmentMetadataCollection.UpdateByID(context.Background(), oi, bson.M{
							"$set": bson.M{
								"failed": true,
							},
						})
					}
				}
			}
		}
	}()
}

func chunkLoop(as *AttachmentServer, ss *socketserver.SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in attachment server chunk loop:", r)
			}
		}()
		chunk := <-as.ChunkChan
		as.Uploaders.mutex.Lock()
		metaData := &models.AttachmentData{}
		if err := colls.AttachmentMetadataCollection.FindOne(context.Background(), bson.M{"_id": chunk.MsgId}).Decode(&metaData); err != nil {
			as.Uploaders.mutex.Unlock()
			as.DeleteChan <- Delete{
				MsgId: chunk.MsgId,
				Uid:   chunk.Uid,
			}
			chunk.RecvChan <- false
			continue
		} else {
			if metaData.Failed {
				chunk.RecvChan <- false
				continue
			}
		}
		nextId := primitive.NewObjectID()
		if _, ok := as.Uploaders.data[chunk.Uid]; !ok {
			// Create uploader data
			uploaderData := make(map[primitive.ObjectID]Upload)
			uploaderData[chunk.MsgId] = Upload{
				Index:      0,
				NextId:     nextId,
				TotalBytes: uint32(metaData.Size),
				LastChunk:  time.Now(),
				MsgId:      chunk.MsgId,
			}
			as.Uploaders.data[chunk.Uid] = uploaderData
		}
		lastChunk := len(chunk.Data) < 4*1024*1024
		var chunkId primitive.ObjectID
		if lastChunk {
			nextId = primitive.NilObjectID
		}
		if as.Uploaders.data[chunk.Uid][chunk.MsgId].Index == 0 {
			chunkId = chunk.MsgId
		} else {
			chunkId = as.Uploaders.data[chunk.Uid][chunk.MsgId].NextId
		}
		// Write chunk
		if _, err := colls.AttachmentChunkCollection.InsertOne(context.Background(), models.AttachmentChunk{
			ID:          chunkId,
			Data:        primitive.Binary{Data: chunk.Data},
			NextChunkID: nextId,
		}); err != nil {
			as.Uploaders.mutex.Unlock()
			as.DeleteChan <- Delete{
				MsgId: chunk.MsgId,
				Uid:   chunk.Uid,
			}
			chunk.RecvChan <- false
			continue
		}
		if lastChunk {
			// Size less than 4mb, its the last chunk, upload is complete
			delete(as.Uploaders.data[chunk.Uid], chunk.MsgId)
			if len(as.Uploaders.data[chunk.Uid]) == 0 {
				delete(as.Uploaders.data, chunk.Uid)
			}
			// Send progress update
			colls.AttachmentMetadataCollection.UpdateByID(context.Background(), chunk.MsgId, bson.M{
				"$set": bson.M{
					"ratio": 1,
				},
			})
			ss.SendDataToUsers <- socketserver.UsersDataMessage{
				Uids: chunk.SendUpdatesTo,
				Data: socketmodels.AttachmentProgress{
					Ratio:  1,
					Failed: false,
					MsgID:  chunk.MsgId.Hex(),
				},
				Type: "ATTACHMENT_PROGRESS",
			}
		} else {
			if upload, ok := as.Uploaders.data[chunk.Uid][chunk.MsgId]; ok {
				if upload.TimedOut {
					chunk.RecvChan <- false
					as.Uploaders.mutex.Unlock()
					continue
				} else {
					// Send progress update
					ratio := (float32(upload.Index) * (4 * 1024 * 1024)) / float32(upload.TotalBytes)
					colls.AttachmentMetadataCollection.UpdateByID(context.Background(), chunk.MsgId, bson.M{
						"$set": bson.M{
							"ratio": ratio,
						},
					})
					ss.SendDataToUsers <- socketserver.UsersDataMessage{
						Uids: chunk.SendUpdatesTo,
						Data: socketmodels.AttachmentProgress{
							Ratio:  ratio,
							Failed: false,
							MsgID:  chunk.MsgId.Hex(),
						},
						Type: "ATTACHMENT_PROGRESS",
					}
					// Increment chunk index
					as.Uploaders.data[chunk.Uid][chunk.MsgId] = Upload{
						Index:      upload.Index + 1,
						TotalBytes: upload.TotalBytes,
						NextId:     nextId,
						LastChunk:  time.Now(),
						MsgId:      chunk.MsgId,
					}
				}
			}
		}
		chunk.RecvChan <- true
		as.Uploaders.mutex.Unlock()
	}
}

func deleteLoop(as *AttachmentServer, ss *socketserver.SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in attachment server delete loop:", r)
			}
		}()
		deleteData := <-as.DeleteChan
		as.Uploaders.mutex.Lock()
		if _, err := colls.AttachmentMetadataCollection.DeleteOne(context.Background(), bson.M{"_id": deleteData.MsgId}); err != nil {
			log.Println("Error deleting attachment metadata:", err)
			delete(as.Uploaders.data[deleteData.Uid], deleteData.MsgId)
			if len(as.Uploaders.data[deleteData.Uid]) == 0 {
				delete(as.Uploaders.data, deleteData.Uid)
			}
			continue
		}
		deleteAttachmentChunks(deleteData.MsgId, deleteData.Uid, deleteData.MsgId, as, colls)
		as.Uploaders.mutex.Unlock()
	}
}

func disconnectLoop(as *AttachmentServer, ss *socketserver.SocketServer, colls *db.Collections) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in attachment server disconnect loop:", r)
			}
		}()
		uid := <-ss.AttachmentServerRemoveUploaderChan
		as.Uploaders.mutex.Lock()
		for msgId := range as.Uploaders.data[uid] {
			deleteAttachmentChunks(msgId, uid, msgId, as, colls)
		}
		delete(as.Uploaders.data, uid)
		as.Uploaders.mutex.Unlock()
	}
}

func deleteAttachmentChunks(chunkId primitive.ObjectID, uid primitive.ObjectID, msgId primitive.ObjectID, as *AttachmentServer, colls *db.Collections) {
	chunkData := &models.AttachmentChunk{}
	if err := colls.AttachmentChunkCollection.FindOne(context.Background(), bson.M{"_id": chunkId}).Decode(&chunkData); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("Error finding attachment chunk:", err)
			delete(as.Uploaders.data[uid], msgId)
			if len(as.Uploaders.data[uid]) == 0 {
				delete(as.Uploaders.data, uid)
			}
		}
		return
	}
	if _, err := colls.AttachmentChunkCollection.DeleteOne(context.Background(), bson.M{"_id": chunkId}); err != nil {
		log.Println("Error deleting attachment chunk:", err)
		delete(as.Uploaders.data[uid], msgId)
		if len(as.Uploaders.data[uid]) == 0 {
			delete(as.Uploaders.data, uid)
		}
		return
	}
	if chunkData.NextChunkID == primitive.NilObjectID {
		return
	}
	deleteAttachmentChunks(chunkData.NextChunkID, uid, msgId, as, colls)
}
