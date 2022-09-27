package lbm

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/v2/relayer/provider"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap/zapcore"
)

var _ provider.RelayerMessage = &LBMMessage{}

type LBMMessage struct {
	Msg sdk.Msg
}

func NewLBMMessage(msg sdk.Msg) provider.RelayerMessage {
	return LBMMessage{
		Msg: msg,
	}
}

func LBMMsg(rm provider.RelayerMessage) sdk.Msg {
	if val, ok := rm.(LBMMessage); !ok {
		fmt.Printf("got data of type %T but wanted provider.LBMMessage \n", val)
		return nil
	} else {
		return val.Msg
	}
}

func LBMMsgs(rm ...provider.RelayerMessage) []sdk.Msg {
	sdkMsgs := make([]sdk.Msg, 0)
	for _, rMsg := range rm {
		if val, ok := rMsg.(LBMMessage); !ok {
			fmt.Printf("got data of type %T but wanted provider.LBMMessage \n", val)
			return nil
		} else {
			sdkMsgs = append(sdkMsgs, val.Msg)
		}
	}
	return sdkMsgs
}

func (cm LBMMessage) Type() string {
	return sdk.MsgTypeURL(cm.Msg)
}

func (cm LBMMessage) MsgBytes() ([]byte, error) {
	return proto.Marshal(cm.Msg)
}

// MarshalLogObject is used to encode cm to a zap logger with the zap.Object field type.
func (cm LBMMessage) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Using plain json.Marshal or calling cm.Msg.String() both fail miserably here.
	// There is probably a better way to encode the message than this.
	j, err := codec.NewLegacyAmino().MarshalJSON(cm.Msg)
	if err != nil {
		return err
	}
	enc.AddByteString("msg_json", j)
	return nil
}
