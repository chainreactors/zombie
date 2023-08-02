package drdynvc

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/plugin"
)

const (
	ChannelName   = plugin.DRDYNVC_SVC_CHANNEL_NAME
	ChannelOption = plugin.CHANNEL_OPTION_INITIALIZED |
		plugin.CHANNEL_OPTION_ENCRYPT_RDP
)

const (
	MAX_DVC_CHANNELS = 20
)

const (
	DYNVC_CREATE_REQ            = 0x01
	DYNVC_DATA_FIRST            = 0x02
	DYNVC_DATA                  = 0x03
	DYNVC_CLOSE                 = 0x04
	DYNVC_CAPABILITIES          = 0x05
	DYNVC_DATA_FIRST_COMPRESSED = 0x06
	DYNVC_DATA_COMPRESSED       = 0x07
	DYNVC_SOFT_SYNC_REQUEST     = 0x08
	DYNVC_SOFT_SYNC_RESPONSE    = 0x09
)

type ChannelClient struct {
	name          string
	id            uint32
	channelSender core.ChannelSender
}

type DvcClient struct {
	w        core.ChannelSender
	channels map[string]ChannelClient
}

func NewDvcClient() *DvcClient {
	return &DvcClient{
		channels: make(map[string]ChannelClient, 100),
	}
}

func (c *DvcClient) LoadAddin(f core.ChannelSender) {

}

type DvcHeader struct {
	cmd    uint8
	sp     uint8
	cbChId uint8
}

func readHeader(r io.Reader) *DvcHeader {
	value, _ := core.ReadUInt8(r)
	cmd := (value & 0xf0) >> 4
	sp := (value & 0x0c) >> 2
	cbChId := (value & 0x03) >> 0
	return &DvcHeader{cmd, sp, cbChId}
}

func (h *DvcHeader) serialize(channelId uint32) []byte {
	b := &bytes.Buffer{}
	core.WriteUInt8((h.cmd<<4)|(h.sp<<2)|h.cbChId, b)
	if h.cbChId == 0 {
		core.WriteUInt8(uint8(channelId), b)
	} else if h.cbChId == 1 {
		core.WriteUInt16LE(uint16(channelId), b)
	} else {
		core.WriteUInt32LE(channelId, b)
	}

	return b.Bytes()
}

func (c *DvcClient) Send(s []byte) (int, error) {
	glog.Debug("len:", len(s), "data:", hex.EncodeToString(s))
	name, _ := c.GetType()
	return c.w.SendToChannel(name, s)
}
func (c *DvcClient) Sender(f core.ChannelSender) {
	c.w = f
}
func (c *DvcClient) GetType() (string, uint32) {
	return ChannelName, ChannelOption
}

func (c *DvcClient) Process(s []byte) {
	glog.Debug("recv:", hex.EncodeToString(s))
	r := bytes.NewReader(s)
	hdr := readHeader(r)
	glog.Infof("dvc: Cmd=0x%x, Sp=%d CbChId=%d all=%d", hdr.cmd, hdr.sp, hdr.cbChId, r.Len())

	b, _ := core.ReadBytes(r.Len(), r)

	switch hdr.cmd {
	case DYNVC_CAPABILITIES:
		glog.Info("DYNVC_CAPABILITIES")
		c.processCapsPdu(hdr, b)
	case DYNVC_CREATE_REQ:
		glog.Info("DYNVC_CREATE_REQ")
		c.processCreateReq(hdr, b)
	case DYNVC_DATA_FIRST:
		glog.Info("DYNVC_DATA_FIRST")
	case DYNVC_DATA:
		glog.Info("DYNVC_DATA")
	case DYNVC_CLOSE:
		glog.Info("DYNVC_CLOSE")
	default:
		glog.Errorf("type 0x%x not supported", hdr.cmd)
	}
}
func (c *DvcClient) processCreateReq(hdr *DvcHeader, s []byte) {
	r := bytes.NewReader(s)
	channelId := readDvcId(r, hdr.cbChId)
	name, _ := core.ReadBytes(r.Len(), r)
	channelName := string(name)
	glog.Infof("Server requests channelId=%d, name=%s", channelId, channelName)

	//response
	b := &bytes.Buffer{}
	b.Write(hdr.serialize(channelId))
	core.WriteUInt32LE(0, b)
	c.Send(b.Bytes())
}

func readDvcId(r io.Reader, cbLen uint8) (id uint32) {
	switch cbLen {
	case 0:
		i, _ := core.ReadUInt8(r)
		id = uint32(i)
	case 1:
		i, _ := core.ReadUint16LE(r)
		id = uint32(i)
	default:
		id, _ = core.ReadUInt32LE(r)
	}
	return
}
func (c *DvcClient) processCapsPdu(hdr *DvcHeader, s []byte) {
	r := bytes.NewReader(s)
	core.ReadUInt8(r)
	ver, _ := core.ReadUint16LE(r)
	glog.Infof("Server supports dvc=%d", ver)

	hdr.cmd = DYNVC_CAPABILITIES
	hdr.cbChId = 0
	hdr.sp = 0

	b := &bytes.Buffer{}
	core.WriteUInt16LE(0x0050, b)
	core.WriteUInt16LE(ver, b)
	c.Send(b.Bytes())
}
