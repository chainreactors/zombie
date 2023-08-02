package plugin

import "C"
import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/tomatome/grdp/glog"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/emission"
)

const (
	CHANNEL_RC_OK                         = 0
	CHANNEL_RC_ALREADY_INITIALIZED        = 1
	CHANNEL_RC_NOT_INITIALIZED            = 2
	CHANNEL_RC_ALREADY_CONNECTED          = 3
	CHANNEL_RC_NOT_CONNECTED              = 4
	CHANNEL_RC_TOO_MANY_CHANNELS          = 5
	CHANNEL_RC_BAD_CHANNEL                = 6
	CHANNEL_RC_BAD_CHANNEL_HANDLE         = 7
	CHANNEL_RC_NO_BUFFER                  = 8
	CHANNEL_RC_BAD_INIT_HANDLE            = 9
	CHANNEL_RC_NOT_OPEN                   = 10
	CHANNEL_RC_BAD_PROC                   = 11
	CHANNEL_RC_NO_MEMORY                  = 12
	CHANNEL_RC_UNKNOWN_CHANNEL_NAME       = 13
	CHANNEL_RC_ALREADY_OPEN               = 14
	CHANNEL_RC_NOT_IN_VIRTUALCHANNELENTRY = 15
	CHANNEL_RC_NULL_DATA                  = 16
	CHANNEL_RC_ZERO_LENGTH                = 17
	CHANNEL_RC_INVALID_INSTANCE           = 18
	CHANNEL_RC_UNSUPPORTED_VERSION        = 19
	CHANNEL_RC_INITIALIZATION_ERROR       = 20
)
const (
	VIRTUAL_CHANNEL_VERSION_WIN2000 = 1
)

const (
	CHANNEL_EVENT_INITIALIZED          = 0
	CHANNEL_EVENT_CONNECTED            = 1
	CHANNEL_EVENT_V1_CONNECTED         = 2
	CHANNEL_EVENT_DISCONNECTED         = 3
	CHANNEL_EVENT_TERMINATED           = 4
	CHANNEL_EVENT_REMOTE_CONTROL_START = 5
	CHANNEL_EVENT_REMOTE_CONTROL_STOP  = 6
	CHANNEL_EVENT_ATTACHED             = 7
	CHANNEL_EVENT_DETACHED             = 8
	CHANNEL_EVENT_DATA_RECEIVED        = 10
	CHANNEL_EVENT_WRITE_COMPLETE       = 11
	CHANNEL_EVENT_WRITE_CANCELLED      = 12
)

const (
	CHANNEL_OPTION_INITIALIZED               = 0x80000000
	CHANNEL_OPTION_ENCRYPT_RDP               = 0x40000000
	CHANNEL_OPTION_ENCRYPT_SC                = 0x20000000
	CHANNEL_OPTION_ENCRYPT_CS                = 0x10000000
	CHANNEL_OPTION_PRI_HIGH                  = 0x08000000
	CHANNEL_OPTION_PRI_MED                   = 0x04000000
	CHANNEL_OPTION_PRI_LOW                   = 0x02000000
	CHANNEL_OPTION_COMPRESS_RDP              = 0x00800000
	CHANNEL_OPTION_COMPRESS                  = 0x00400000
	CHANNEL_OPTION_SHOW_PROTOCOL             = 0x00200000
	CHANNEL_OPTION_REMOTE_CONTROL_PERSISTENT = 0x00100000
)

type ChannelDef struct {
	Name    string
	Options uint32
}
type CHANNEL_INIT_EVENT_EX_FN func(lpUserParam interface{},
	pInitHandle interface{}, event uint, pData uintptr, dataLength uint)
type VIRTUALCHANNELINITEX func(lpUserParam interface{}, clientContext interface{},
	pInitHandle interface{}, pChannel []ChannelDef,
	channelCount int, versionRequested uint32,
	pChannelInitEventProcEx CHANNEL_INIT_EVENT_EX_FN) uint

type CHANNEL_OPEN_EVENT_EX_FN func(lpUserParam uintptr,
	openHandle uint32, event uint,
	pData uintptr, dataLength uint32, totalLength uint32, dataFlags uint32)
type VIRTUALCHANNELOPENEX func(pInitHandle interface{}, pOpenHandle *uint32,
	pChannelName string,
	pChannelOpenEventProcEx *CHANNEL_OPEN_EVENT_EX_FN) uint

type VIRTUALCHANNELCLOSEEX func(pInitHandle interface{}, openHandle uint32) uint

type VIRTUALCHANNELWRITEEX func(pInitHandle interface{}, openHandle uint32, pData interface{},
	dataLength uint32, pUserData interface{}) uint

type ChannelEntryPointsEx struct {
	CbSize                 uint32
	ProtocolVersion        uint32
	PVirtualChannelInitEx  VIRTUALCHANNELINITEX
	PVirtualChannelOpenEx  VIRTUALCHANNELOPENEX
	PVirtualChannelCloseEx VIRTUALCHANNELCLOSEEX
	PVirtualChannelWriteEx VIRTUALCHANNELWRITEEX
}

