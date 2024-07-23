@0xb068ff5fb1c4f77e;

using Go = import "/go.capnp";
$Go.package("schema");
$Go.import("github.com/BradMyrick/chatnp/schema");

struct UserId {
  id @0 :Data;
}

struct ChatMessage {
  id @0 :UInt64;
  timestamp @1 :Int64;
  sender @2 :UserId;
  content @3 :Text;
}

struct ChatRoom {
  id @0 :UInt64;
  name @1 :Text;
  participants @2 :List(UserId);
  messages @3 :List(ChatMessage);
}

struct LocalChatHistory {
  rooms @0 :List(ChatRoom);
}

interface PeerDiscoveryService {
  discoverPeers @0 () -> (peerIds :List(UserId));
}

interface SecureMessagingService {
  sendMessage @0 (roomId :UInt64, msg :ChatMessage) -> ();
  getMessages @1 (roomId :UInt64, lastMessageId :UInt64) -> (messages :List(ChatMessage));
}

interface ChatRoomService {
  createRoom @0 (name :Text, participants :List(UserId)) -> (roomId :UInt64);
  joinRoom @1 (roomId :UInt64, userId :UserId) -> (success :Bool);
  leaveRoom @2 (roomId :UInt64, userId :UserId) -> ();
  getParticipants @3 (roomId :UInt64) -> (participants :List(UserId));
}

interface LocalHistoryService {
  saveHistory @0 (history :LocalChatHistory) -> ();
  loadHistory @1 () -> (history :LocalChatHistory);
}
