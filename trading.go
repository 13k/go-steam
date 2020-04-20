package steam

import (
	pb "github.com/13k/go-steam-resources/protobuf/steam"
	"github.com/13k/go-steam-resources/steamlang"
	"github.com/13k/go-steam/protocol"
	"github.com/13k/go-steam/steamid"
	"google.golang.org/protobuf/proto"
)

// Trading provides access to the Steam client's part of Steam Trading, that is bootstrapping the
// trade.
//
// The trade itself is not handled by the Steam client itself, but it's a part of the Steam website.
//
// You'll receive a TradeProposedEvent when a friend proposes a trade. You can accept it with the
// RespondRequest method. You can request a trade yourself with RequestTrade.
type Trading struct {
	client *Client
}

type TradeRequestId uint32

func (t *Trading) HandlePacket(packet *protocol.Packet) {
	switch packet.EMsg {
	case steamlang.EMsg_EconTrading_InitiateTradeProposed:
		msg := &pb.CMsgTrading_InitiateTradeRequest{}
		packet.ReadProtoMsg(msg)
		t.client.Emit(&TradeProposedEvent{
			RequestId: TradeRequestId(msg.GetTradeRequestId()),
			Other:     steamid.SteamId(msg.GetOtherSteamid()),
		})
	case steamlang.EMsg_EconTrading_InitiateTradeResult:
		msg := &pb.CMsgTrading_InitiateTradeResponse{}
		packet.ReadProtoMsg(msg)
		t.client.Emit(&TradeResultEvent{
			RequestId: TradeRequestId(msg.GetTradeRequestId()),
			Response:  steamlang.EEconTradeResponse(msg.GetResponse()),
			Other:     steamid.SteamId(msg.GetOtherSteamid()),

			NumDaysSteamGuardRequired:            msg.GetSteamguardRequiredDays(),
			NumDaysNewDeviceCooldown:             msg.GetNewDeviceCooldownDays(),
			DefaultNumDaysPasswordResetProbation: msg.GetDefaultPasswordResetProbationDays(),
			NumDaysPasswordResetProbation:        msg.GetPasswordResetProbationDays(),
		})
	case steamlang.EMsg_EconTrading_StartSession:
		msg := &pb.CMsgTrading_StartSession{}
		packet.ReadProtoMsg(msg)
		t.client.Emit(&TradeSessionStartEvent{
			Other: steamid.SteamId(msg.GetOtherSteamid()),
		})
	}
}

// Requests a trade. You'll receive a TradeResultEvent if the request fails or
// if the friend accepted the trade.
func (t *Trading) RequestTrade(other steamid.SteamId) {
	pbmsg := &pb.CMsgTrading_InitiateTradeRequest{
		OtherSteamid: proto.Uint64(uint64(other)),
	}

	t.client.Write(protocol.NewClientMsgProtobuf(steamlang.EMsg_EconTrading_InitiateTradeRequest, pbmsg))
}

// Responds to a TradeProposedEvent.
func (t *Trading) RespondRequest(requestId TradeRequestId, accept bool) {
	var resp uint32

	if !accept {
		resp = 1
	}

	pbmsg := &pb.CMsgTrading_InitiateTradeResponse{
		TradeRequestId: proto.Uint32(uint32(requestId)),
		Response:       proto.Uint32(resp),
	}

	t.client.Write(protocol.NewClientMsgProtobuf(steamlang.EMsg_EconTrading_InitiateTradeResponse, pbmsg))
}

// This cancels a request made with RequestTrade.
func (t *Trading) CancelRequest(other steamid.SteamId) {
	pbmsg := &pb.CMsgTrading_CancelTradeRequest{
		OtherSteamid: proto.Uint64(uint64(other)),
	}

	t.client.Write(protocol.NewClientMsgProtobuf(steamlang.EMsg_EconTrading_CancelTradeRequest, pbmsg))
}
