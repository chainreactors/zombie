package rdpgfx

import (
	"encoding/hex"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/plugin"
)

const (
	ChannelName = plugin.RDPGFX_DVC_CHANNEL_NAME
)

type gfxClient struct {
}

func (c *gfxClient) Send(s []byte) (int, error) {
	glog.Debug("len:", len(s), "data:", hex.EncodeToString(s))
	name, _ := c.GetType()
	return c.w.SendToChannel(name, s)
}
func (c *gfxClient) Sender(f core.ChannelSender) {
	c.w = f
}
func (c *gfxClient) GetType() (string, uint32) {
	return ChannelName, 0
}
