package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	rpc "net/rpc"
	"os"
	"sync"

	"capnproto.org/go/capnp/v3"
	"github.com/BradMyrick/chatnp/schema"
)

type ChatServer struct {
	schema.ChatRoomService_Server
	userId    schema.UserId
	chatRooms sync.Map // map[uint64]*schema.ChatRoom
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		userId: generateUserId(),
	}
}

func (s *ChatServer) CreateRoom(ctx context.Context, call schema.ChatRoomService_createRoom) error {
	params := call.Args()

	name, err := params.Name()
	if err != nil {
		return err
	}
	participants, err := params.Participants()
	if err != nil {
		return err
	}
	roomId := s.createRoom(name, participants)
	result, err := call.AllocResults()
	if err != nil {
		return err
	}
	result.SetRoomId(roomId)
	return nil
}

func (s *ChatServer) JoinRoom(ctx context.Context, call schema.ChatRoomService_joinRoom) error {
	params := call.Args()

	roomId := params.RoomId()
	userId, err := params.UserId()
	if err != nil {
		return err
	}
	success := s.joinRoom(roomId, userId)
	result, err := call.AllocResults()
	if err != nil {
		return err
	}
	result.SetSuccess(success)
	return nil
}

func (s *ChatServer) LeaveRoom(ctx context.Context, call schema.ChatRoomService_leaveRoom) error {
	params := call.Args()

	roomId := params.RoomId()
	userId, err := params.UserId()
	if err != nil {
		return err
	}
	s.leaveRoom(roomId, userId)
	return nil
}

func (s *ChatServer) GetParticipants(ctx context.Context, call schema.ChatRoomService_getParticipants) error {
	params := call.Args()

	roomId := params.RoomId()
	participants := s.getParticipants(roomId)
	result, err := call.AllocResults()
	if err != nil {
		return err
	}
	result.SetParticipants(participants)
	return nil
}

func (s *ChatServer) createRoom(name string, participants schema.UserId_List) uint64 {
	roomId := generateRoomId()
    var b []byte
    _, seg := capnp.NewSingleSegmentMessage(b)
	room, err := schema.NewChatRoom(seg)
	if err != nil {
		log.Fatal(err)
	}
	room.SetId(roomId)
	room.SetName(name)
	room.SetParticipants(participants)

	s.chatRooms.Store(roomId, room)
	return roomId
}

func (s *ChatServer) joinRoom(roomId uint64, userId schema.UserId) bool {
	room, ok := s.chatRooms.Load(roomId)
	if !ok {
		return false
	}
    // participants are the authorized userId's that can join
	participants, err := room.(*schema.ChatRoom).Participants()
	if err != nil {
		log.Println("Error getting participants:", err)
		return false
	}

    // check if the userId is in the authorized participants.  
    // if so join the room
	for i := 0; i < participants.Len(); i++ {
		participant := participants.At(i)
		participantId, err := participant.Id()
		if err != nil {
			log.Println("Error getting participant ID:", err)
			return false
		}
		suid, err := userId.Id()
		if err != nil {
			log.Println("Error getting user ID:", err)
			return false
		}
		if bytes.Equal(suid, participantId) {
			return true
		}
	}
    return false
}


func (s *ChatServer) getParticipants(roomId uint64) schema.UserId_List {
	room, ok := s.chatRooms.Load(roomId)
	if !ok {
		return schema.UserId_List{}
	}
	participants, err := room.(*schema.ChatRoom).Participants()
	if err != nil {
		log.Println("Error getting participants:", err)
		return schema.UserId_List{}
	}
	return participants
}

func generateUserId() schema.UserId {
    var b []byte
    _, seg := capnp.NewSingleSegmentMessage(b)
	id, err := schema.NewUserId(seg)
	if err != nil {
		log.Fatal(err)
	}
	b = make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	err = id.SetId(b)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func generateRoomId() uint64 {
	
	return uint64(0) // TODO
}

func StartServer() {
	s := NewChatServer()
    suid, err := s.userId.Id()
    if err != nil {
        log.Fatal(err)
    } 
	fmt.Printf("Generated user ID: %s\n", hex.EncodeToString(suid))

	go rpc.Register(s)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	for {
		fmt.Println("Enter command (join <roomId>, create <name>, quit):")
		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "join":
			var roomId uint64
			fmt.Scanln(&roomId)
			// TODO: check if room exists
			_, seg := capnp.NewSingleSegmentMessage(nil)
			userId, err := schema.NewUserId(seg)
			if err != nil {
				log.Fatal(err)
			}
			var id []byte
			id, err = s.userId.Id()
			if err != nil {
				log.Fatal(err)
			}
			userId.SetId(id)
			if s.joinRoom(roomId, userId) {
				log.Printf("Joined room %d\n", roomId)
			} else {
				log.Printf("Failed to join room %d\n", roomId)
			} 
		case "create":
			var name string
			fmt.Scanln(&name)
			_, seg := capnp.NewSingleSegmentMessage(nil)

			userId, err := schema.NewUserId(seg)
			if err != nil {
				log.Fatal(err)
			}
			var id []byte
			id, err = s.userId.Id()
			if err != nil {
				log.Fatal(err)
			}
			userId.SetId(id)
			_, seg = capnp.NewSingleSegmentMessage(nil)
			participants, err := schema.NewUserId_List(seg, 1)
			if err != nil {
				log.Fatal(err)
			}
			participants.Set(0, userId)
			roomId := s.createRoom(name, participants)
			log.Printf("Created room %d\n", roomId)
		case "quit":
			log.Println("Quitting")
			l.Close()
			os.Exit(0)
		default:
			log.Println("Unknown command")
		}
	}
}