func NewChannelEntryPointsEx() *ChannelEntryPointsEx {
	e := &ChannelEntryPointsEx{}
	e.CbSize = uint32(unsafe.Sizeof(e))
	e.ProtocolVersion = VIRTUAL_CHANNEL_VERSION_WIN2000
	return e
}

type VIRTUALCHANNELENTRYEX func(pEntryPointsEx *ChannelEntryPointsEx,
	pInitHandle interface{}) error

/*
type ChannelEntryPoints struct {
	CbSize               uint32
	ProtocolVersion      uint32
	PVirtualChannelInit  PVIRTUALCHANNELINIT
	PVirtualChannelOpen  PVIRTUALCHANNELOPEN
	PVirtualChannelClose PVIRTUALCHANNELCLOSE
	PVirtualChannelWrite PVIRTUALCHANNELWRITE
}
typedef VOID VCAPITYPE CHANNEL_INIT_EVENT_FN(LPVOID pInitHandle,
        UINT event, LPVOID pData, UINT dataLength);

typedef CHANNEL_INIT_EVENT_FN* PCHANNEL_INIT_EVENT_FN;
typedef VOID VCAPITYPE CHANNEL_OPEN_EVENT_FN(DWORD openHandle, UINT event,
        LPVOID pData, UINT32 dataLength, UINT32 totalLength, UINT32 dataFlags);

typedef CHANNEL_OPEN_EVENT_FN* PCHANNEL_OPEN_EVENT_FN;
typedef UINT VCAPITYPE VIRTUALCHANNELINIT(LPVOID* ppInitHandle, PCHANNEL_DEF pChannel,
                                          INT channelCount, ULONG versionRequested,
                                          PCHANNEL_INIT_EVENT_FN pChannelInitEventProc);
typedef VIRTUALCHANNELINIT* PVIRTUALCHANNELINIT;

typedef UINT VCAPITYPE VIRTUALCHANNELOPEN(LPVOID pInitHandle, LPDWORD pOpenHandle,
                                          PCHAR pChannelName,
                                          PCHANNEL_OPEN_EVENT_FN pChannelOpenEventProc);

typedef VIRTUALCHANNELOPEN* PVIRTUALCHANNELOPEN;

typedef UINT VCAPITYPE VIRTUALCHANNELCLOSE(DWORD openHandle);
typedef VIRTUALCHANNELCLOSE* PVIRTUALCHANNELCLOSE;

typedef UINT VCAPITYPE VIRTUALCHANNELWRITE(DWORD openHandle, LPVOID pData, ULONG dataLength,
                                           LPVOID pUserData);
typedef VIRTUALCHANNELWRITE* PVIRTUALCHANNELWRITE;

typedef UINT VCAPITYPE VIRTUALCHANNELINITEX(LPVOID lpUserParam, LPVOID clientContext,
                                            LPVOID pInitHandle, PCHANNEL_DEF pChannel,
                                            INT channelCount, ULONG versionRequested,
                                            PCHANNEL_INIT_EVENT_EX_FN pChannelInitEventProcEx);
typedef VIRTUALCHANNELINITEX* PVIRTUALCHANNELINITEX;

typedef UINT VCAPITYPE VIRTUALCHANNELOPENEX(LPVOID pInitHandle, LPDWORD pOpenHandle,
                                            PCHAR pChannelName,
                                            PCHANNEL_OPEN_EVENT_EX_FN pChannelOpenEventProcEx);
typedef VIRTUALCHANNELOPENEX* PVIRTUALCHANNELOPENEX;


typedef UINT VCAPITYPE VIRTUALCHANNELCLOSEEX(LPVOID pInitHandle, DWORD openHandle);
typedef VIRTUALCHANNELCLOSEEX* PVIRTUALCHANNELCLOSEEX;

typedef UINT VCAPITYPE VIRTUALCHANNELWRITEEX(LPVOID pInitHandle, DWORD openHandle, LPVOID pData,
                                             ULONG dataLength, LPVOID pUserData);
typedef VIRTUALCHANNELWRITEEX* PVIRTUALCHANNELWRITEEX;
*/

//static channel name
const (
	CLIPRDR_SVC_CHANNEL_NAME = "cliprdr" //剪切板
	RDPDR_SVC_CHANNEL_NAME   = "rdpdr"   //设备重定向(打印机，磁盘，端口，智能卡等)
	RDPSND_SVC_CHANNEL_NAME  = "rdpsnd"  //音频输出
	RAIL_SVC_CHANNEL_NAME    = "rail"    //远程应用
	DRDYNVC_SVC_CHANNEL_NAME = "drdynvc" //动态虚拟通道
	REMDESK_SVC_CHANNEL_NAME = "remdesk" //远程协助
)

