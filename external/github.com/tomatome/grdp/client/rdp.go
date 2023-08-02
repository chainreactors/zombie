package client

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
)

type RdpClient struct {
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	//channels *plugin.Channels
}

func newRdpClient(s *Setting) *RdpClient {
	return &RdpClient{}
}

func bitmapDecompress(bitmap *pdu.BitmapData) []byte {
	return core.Decompress(bitmap.BitmapDataStream, int(bitmap.Width), int(bitmap.Height), Bpp(bitmap.BitsPerPixel))
}
func split(user string) (domain string, uname string) {
	if strings.Index(user, "\\") != -1 {
		t := strings.Split(user, "\\")
		domain = t[0]
		uname = t[len(t)-1]
	} else if strings.Index(user, "/") != -1 {
		t := strings.Split(user, "/")
		domain = t[0]
		uname = t[len(t)-1]
	} else {
		uname = user
	}
	return
}
func (c *RdpClient) Login(host, user, pwd string, width, height int) error {
	conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}

	domain, user := split(user)
	c.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
	c.x224 = x224.New(c.tpkt)
	c.mcs = t125.NewMCSClient(c.x224)
	c.sec = sec.NewClient(c.mcs)
	c.pdu = pdu.NewClient(c.sec)
	//c.channels = plugin.NewChannels(c.sec)

	c.mcs.SetClientDesktop(uint16(width), uint16(height))

	c.sec.SetUser(user)
	c.sec.SetPwd(pwd)
	c.sec.SetDomain(domain)

	c.tpkt.SetFastPathListener(c.sec)
	c.sec.SetFastPathListener(c.pdu)
	c.sec.SetChannelSender(c.mcs)
	//c.channels.SetChannelSender(c.sec)

	//c.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)
	//c.x224.SetRequestedProtocol(x224.PROTOCOL_SSL)

	err = c.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	return nil
}
func (c *RdpClient) On(event string, f interface{}) {
	c.pdu.On(event, f)
}
func (c *RdpClient) KeyUp(sc int, name string) {
	p := &pdu.ScancodeKeyEvent{}
	p.KeyCode = uint16(sc)
	p.KeyboardFlags |= pdu.KBDFLAGS_RELEASE
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}
func (c *RdpClient) KeyDown(sc int, name string) {
	p := &pdu.ScancodeKeyEvent{}
	p.KeyCode = uint16(sc)
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}

func (c *RdpClient) MouseMove(x, y int) {
	p := &pdu.PointerEvent{}
	p.PointerFlags |= pdu.PTRFLAGS_MOVE
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}

func (c *RdpClient) MouseWheel(scroll, x, y int) {
	p := &pdu.PointerEvent{}
	p.PointerFlags |= pdu.PTRFLAGS_WHEEL
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}

func (c *RdpClient) MouseUp(button int, x, y int) {
	p := &pdu.PointerEvent{}

	switch button {
	case 0:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON1
	case 2:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON2
	case 1:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON3
	default:
		p.PointerFlags |= pdu.PTRFLAGS_MOVE
	}

	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}
func (c *RdpClient) MouseDown(button int, x, y int) {
	p := &pdu.PointerEvent{}

	p.PointerFlags |= pdu.PTRFLAGS_DOWN

	switch button {
	case 0:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON1
	case 2:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON2
	case 1:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON3
	default:
		p.PointerFlags |= pdu.PTRFLAGS_MOVE
	}

	p.XPos = uint16(x)
	p.YPos = uint16(y)
	c.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}
func (c *RdpClient) Close() {
	if c != nil && c.tpkt != nil {
		c.tpkt.Close()
	}
}
