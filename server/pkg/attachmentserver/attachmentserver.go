package attachmentserver

import (
	"context"
	"log"
	"sync"

	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"
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
	DeleteChan chan primitive.ObjectID
}

type Uploaders struct {
	// Outer map key is UID, inner map key is MsgId
	data  map[primitive.ObjectID]map[primitive.ObjectID]Upload
	mutex sync.Mutex
}

type Upload struct {
	Index  uint16
	NextId primitive.ObjectID
}

type InChunk struct {
	Uid   primitive.ObjectID
	MsgId primitive.ObjectID
	Data  []byte
}

func Init(colls *db.Collections) *AttachmentServer {
	as := &AttachmentServer{
		Uploaders: Uploaders{
			data: make(map[primitive.ObjectID]map[primitive.ObjectID]Upload),
		},

		ChunkChan:  make(chan InChunk),
		DeleteChan: make(chan primitive.ObjectID),
	}
	runServer(as, colls)
	return as
}

func runServer(as *AttachmentServer, colls *db.Collections) {
	/* ------- Chunk loop ------- */
	go func() {
		for {
			chunk := <-as.ChunkChan
			as.Uploaders.mutex.Lock()
			nextId := primitive.NewObjectID()
			if _, ok := as.Uploaders.data[chunk.Uid]; !ok {
				// Create uploader data
				uploaderData := make(map[primitive.ObjectID]Upload)
				uploaderData[chunk.MsgId] = Upload{
					Index:  0,
					NextId: nextId,
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
				log.Println("Attachment chunk error:", err)
				as.Uploaders.mutex.Unlock()
				as.DeleteChan <- chunk.MsgId
				continue
			}
			if lastChunk {
				// Size less than 4mb, its the last chunk, upload is complete
				delete(as.Uploaders.data[chunk.Uid], chunk.MsgId)
				if len(as.Uploaders.data[chunk.Uid]) == 0 {
					delete(as.Uploaders.data, chunk.Uid)
				}
			} else {
				if _, ok := as.Uploaders.data[chunk.MsgId][chunk.MsgId]; ok {
					// Increment chunk index
					as.Uploaders.data[chunk.Uid][chunk.MsgId] = Upload{
						Index:  as.Uploaders.data[chunk.Uid][chunk.MsgId].Index + 1,
						NextId: nextId,
					}
				}
			}
			as.Uploaders.mutex.Unlock()
		}
	}()

	/* ------- Delete loop ------- */
	go func() {
		for {
			msgId := <-as.DeleteChan
			as.Uploaders.mutex.Lock()
			if _, err := colls.AttachmentMetadataCollection.DeleteOne(context.Background(), bson.M{"_id": msgId}); err != nil {
				log.Println("Error deleting attachment metadata:", err)
				continue
			}
			deleteAttachmentChunks(msgId, colls)
			as.Uploaders.mutex.Unlock()
		}
	}()
}

func deleteAttachmentChunks(chunkId primitive.ObjectID, colls *db.Collections) {
	chunkData := &models.AttachmentChunk{}
	if err := colls.AttachmentChunkCollection.FindOne(context.Background(), bson.M{"_id": chunkId}).Decode(&chunkData); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("Error finding attachment chunk:", err)
		}
		return
	}
	if _, err := colls.AttachmentChunkCollection.DeleteOne(context.Background(), bson.M{"_id": chunkId}); err != nil {
		log.Println("Error deleting attachment chunk:", err)
		return
	}
	if chunkData.NextChunkID == primitive.NilObjectID {
		return
	}
	deleteAttachmentChunks(chunkData.NextChunkID, colls)
}
