// client.go
package client

import (
	"log"
	"os"

	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/rfb"
)

const (
	CLIP_OFF = 0
	CLIP_IN  = 0x1
	CLIP_OUT = 0x2
)

const (
	TC_RDP = 0
	TC_VNC = 1
)

type Control interface {
	Login(host, user, passwd string, width, height int) error
	KeyUp(sc int, name string)
	KeyDown(sc int, name string)
	MouseMove(x, y int)
	MouseWheel(scroll, x, y int)
	MouseUp(button int, x, y int)
	MouseDown(button int, x, y int)
	On(event string, msg interface{})
	Close()
}

func init() {
	glog.SetLevel(glog.INFO)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
}

type Client struct {
	host    string
	user    string
	passwd  string
	ctl     Control
	tc      int
	setting *Setting
}

func NewClient(host, user, passwd string, t int, s *Setting) *Client {
	if s == nil {
		s = NewSetting()
	}
	c := &Client{
		host:    host,
		user:    user,
		passwd:  passwd,
		tc:      t,
		setting: s,
	}

	switch t {
	case TC_VNC:
		c.ctl = newVncClient(s)
	default:
		c.ctl = newRdpClient(s)
	}

	s.SetLogLevel()
	return c
}

func (c *Client) Login() error {
	return c.ctl.Login(c.host, c.user, c.passwd, c.setting.Width, c.setting.Height)
}

func (c *Client) KeyUp(sc int, name string) {
	c.ctl.KeyUp(sc, name)
}
func (c *Client) KeyDown(sc int, name string) {
	c.ctl.KeyDown(sc, name)
}
func (c *Client) MouseMove(x, y int) {
	c.ctl.MouseMove(x, y)
}
func (c *Client) MouseWheel(scroll, x, y int) {
	c.ctl.MouseWheel(scroll, x, y)
}
func (c *Client) MouseUp(button, x, y int) {
	c.ctl.MouseUp(button, x, y)
}
func (c *Client) MouseDown(button, x, y int) {
	c.ctl.MouseDown(button, x, y)
}
func (c *Client) OnError(f func(e error)) {
	c.ctl.On("error", f)
}
func (c *Client) OnClose(f func()) {
	c.ctl.On("close", f)
}
func (c *Client) OnSuccess(f func()) {
	c.ctl.On("success", f)
}
func (c *Client) OnReady(f func()) {
	c.ctl.On("ready", f)
}
func (c *Client) OnBitmap(f func([]Bitmap)) {
	f1 := func(data interface{}) {
		bs := make([]Bitmap, 0, 50)
		if c.tc == TC_VNC {
			br := data.(*rfb.BitRect)
			for _, v := range br.Rects {
				b := Bitmap{int(v.Rect.X), int(v.Rect.Y), int(v.Rect.X + v.Rect.Width), int(v.Rect.Y + v.Rect.Height),
					int(v.Rect.Width), int(v.Rect.Height),
					Bpp(uint16(br.Pf.BitsPerPixel)), false, v.Data}
				bs = append(bs, b)
			}
		} else {
			for _, v := range data.([]pdu.BitmapData) {
				IsCompress := v.IsCompress()
				stream := v.BitmapDataStream
				if IsCompress {
					stream = bitmapDecompress(&v)
					IsCompress = false
				}

				b := Bitmap{int(v.DestLeft), int(v.DestTop), int(v.DestRight), int(v.DestBottom),
					int(v.Width), int(v.Height), Bpp(v.BitsPerPixel), IsCompress, stream}
				bs = append(bs, b)
			}
		}
		f(bs)
	}

	c.ctl.On("bitmap", f1)
}

type Bitmap struct {
	DestLeft     int    `json:"destLeft"`
	DestTop      int    `json:"destTop"`
	DestRight    int    `json:"destRight"`
	DestBottom   int    `json:"destBottom"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	BitsPerPixel int    `json:"bitsPerPixel"`
	IsCompress   bool   `json:"isCompress"`
	Data         []byte `json:"data"`
}

func Bpp(bp uint16) int {
	return int(bp / 8)
}

type Setting struct {
	Width    int
	Height   int
	Protocol string
	LogLevel glog.LEVEL
}

func NewSetting() *Setting {
	return &Setting{
		Width:    1024,
		Height:   768,
		LogLevel: glog.NONE,
	}
}
func (s *Setting) SetLogLevel() {
	glog.SetLevel(s.LogLevel)
}

func (s *Setting) SetRequestedProtocol(p uint32) {}
func (s *Setting) SetClipboard(c int)            {}
