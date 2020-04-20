package steam

import (
	"bytes"

	pb "github.com/13k/go-steam-resources/protobuf/steam"
	"github.com/13k/go-steam-resources/steamlang"
	"github.com/13k/go-steam/protocol"
	gc "github.com/13k/go-steam/protocol/gamecoordinator"
	"google.golang.org/protobuf/proto"
)

type GameCoordinator struct {
	client   *Client
	handlers []GCPacketHandler
}

func newGC(client *Client) *GameCoordinator {
	return &GameCoordinator{
		client:   client,
		handlers: make([]GCPacketHandler, 0),
	}
}

type GCPacketHandler interface {
	HandleGCPacket(*gc.GCPacket)
}

func (g *GameCoordinator) RegisterPacketHandler(handler GCPacketHandler) {
	g.handlers = append(g.handlers, handler)
}

func (g *GameCoordinator) HandlePacket(packet *protocol.Packet) {
	if packet.EMsg != steamlang.EMsg_ClientFromGC {
		return
	}

	msg := &pb.CMsgGCClient{}

	packet.ReadProtoMsg(msg)

	p, err := gc.NewGCPacket(msg)

	if err != nil {
		g.client.Errorf("Error reading GC message: %v", err)
		return
	}

	for _, handler := range g.handlers {
		handler.HandleGCPacket(p)
	}
}

func (g *GameCoordinator) Write(msg gc.IGCMsg) {
	buf := &bytes.Buffer{}

	msg.Serialize(buf)

	msgType := msg.GetMsgType()

	if msg.IsProto() {
		msgType = msgType | 0x80000000 // mask with protoMask
	}

	g.client.Write(protocol.NewClientMsgProtobuf(steamlang.EMsg_ClientToGC, &pb.CMsgGCClient{
		Msgtype: proto.Uint32(msgType),
		Appid:   proto.Uint32(msg.GetAppId()),
		Payload: buf.Bytes(),
	}))
}

// Sets you in the given games. Specify none to quit all games.
func (g *GameCoordinator) SetGamesPlayed(appIds ...uint64) {
	games := make([]*pb.CMsgClientGamesPlayed_GamePlayed, 0)

	for _, appId := range appIds {
		games = append(games, &pb.CMsgClientGamesPlayed_GamePlayed{
			GameId: proto.Uint64(appId),
		})
	}

	g.client.Write(protocol.NewClientMsgProtobuf(steamlang.EMsg_ClientGamesPlayed, &pb.CMsgClientGamesPlayed{
		GamesPlayed: games,
	}))
}
