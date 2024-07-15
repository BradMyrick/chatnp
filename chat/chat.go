package chat

import (
	"log"
	"sort"

	"github.com/BradMyrick/chatnp/chat"
)

type chatServiceImpl struct {
	messages []chat.ChatMessage
}

func (s *chatServiceImpl) SendMessage(call chat.ChatService_sendMessage) error {
	msg, err := call.Params.Msg()
	if err != nil {
		return err
	}
	s.messages = append(s.messages, msg)
	return nil
}

func (s *chatServiceImpl) GetMessages(call chat.ChatService_getMessages) error {
	lastId := call.Params.LastMessageId()
	idx := sort.Search(len(s.messages), func(i int) bool {
		return s.messages[i].Id() > lastId
	})

	msgs := s.messages[idx:]
	result, err := chat.NewChatService_getMessages_Results(call.Results.Segment())
	if err != nil {
		return err
	}
	return result.SetMessages(msgs)
}

func ServerConnect() {
	server, err := chat.NewServer(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
	chatService := &chatServiceImpl{}
	chat.ChatService_ServerToClient(chatService).Export(server)

}
