@0xb068ff5fb1c4f77e;

using Go = import "/go.capnp";
$Go.package("chat");
$Go.import("github.com/BradMyrick/chatnp/chat");

struct ChatMessage {
  id @0 :UInt64;
  timestamp @1 :Int64;
  sender @2 :Text;
  content @3 :Text;
}

interface ChatService {
  sendMessage @0 (msg :ChatMessage) -> ();
  getMessages @1 (lastMessageId :UInt64) -> (messages :List(ChatMessage));
}
