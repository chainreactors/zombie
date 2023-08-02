// rfb.go
package client

import (
	"fmt"
	"net"
	"time"

	"github.com/tomatome/grdp/protocol/rfb"
)

type VncClient struct {
	vnc *rfb.RFB
}

func newVncClient(s *Setting) *VncClient {
	return &VncClient{}
}

func (c *VncClient) Login(host, user, pwd string, width, height int) error {
	conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}

	c.vnc = rfb.NewRFB(rfb.NewRFBConn(conn, pwd))

	err = c.vnc.Connect()
	if err != nil {
		return fmt.Errorf("[vnc connect err] %v", err)
	}

	return nil
}
func (c *VncClient) On(event string, f interface{}) {
	c.vnc.On(event, f)
}

func (c *VncClient) KeyUp(sc int, name string) {
	k := &rfb.KeyEvent{}
	k.Key = uint32(sc)
	c.vnc.SendKeyEvent(k)
}
func (c *VncClient) KeyDown(sc int, name string) {
	k := &rfb.KeyEvent{}
	k.DownFlag = 1
	k.Key = uint32(sc)
	c.vnc.SendKeyEvent(k)
}

func (c *VncClient) MouseMove(x, y int) {
	p := &rfb.PointerEvent{}
	time.Sleep(8 * time.Millisecond)
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.vnc.SendPointEvent(p)
}

func (c *VncClient) MouseWheel(scroll, x, y int) {
}

func (c *VncClient) MouseUp(button int, x, y int) {
	p := &rfb.PointerEvent{}

	switch button {
	case 0:
		p.Mask = 1
	case 2:
		p.Mask = 1<<3 - 1
	case 1:
		p.Mask = 1<<2 - 1
	default:
		p.Mask = 0
	}
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.vnc.SendPointEvent(p)
}
func (c *VncClient) MouseDown(button int, x, y int) {
	p := &rfb.PointerEvent{}

	switch button {
	case 0:
		p.Mask = 1
	case 2:
		p.Mask = 1<<3 - 1
	case 1:
		p.Mask = 1<<2 - 1
	default:
		p.Mask = 0
	}

	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.MouseMove(x, y)
	c.vnc.SendPointEvent(p)
}

func (c *VncClient) Close() {
	if c.vnc != nil {
		c.vnc.Close()
	}
}