const (
	RDPGFX_DVC_CHANNEL_NAME = "Microsoft::Windows::RDS::Graphics" //图形扩展
)

var StaticVirtualChannels = map[string]int{
	CLIPRDR_SVC_CHANNEL_NAME: CHANNEL_OPTION_INITIALIZED | CHANNEL_OPTION_ENCRYPT_RDP |
		CHANNEL_OPTION_COMPRESS_RDP | CHANNEL_OPTION_SHOW_PROTOCOL,
	RDPDR_SVC_CHANNEL_NAME: CHANNEL_OPTION_INITIALIZED | CHANNEL_OPTION_ENCRYPT_RDP | CHANNEL_OPTION_COMPRESS_RDP,
	RDPSND_SVC_CHANNEL_NAME: CHANNEL_OPTION_INITIALIZED | CHANNEL_OPTION_ENCRYPT_RDP |
		CHANNEL_OPTION_COMPRESS_RDP | CHANNEL_OPTION_SHOW_PROTOCOL,
	RAIL_SVC_CHANNEL_NAME: CHANNEL_OPTION_INITIALIZED | CHANNEL_OPTION_ENCRYPT_RDP |
		CHANNEL_OPTION_COMPRESS_RDP | CHANNEL_OPTION_SHOW_PROTOCOL,
}

const (
	CHANNEL_CHUNK_LENGTH       = 1600
	CHANNEL_FLAG_FIRST         = 0x01
	CHANNEL_FLAG_LAST          = 0x02
	CHANNEL_FLAG_SHOW_PROTOCOL = 0x10
)

type ChannelTransport interface {
	GetType() (string, uint32)
	Sender(core.ChannelSender)
	Process(s []byte)
}
type ChannelClient struct {
	ChannelDef
	t ChannelTransport
}

type Channels struct {
	emission.Emitter
	channels      map[string]ChannelClient
	transport     core.Transport
	buff          *bytes.Buffer
	channelSender core.ChannelSender
}

func NewChannels(t core.Transport) *Channels {
	c := &Channels{
		Emitter:   *emission.NewEmitter(),
		channels:  make(map[string]ChannelClient, 20),
		transport: t,
		buff:      &bytes.Buffer{},
	}
	t.On("channel", c.process)
	return c
}

func (c *Channels) SetChannelSender(f core.ChannelSender) {
	c.channelSender = f
}
func (c *Channels) Register(t ChannelTransport) {
	name, option := t.GetType()
	_, ok := c.channels[name]
	if ok {
		glog.Warn("Already register channel:", name)
		return
	}
	t.Sender(c)
	c.channels[name] = ChannelClient{ChannelDef{name, option}, t}
}

func (c *Channels) SendToChannel(channel string, s []byte) (int, error) {
	cli, ok := c.channels[channel]
	if !ok {
		glog.Warn("No register channel:", channel)
		return 0, fmt.Errorf("No register channel: %s", channel)
	}
	idx := 0
	ln := len(s)
	b := &bytes.Buffer{}
	for ln > 0 {
		var flag uint32 = 0
		if cli.Options&CHANNEL_OPTION_SHOW_PROTOCOL != 0 {
			flag |= CHANNEL_FLAG_SHOW_PROTOCOL
		}
		if idx == 0 {
			flag |= CHANNEL_FLAG_FIRST
		}

		var ss []byte
		if ln > CHANNEL_CHUNK_LENGTH {
			ss = s[idx : idx+CHANNEL_CHUNK_LENGTH]
			idx += CHANNEL_CHUNK_LENGTH
		} else {
			flag |= CHANNEL_FLAG_LAST
			ss = s[idx : idx+ln]
		}
		glog.Debug("len:", len(ss), "flag:", flag)
		ln -= len(ss)
		b.Reset()
		core.WriteUInt32LE(uint32(len(s)), b)
		core.WriteUInt32LE(flag, b)
		b.Write(ss)
		c.channelSender.SendToChannel(channel, b.Bytes())
	}
	return ln, nil
}

func (c *Channels) process(channel string, s []byte) {
	cli, ok := c.channels[channel]
	if !ok {
		glog.Warn("No found channel:", channel)
		return
	}
	r := bytes.NewReader(s)
	ln, _ := core.ReadUInt32LE(r)
	flags, _ := core.ReadUInt32LE(r)
	glog.Debugf("channel:%s length: %d, flags: %d", channel, ln, flags)
	if flags&CHANNEL_FLAG_FIRST == 0 || flags&CHANNEL_FLAG_LAST == 0 {
		if flags&CHANNEL_FLAG_FIRST != 0 {
			c.buff.Reset()
		}
		b, _ := core.ReadBytes(r.Len(), r)
		c.buff.Write(b)
		if flags&CHANNEL_FLAG_LAST == 0 {
			return
		}
		s = c.buff.Bytes()
	} else {
		s, _ = core.ReadBytes(r.Len(), r)
	}

	cli.t.Process(s)
}
