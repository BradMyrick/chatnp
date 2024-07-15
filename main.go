package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"capnproto.org/go/capnp/v3"
)


func main() {
	clientHook := errorClient{e}
	conn := capnp.NewClient(capnp.ClientHook{})
	go func() {
	  var lastId uint64
	  for {
		msgs, err := client.GetMessages(ctx, func(p chat.ChatService_getMessages_Params) error {
		  return p.SetLastMessageId(lastId)
		}).Struct()
		
		if err != nil {
		  log.Println("Failed to get messages:", err)
		  continue
		}
		
		for _, m := range msgs.Messages() {
		  log.Printf("%s: %s\n", m.Sender(), m.Content())
		  if m.Id() > lastId {
			lastId = m.Id()
		  }
		}
		
		time.Sleep(time.Second)
	  }
	}()
	
	reader := bufio.NewReader(os.Stdin)
	for {
	  fmt.Print("> ")
	  text, _ := reader.ReadString('\n')
	  
	  _, err := client.SendMessage(ctx, func(p chat.ChatService_sendMessage_Params) error {
		msg, err := p.NewMsg()
		if err != nil {
		  return err
		}
		
		msg.SetId(uint64(time.Now().UnixNano()))
		msg.SetTimestamp(time.Now().Unix())
		msg.SetSender("user")
		msg.SetContent(strings.TrimSpace(text))
		
		return nil
	  })
	  
	  if err != nil {
		log.Println("Failed to send message:", err)
	  }
	}
  }
  