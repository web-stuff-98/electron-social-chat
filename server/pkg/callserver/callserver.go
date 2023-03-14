package callserver

import (
	"log"
	"sync"

	"github.com/web-stuff-98/electron-social-chat/pkg/socketmodels"
	"github.com/web-stuff-98/electron-social-chat/pkg/socketserver"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	This is for 1-1 video calls
*/

type CallServer struct {
	// Calls that have not yet been answered
	CallsPending CallsPending
	// Channel for creating pending calls
	CallsPendingChan chan InCall
	// Channel for responding to pending calls
	ResponseToCallChan chan InCallResponse
	// Mutex protected map for active calls
	CallsActive CallsActive
	// Channel for closing active calls
	LeaveCallChan chan primitive.ObjectID
	// Channel for sending call recipient offer
	SendCallRecipientOffer chan CallerSignal
	// Channel for sending answer from called back to caller
	SendCalledAnswer chan CalledSignal
	// Channel for recipient requesting WebRTC re-initialization (necessary for changing/adding media devices)
	CallRecipientRequestedReInitialization chan primitive.ObjectID
}

/* --------------- MUTEX PROTECTED MAPS --------------- */
type CallsPending struct {
	// outer map is caller ID, inner map is the user that was called ID
	data  map[primitive.ObjectID]primitive.ObjectID
	mutex sync.Mutex
}
type CallsActive struct {
	// outer map is caller ID, inner map is the user that was called ID
	data  map[primitive.ObjectID]primitive.ObjectID
	mutex sync.Mutex
}

/* --------------- STRUCTS --------------- */

type CallerSignal struct {
	Caller primitive.ObjectID
	Signal string

	UserMediaStreamID string
	UserMediaVid      bool
	DisplayMediaVid   bool
}
type CalledSignal struct {
	Called primitive.ObjectID
	Signal string

	UserMediaStreamID string
	UserMediaVid      bool
	DisplayMediaVid   bool
}
type InCall struct {
	Caller primitive.ObjectID
	Called primitive.ObjectID
}
type InCallResponse struct {
	Caller primitive.ObjectID
	Called primitive.ObjectID
	Accept bool
}

func Init(ss *socketserver.SocketServer, dc chan primitive.ObjectID) *CallServer {
	cs := &CallServer{
		CallsPending: CallsPending{
			data: make(map[primitive.ObjectID]primitive.ObjectID),
		},
		CallsPendingChan:   make(chan InCall),
		ResponseToCallChan: make(chan InCallResponse),
		CallsActive: CallsActive{
			data: make(map[primitive.ObjectID]primitive.ObjectID),
		},
		LeaveCallChan:                          make(chan primitive.ObjectID),
		SendCallRecipientOffer:                 make(chan CallerSignal),
		SendCalledAnswer:                       make(chan CalledSignal),
		CallRecipientRequestedReInitialization: make(chan primitive.ObjectID),
	}
	runServer(ss, cs, dc)
	return cs
}

func runServer(ss *socketserver.SocketServer, cs *CallServer, dc chan primitive.ObjectID) {
	/* ----- Call pending loop ----- */
	go callPendingChanLoop(ss, cs)
	/* ----- Call response loop ----- */
	go callResponseChanLoop(ss, cs)
	/* ----- Leave call channel loop ----- */
	go leaveCallChanLoop(ss, cs)
	/* ----- Send call recipient webRTC offer loop ----- */
	go sendCallRecipientOfferLoop(ss, cs)
	/* ----- Send call recipient webRTC offer loop ----- */
	go sendCallerAnswerLoop(ss, cs)
	/* ----- Call recipient request webRTC reinitialization loop ----- */
	go callRecipientRequestReInitializationLoop(ss, cs)
	/* ----- Socket disconnect registration loop ----- */
	go socketDisconnectRegistrationLoop(ss, cs, dc)
}

func callPendingChanLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in pending call channel:", r)
			}
			go callPendingChanLoop(ss, cs)
		}()
		data := <-cs.CallsPendingChan
		cs.CallsPending.mutex.Lock()
		if called, ok := cs.CallsPending.data[data.Caller]; ok {
			if called != data.Called {
				// pending call switching to different user. cancel previous pending call.
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[called] = struct{}{}
				Uids[data.Caller] = struct{}{}
				ss.SendDataToUsers <- socketserver.UsersDataMessage{
					Uids: Uids,
					Type: "CALL_USER_RESPONSE",
					Data: socketmodels.CallResponse{
						Called: data.Called.Hex(),
						Caller: data.Caller.Hex(),
						Accept: false,
					},
				}
				cs.CallsPending.data[data.Caller] = data.Called
			}
		} else {
			cs.CallsPending.data[data.Caller] = data.Called
		}
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[data.Called] = struct{}{}
		Uids[data.Caller] = struct{}{}
		ss.SendDataToUsers <- socketserver.UsersDataMessage{
			Uids: Uids,
			Type: "CALL_USER_ACKNOWLEDGE",
			Data: socketmodels.CallAcknowledge{
				Caller: data.Caller.Hex(),
				Called: data.Called.Hex(),
			},
		}
		cs.CallsPending.mutex.Unlock()
	}
}

func callResponseChanLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in call response channel:", r)
			}
			go callResponseChanLoop(ss, cs)
		}()
		data := <-cs.ResponseToCallChan
		cs.CallsPending.mutex.Lock()
		cs.CallsActive.mutex.Lock()
		delete(cs.CallsPending.data, data.Caller)

		if data.Accept {
			// Close any call that either user is currently in.
			// Clients can only be in a single call.
			// Confusing variable names here.
			closedCallerCall := false
			closedCalledCall := false
			if callerCalled, ok := cs.CallsActive.data[data.Caller]; ok {
				closedCallerCall = true
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[data.Caller] = struct{}{}
				Uids[callerCalled] = struct{}{}
				ss.SendDataToUsers <- socketserver.UsersDataMessage{
					Type: "CALL_LEFT",
					Data: socketmodels.CallLeft{},
					Uids: Uids,
				}
				delete(cs.CallsActive.data, data.Caller)
			}
			if calledCalled, ok := cs.CallsActive.data[data.Called]; ok {
				closedCalledCall = true
				Uids := make(map[primitive.ObjectID]struct{})
				Uids[data.Called] = struct{}{}
				Uids[calledCalled] = struct{}{}
				ss.SendDataToUsers <- socketserver.UsersDataMessage{
					Type: "CALL_LEFT",
					Data: socketmodels.CallLeft{},
					Uids: Uids,
				}
				delete(cs.CallsActive.data, data.Called)
			}
			// make sure that the caller is not in a call. If they are exit the call they are already in
			if !closedCallerCall {
				for caller, called := range cs.CallsActive.data {
					if data.Caller == called {
						Uids := make(map[primitive.ObjectID]struct{})
						Uids[caller] = struct{}{}
						Uids[called] = struct{}{}
						ss.SendDataToUsers <- socketserver.UsersDataMessage{
							Type: "CALL_LEFT",
							Data: socketmodels.CallLeft{},
							Uids: Uids,
						}
						delete(cs.CallsActive.data, caller)
						break
					}
				}
			}
			// make sure that the called user is not in a call. If they are exit the call they are already in
			if !closedCalledCall {
				for caller, called := range cs.CallsActive.data {
					if data.Called == called {
						Uids := make(map[primitive.ObjectID]struct{})
						Uids[caller] = struct{}{}
						Uids[called] = struct{}{}
						ss.SendDataToUsers <- socketserver.UsersDataMessage{
							Type: "CALL_LEFT",
							Data: socketmodels.CallLeft{},
							Uids: Uids,
						}
						delete(cs.CallsActive.data, caller)
						break
					}
				}
			}

			// Any active calls that either user in have now been closed. Proceed.
			cs.CallsActive.data[data.Caller] = data.Called
		}

		// Send the response to both clients
		Uids := make(map[primitive.ObjectID]struct{})
		Uids[data.Called] = struct{}{}
		Uids[data.Caller] = struct{}{}
		ss.SendDataToUsers <- socketserver.UsersDataMessage{
			Uids: Uids,
			Type: "CALL_USER_RESPONSE",
			Data: socketmodels.CallResponse{
				Caller: data.Caller.Hex(),
				Called: data.Called.Hex(),
				Accept: data.Accept,
			},
		}

		cs.CallsPending.mutex.Unlock()
		cs.CallsActive.mutex.Unlock()
	}
}

func leaveCallChanLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in leave call channel:", r)
			}
			go leaveCallChanLoop(ss, cs)
		}()
		uid := <-cs.LeaveCallChan
		cs.CallsActive.mutex.Lock()
		if called, ok := cs.CallsActive.data[uid]; ok {
			ss.SendDataToUser <- socketserver.UserDataMessage{
				Type: "CALL_LEFT",
				Data: socketmodels.CallLeft{},
				Uid:  called,
			}
			delete(cs.CallsActive.data, uid)
		} else {
			for caller, called := range cs.CallsActive.data {
				if called == uid {
					ss.SendDataToUser <- socketserver.UserDataMessage{
						Type: "CALL_LEFT",
						Data: socketmodels.CallLeft{},
						Uid:  caller,
					}
					delete(cs.CallsActive.data, caller)
					break
				}
			}
		}
		cs.CallsActive.mutex.Unlock()
	}
}

func sendCallRecipientOfferLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in send call recipient offer loop :", r)
			}
			go sendCallRecipientOfferLoop(ss, cs)
		}()
		data := <-cs.SendCallRecipientOffer
		cs.CallsActive.mutex.Lock()
		if called, ok := cs.CallsActive.data[data.Caller]; ok {
			ss.SendDataToUser <- socketserver.UserDataMessage{
				Uid:  called,
				Type: "CALL_WEBRTC_OFFER_FROM_INITIATOR",
				Data: socketmodels.CallWebRTCOfferFromInitiator{
					Signal: data.Signal,

					UserMediaStreamID: data.UserMediaStreamID,
					UserMediaVid:      data.UserMediaVid,
					DisplayMediaVid:   data.DisplayMediaVid,
				},
			}
		}
		cs.CallsActive.mutex.Unlock()
	}
}

func sendCallerAnswerLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in sending offer answer loop :", r)
			}
			go sendCallerAnswerLoop(ss, cs)
		}()
		data := <-cs.SendCalledAnswer
		cs.CallsActive.mutex.Lock()
		for caller, oi2 := range cs.CallsActive.data {
			if oi2 == data.Called {
				ss.SendDataToUser <- socketserver.UserDataMessage{
					Uid:  caller,
					Type: "CALL_WEBRTC_ANSWER_FROM_RECIPIENT",
					Data: socketmodels.CallWebRTCOfferAnswer{
						Signal: data.Signal,

						UserMediaStreamID: data.UserMediaStreamID,
						UserMediaVid:      data.UserMediaVid,
						DisplayMediaVid:   data.DisplayMediaVid,
					},
				}
				break
			}
		}
		cs.CallsActive.mutex.Unlock()
	}
}

func callRecipientRequestReInitializationLoop(ss *socketserver.SocketServer, cs *CallServer) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in caller recipient request re-initialization loop :", r)
			}
			go callRecipientRequestReInitializationLoop(ss, cs)
		}()
		callerUid := <-cs.CallRecipientRequestedReInitialization
		cs.CallsActive.mutex.Lock()
		for caller, oi2 := range cs.CallsActive.data {
			if oi2 == callerUid {
				ss.SendDataToUser <- socketserver.UserDataMessage{
					Uid:  caller,
					Type: "CALL_WEBRTC_REQUESTED_REINITIALIZATION",
					Data: socketmodels.CallWebRTCRequestedReInitialization{},
				}
				break
			}
		}
		cs.CallsActive.mutex.Unlock()
	}
}

func socketDisconnectRegistrationLoop(ss *socketserver.SocketServer, cs *CallServer, dc chan primitive.ObjectID) {
	for {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("Recovered from panic in caller socket disconnect registration loop :", r)
			}
			go socketDisconnectRegistrationLoop(ss, cs, dc)
		}()
		uid := <-dc
		cs.CallsPending.mutex.Lock()
		if callPending, ok := cs.CallsPending.data[uid]; ok {
			ss.SendDataToUser <- socketserver.UserDataMessage{
				Uid:  callPending,
				Type: "CALL_USER_RESPONSE",
				Data: socketmodels.CallResponse{
					Caller: uid.Hex(),
					Called: callPending.Hex(),
					Accept: false,
				},
			}
			delete(cs.CallsActive.data, uid)
		}
		for caller, called := range cs.CallsPending.data {
			if called == uid {
				ss.SendDataToUser <- socketserver.UserDataMessage{
					Uid:  caller,
					Type: "CALL_USER_RESPONSE",
					Data: socketmodels.CallResponse{
						Caller: caller.Hex(),
						Called: uid.Hex(),
						Accept: false,
					},
				}
				delete(cs.CallsPending.data, caller)
			}
		}
		cs.CallsPending.mutex.Unlock()

		cs.CallsActive.mutex.Lock()
		if called, ok := cs.CallsActive.data[uid]; ok {
			ss.SendDataToUser <- socketserver.UserDataMessage{
				Uid:  called,
				Type: "CALL_LEFT",
				Data: socketmodels.CallLeft{},
			}
			delete(cs.CallsActive.data, uid)
		} else {
			for caller, called := range cs.CallsActive.data {
				if called == uid {
					ss.SendDataToUser <- socketserver.UserDataMessage{
						Type: "CALL_LEFT",
						Uid:  caller,
						Data: socketmodels.CallLeft{},
					}
					delete(cs.CallsActive.data, caller)
					break
				}
			}
		}
		cs.CallsActive.mutex.Unlock()
	}
}
