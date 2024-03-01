package keeper

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"polaris/x/blog/types"
	"strconv"
)

// TransmitIbcPostPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitIbcPostPacket(
	ctx sdk.Context,
	packetData types.IbcPostPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvIbcPostPacket processes packet reception
func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}
	id := k.AppendPost(
		ctx,
		types.Post{
			Creator: packet.SourcePort + "_" + packet.SourceChannel + "_" + data.Creator,
			Title:   data.Title,
			Content: data.Content,
		})

	// packet reception logic
	packetAck.PostID = strconv.FormatUint(id, 10)

	return packetAck, nil
}

// OnAcknowledgementIbcPostPacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.IbcPostPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// successful acknowledgement logic
		k.AppendSendPost(ctx, types.SendPost{
			Creator: data.Creator,
			PostID:  packetAck.PostID,
			Title:   data.Title,
			Chain:   packet.DestinationPort + "_" + packet.DestinationChannel + "_",
		})

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("the counter-party module does not implement the correct acknowledgment format")
	}
}

// OnTimeoutIbcPostPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {

	// packet timeout logic

	k.AppendTimedoutPost(ctx, types.TimedoutPost{
		Creator: data.Creator,
		Title:   data.Title,
		Chain:   packet.DestinationPort + "_" + packet.DestinationChannel,
	})

	return nil
}
