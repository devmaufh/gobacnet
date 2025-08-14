package gobacnet

import (
	"context"
	"fmt"
	"github.com/devmaufh/gobacnet/encoding"
	"github.com/devmaufh/gobacnet/types"
	"log"
	"time"
)

func (c *Client) WriteProperty(dest types.Device, req types.WritePropertyData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	invokeID, err := c.tsm.ID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get a transaction id: %v", err)
	}
	defer c.tsm.Put(invokeID)

	udp, err := c.localUDPAddress()
	if err != nil {
		return err
	}

	src := types.UDPToAddress(udp)
	enc := encoding.NewEncoder()
	enc.NPDU(types.NPDU{
		Version:               types.ProtocolVersion,
		Destination:           &dest.Addr,
		Source:                &src,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              types.Normal,
		HopCount:              types.DefaultHopCount,
	})

	enc.WriteProperty(uint8(invokeID), req)
	if enc.Error() != nil {
		return fmt.Errorf("error encoding WriteProperty request: %v", enc.Error())
	}

	_, err = c.send(dest.Addr, enc.Bytes())
	if err != nil {
		log.Printf("Error sending WriteProperty request: %v", err)
		return fmt.Errorf("error sending WriteProperty request: %v", err)
	}

	return nil
}
