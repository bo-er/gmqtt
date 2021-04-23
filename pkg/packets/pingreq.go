package packets

import (
	"fmt"
	"io"
	"strings"

	"github.com/DrmagicE/gmqtt/pkg/codes"
	log "github.com/sirupsen/logrus"
)

// Pingreq represents the MQTT Pingreq  packet
type Pingreq struct {
	FixHeader *FixHeader
}

func (p *Pingreq) String() string {
	return fmt.Sprintf("Pingreq")
}

// NewPingreqPacket returns a Pingreq instance by the given FixHeader and io.Reader
func NewPingreqPacket(fh *FixHeader, r io.Reader) (*Pingreq, error) {
	log.WithFields(log.Fields{
		"CodeAddress":  "pingreq.go",
		"PacketType":   fh.PacketType,
		"Flags":        fh.Flags,
		"RemainLength": fh.RemainLength,
		"ErrMalformed": fh.Flags != FlagReserved,
	}).Info("Received New Connect Packet")
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		log.Errorf("failed to copy io.Reader to stringBuilder. Err:%s", err)
	}
	log.Info("======================================================================================")
	log.Infoln("received packet is:", buf.String())
	log.Info("=======================================================================================")
	if fh.Flags != FlagReserved {
		return nil, codes.ErrMalformed
	}
	p := &Pingreq{FixHeader: fh}
	err = p.Unpack(r)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// NewPingresp returns a Pingresp struct
func (p *Pingreq) NewPingresp() *Pingresp {
	fh := &FixHeader{PacketType: PINGRESP, Flags: 0, RemainLength: 0}
	return &Pingresp{FixHeader: fh}
}

// Pack encodes the packet struct into bytes and writes it into io.Writer.
func (p *Pingreq) Pack(w io.Writer) error {
	p.FixHeader = &FixHeader{PacketType: PINGREQ, Flags: 0, RemainLength: 0}
	return p.FixHeader.Pack(w)
}

// Unpack read the packet bytes from io.Reader and decodes it into the packet struct.
func (p *Pingreq) Unpack(r io.Reader) error {
	if p.FixHeader.RemainLength != 0 {
		return codes.ErrMalformed
	}
	return nil
}
