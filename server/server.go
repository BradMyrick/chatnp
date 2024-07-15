package server

import (
	"context"

	"capnproto.org/go/capnp/v3/server"

	chat "github.com/BradMyrick/chatnp/capnp"
)

type chatServer struct {
	messages []chat.ChatMessage
}

func (s *chatServer) SendMessage(ctx context.Context, call chat.ChatService_sendMessage) error {
	msg, err := call.Args().Msg()
	if err != nil {
		return err
	}
	s.messages = append(s.messages, msg)
	return nil
}

func (s *chatServer) GetMessages(ctx context.Context, call chat.ChatService_getMessages) error {
	lastId := call.Args().LastMessageId()
	idx := 0
	for i, m := range s.messages {
		if m.Id() > lastId {
			idx = i
			break
		}
	}

	msgs := s.messages[idx:]
	result, err := call.AllocResults()
	if err != nil {
		return err
	}
	msgList, err := chat.NewChatMessage_List(result.Segment(), int32(len(msgs)))
	if err != nil {
		return err
	}
	for i := 0; i < len(msgs); i++ {
		if err := msgList.Set(i, msgs[i]); err != nil {
			return err
		}
	}
	return result.SetMessages(msgList)

}


func NewServer() (chat.ChatService, *server.Server) {
    s := &chatServer{}
    srv := server.New(chat.ChatService_Methods(nil, s), nil, nil)
    client := chat.ChatService_ServerToClient(s)
    return client, srv
}
