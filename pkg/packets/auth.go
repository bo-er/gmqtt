package packets

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/DrmagicE/gmqtt/pkg/codes"
	log "github.com/sirupsen/logrus"
)

type Auth struct {
	FixHeader  *FixHeader
	Code       byte
	Properties *Properties
}

func (a *Auth) String() string {
	return fmt.Sprintf("Auth, Code: %v, Properties: %s", a.Code, a.Properties)
}

func (a *Auth) Pack(w io.Writer) error {
	a.FixHeader = &FixHeader{PacketType: AUTH, Flags: FlagReserved}
	bufw := &bytes.Buffer{}
	if a.Code != codes.Success || a.Properties != nil {
		bufw.WriteByte(a.Code)
		a.Properties.Pack(bufw, AUTH)
	}
	a.FixHeader.RemainLength = bufw.Len()
	err := a.FixHeader.Pack(w)
	if err != nil {
		return err
	}
	_, err = bufw.WriteTo(w)
	return err
}

func (a *Auth) Unpack(r io.Reader) error {
	if a.FixHeader.RemainLength == 0 {
		a.Code = codes.Success
		return nil
	}
	restBuffer := make([]byte, a.FixHeader.RemainLength)
	_, err := io.ReadFull(r, restBuffer)
	if err != nil {
		return codes.ErrMalformed
	}
	bufr := bytes.NewBuffer(restBuffer)
	a.Code, err = bufr.ReadByte()
	if err != nil {
		return codes.ErrMalformed
	}
	if !ValidateCode(AUTH, a.Code) {
		return codes.ErrProtocol
	}
	a.Properties = &Properties{}
	return a.Properties.Unpack(bufr, AUTH)
}

func NewAuthPacket(fh *FixHeader, r io.Reader) (*Auth, error) {
	p := &Auth{FixHeader: fh}
	log.WithFields(log.Fields{
		"CodeAddress":  "auth.go",
		"PacketType":   fh.PacketType,
		"Flags":        fh.Flags,
		"RemainLength": fh.RemainLength,
		"ErrMalformed":fh.Flags != FlagReserved,
	}).Info("Received New Connect Packet")
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		log.Errorf("failed to copy io.Reader to stringBuilder. Err:%s", err)
	}
	log.Info("======================================================================================")
	log.Infoln("received packet is:", buf.String())
	log.Info("=======================================================================================")
	//判断 标志位 flags 是否合法[MQTT-2.2.2-2]
	if fh.Flags != FlagReserved {
		return nil, codes.ErrMalformed
	}
	err = p.Unpack(r)
	if err != nil {
		return nil, err
	}
	return p, err
}
