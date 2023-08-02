// rail.go
package rail

import (
	"bytes"
	"encoding/hex"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/plugin"
)

const (
	ChannelName   = plugin.RAIL_SVC_CHANNEL_NAME
	ChannelOption = plugin.CHANNEL_OPTION_INITIALIZED | plugin.CHANNEL_OPTION_ENCRYPT_RDP |
		plugin.CHANNEL_OPTION_COMPRESS_RDP | plugin.CHANNEL_OPTION_SHOW_PROTOCOL
)

const (
	TS_RAIL_ORDER_EXEC                  = 0x0001
	TS_RAIL_ORDER_ACTIVATE              = 0x0002
	TS_RAIL_ORDER_SYSPARAM              = 0x0003
	TS_RAIL_ORDER_SYSCOMMAND            = 0x0004
	TS_RAIL_ORDER_HANDSHAKE             = 0x0005
	TS_RAIL_ORDER_NOTIFY_EVENT          = 0x0006
	TS_RAIL_ORDER_WINDOWMOVE            = 0x0008
	TS_RAIL_ORDER_LOCALMOVESIZE         = 0x0009
	TS_RAIL_ORDER_MINMAXINFO            = 0x000A
	TS_RAIL_ORDER_CLIENTSTATUS          = 0x000B
	TS_RAIL_ORDER_SYSMENU               = 0x000C
	TS_RAIL_ORDER_LANGBARINFO           = 0x000D
	TS_RAIL_ORDER_GET_APPID_REQ         = 0x000E
	TS_RAIL_ORDER_GET_APPID_RESP        = 0x000F
	TS_RAIL_ORDER_TASKBARINFO           = 0x0010
	TS_RAIL_ORDER_LANGUAGEIMEINFO       = 0x0011
	TS_RAIL_ORDER_COMPARTMENTINFO       = 0x0012
	TS_RAIL_ORDER_HANDSHAKE_EX          = 0x0013
	TS_RAIL_ORDER_ZORDER_SYNC           = 0x0014
	TS_RAIL_ORDER_CLOAK                 = 0x0015
	TS_RAIL_ORDER_POWER_DISPLAY_REQUEST = 0x0016
	TS_RAIL_ORDER_SNAP_ARRANGE          = 0x0017
	TS_RAIL_ORDER_GET_APPID_RESP_EX     = 0x0018
	TS_RAIL_ORDER_EXEC_RESULT           = 0x0080
)

type RailClient struct {
	w                        core.ChannelSender
	DesktopWidth             uint16
	DesktopHeight            uint16
	RemoteApplicationProgram string
	ShellWorkingDirectory    string
	RemoteApplicationCmdLine string
}

func NewClient() *RailClient {
	return &RailClient{
		DesktopWidth:             800,
		DesktopHeight:            600,
		RemoteApplicationProgram: "calc",
		ShellWorkingDirectory:    "/tmp",
	}
}

type RailPDUHeader struct {
	OrderType   uint16 `struc:"little"`
	OrderLength uint16 `struc:"little"`
}

func NewRailPDUHeader(mType, ln uint16) *RailPDUHeader {
	return &RailPDUHeader{
		OrderType:   mType,
		OrderLength: ln,
	}
}

func (h *RailPDUHeader) serialize() []byte {
	b := &bytes.Buffer{}
	core.WriteUInt16LE(h.OrderType, b)
	core.WriteUInt16LE(h.OrderLength, b)
	return b.Bytes()
}

func (c *RailClient) sendData(mType uint16, ln int, s []byte) {
	glog.Debug(ln, ":ln:", len(s), "data:", hex.EncodeToString(s))
	header := NewRailPDUHeader(mType, uint16(ln))

	b := &bytes.Buffer{}
	core.WriteBytes(header.serialize(), b)
	core.WriteBytes(s, b)

	c.Send(b.Bytes())
}

func (c *RailClient) Send(s []byte) (int, error) {
	glog.Debug("len:", len(s), "data:", hex.EncodeToString(s))
	name, _ := c.GetType()
	return c.w.SendToChannel(name, s)
}
func (c *RailClient) Sender(f core.ChannelSender) {
	c.w = f
}
func (c *RailClient) GetType() (string, uint32) {
	return ChannelName, ChannelOption
}

func (c *RailClient) Process(s []byte) {
	glog.Debug("recv:", hex.EncodeToString(s))
	r := bytes.NewReader(s)
	msgType, _ := core.ReadUint16LE(r)
	length, _ := core.ReadUint16LE(r)

	glog.Infof("rail: type=0x%x length=%d, all=%d", msgType, length, r.Len())

	b, _ := core.ReadBytes(int(length), r)
	glog.Info("b:", hex.EncodeToString(b))

	switch msgType {
	case TS_RAIL_ORDER_HANDSHAKE:
		glog.Info("TS_RAIL_ORDER_HANDSHAKE")
		c.processOrderHandshake(b)
	case TS_RAIL_ORDER_SYSPARAM:
		glog.Info("TS_RAIL_ORDER_SYSPARAM")
		c.processOrderSysparam(b)
	case TS_RAIL_ORDER_EXEC_RESULT:
		glog.Info("TS_RAIL_ORDER_EXEC_RESULT")
		c.processExecResult(b)

	default:
		glog.Errorf("type 0x%x not supported", msgType)
	}
}

func (c *RailClient) processOrderHandshake(b []byte) {
	r := bytes.NewReader(b)
	buildNumber, _ := core.ReadUInt32LE(r)
	glog.Info("buildNumber:", buildNumber)

	//send client info
	c.sendClientStatus()

	//send client systemparam
	c.sendClientSystemparam()

	//send client execute
	c.sendClientExecute()
}

const (
	TS_RAIL_CLIENTSTATUS_ALLOWLOCALMOVESIZE              = 0x00000001
	TS_RAIL_CLIENTSTATUS_AUTORECONNECT                   = 0x00000002
	TS_RAIL_CLIENTSTATUS_ZORDER_SYNC                     = 0x00000004
	TS_RAIL_CLIENTSTATUS_WINDOW_RESIZE_MARGIN_SUPPORTED  = 0x00000010
	TS_RAIL_CLIENTSTATUS_HIGH_DPI_ICONS_SUPPORTED        = 0x00000020
	TS_RAIL_CLIENTSTATUS_APPBAR_REMOTING_SUPPORTED       = 0x00000040
	TS_RAIL_CLIENTSTATUS_POWER_DISPLAY_REQUEST_SUPPORTED = 0x00000080
	TS_RAIL_CLIENTSTATUS_GET_APPID_RESPONSE_EX_SUPPORTED = 0x00000100
	TS_RAIL_CLIENTSTATUS_BIDIRECTIONAL_CLOAK_SUPPORTED   = 0x00000200
)

func (c *RailClient) sendClientStatus() {
	glog.Info("Send client Status")
	var flags uint32 = TS_RAIL_CLIENTSTATUS_ALLOWLOCALMOVESIZE

	//if (settings->AutoReconnectionEnabled)
	//clientStatus.flags |= TS_RAIL_CLIENTSTATUS_AUTORECONNECT;

	flags |= TS_RAIL_CLIENTSTATUS_ZORDER_SYNC
	flags |= TS_RAIL_CLIENTSTATUS_WINDOW_RESIZE_MARGIN_SUPPORTED
	flags |= TS_RAIL_CLIENTSTATUS_APPBAR_REMOTING_SUPPORTED
	flags |= TS_RAIL_CLIENTSTATUS_POWER_DISPLAY_REQUEST_SUPPORTED
	flags |= TS_RAIL_CLIENTSTATUS_BIDIRECTIONAL_CLOAK_SUPPORTED

	b := &bytes.Buffer{}
	core.WriteUInt32LE(flags, b)

	c.sendData(TS_RAIL_ORDER_CLIENTSTATUS, 4, b.Bytes())
}

const (
	SPI_SET_SCREEN_SAVE_ACTIVE = 0x00000011
	SPI_SET_SCREEN_SAVE_SECURE = 0x00000077
)
const (
	/*Bit mask values for SPI_ parameters*/
	SPI_MASK_SET_DRAG_FULL_WINDOWS      = 0x00000001
	SPI_MASK_SET_KEYBOARD_CUES          = 0x00000002
	SPI_MASK_SET_KEYBOARD_PREF          = 0x00000004
	SPI_MASK_SET_MOUSE_BUTTON_SWAP      = 0x00000008
	SPI_MASK_SET_WORK_AREA              = 0x00000010
	SPI_MASK_DISPLAY_CHANGE             = 0x00000020
	SPI_MASK_TASKBAR_POS                = 0x00000040
	SPI_MASK_SET_HIGH_CONTRAST          = 0x00000080
	SPI_MASK_SET_SCREEN_SAVE_ACTIVE     = 0x00000100
	SPI_MASK_SET_SET_SCREEN_SAVE_SECURE = 0x00000200
	SPI_MASK_SET_CARET_WIDTH            = 0x00000400
	SPI_MASK_SET_STICKY_KEYS            = 0x00000800
	SPI_MASK_SET_TOGGLE_KEYS            = 0x00001000
	SPI_MASK_SET_FILTER_KEYS            = 0x00002000
)
const (
	SPI_SET_DRAG_FULL_WINDOWS = 0x00000025
	SPI_SET_KEYBOARD_CUES     = 0x0000100B
	SPI_SET_KEYBOARD_PREF     = 0x00000045
	SPI_SET_MOUSE_BUTTON_SWAP = 0x00000021
	SPI_SET_WORK_AREA         = 0x0000002F
	SPI_DISPLAY_CHANGE        = 0x0000F001
	SPI_TASKBAR_POS           = 0x0000F000
	SPI_SET_HIGH_CONTRAST     = 0x00000043
	SPI_SETCARETWIDTH         = 0x00002007
	SPI_SETSTICKYKEYS         = 0x0000003B
	SPI_SETTOGGLEKEYS         = 0x00000035
	SPI_SETFILTERKEYS         = 0x00000033
)

type TsFilterKeys struct {
	Flags      uint32
	WaitTime   uint32
	DelayTime  uint32
	RepeatTime uint32
	BounceTime uint32
}
type RailHighContrast struct {
	flags             uint32
	colorSchemeLength uint32
	colorScheme       string
}
type Rectangle16 struct {
	left   uint16
	top    uint16
	right  uint16
	bottom uint16
}
type RailSysparamOrder struct {
	param               uint32
	params              uint32
	dragFullWindows     uint8
	keyboardCues        uint8
	keyboardPref        uint8
	mouseButtonSwap     uint8
	workArea            Rectangle16
	displayChange       Rectangle16
	taskbarPos          Rectangle16
	highContrast        RailHighContrast
	caretWidth          uint32
	stickyKeys          uint32
	toggleKeys          uint32
	filterKeys          TsFilterKeys
	setScreenSaveActive uint8
	setScreenSaveSecure uint8
}

func (c *RailClient) sendClientSystemparam() {
	glog.Info("Send client Systemparam")

	var sp RailSysparamOrder
	sp.params = 0
	sp.params |= SPI_MASK_SET_HIGH_CONTRAST
	sp.highContrast.colorScheme = ""
	sp.highContrast.colorSchemeLength = 0
	sp.highContrast.flags = 0x7E
	sp.params |= SPI_MASK_SET_MOUSE_BUTTON_SWAP
	sp.mouseButtonSwap = 0
	sp.params |= SPI_MASK_SET_KEYBOARD_PREF
	sp.keyboardPref = 0
	sp.params |= SPI_MASK_SET_DRAG_FULL_WINDOWS
	sp.dragFullWindows = 0
	sp.params |= SPI_MASK_SET_KEYBOARD_CUES
	sp.keyboardCues = 0
	sp.params |= SPI_MASK_SET_WORK_AREA
	sp.workArea.left = 0
	sp.workArea.top = 0
	sp.workArea.right = c.DesktopWidth
	sp.workArea.bottom = c.DesktopHeight

	if sp.params&SPI_MASK_SET_HIGH_CONTRAST != 0 {
		sp.param = SPI_SET_HIGH_CONTRAST
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_TASKBAR_POS != 0 {
		sp.param = SPI_TASKBAR_POS
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_SET_MOUSE_BUTTON_SWAP != 0 {
		sp.param = SPI_SET_MOUSE_BUTTON_SWAP
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_SET_KEYBOARD_PREF != 0 {
		sp.param = SPI_SET_KEYBOARD_PREF
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_SET_DRAG_FULL_WINDOWS != 0 {
		sp.param = SPI_SET_DRAG_FULL_WINDOWS
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_SET_KEYBOARD_CUES != 0 {
		sp.param = SPI_SET_KEYBOARD_CUES
		c.sendOneClientSysparam(&sp)
	}

	if sp.params&SPI_MASK_SET_WORK_AREA != 0 {
		sp.param = SPI_SET_WORK_AREA
		glog.Debug("SPI_SET_WORK_AREA")
		c.sendOneClientSysparam(&sp)
	}
}

func (c *RailClient) sendOneClientSysparam(sp *RailSysparamOrder) {
	length := 0
	b := &bytes.Buffer{}
	core.WriteUInt32LE(sp.param, b)
	switch sp.param {
	case SPI_SET_DRAG_FULL_WINDOWS:
		core.WriteUInt8(sp.dragFullWindows, b)

	case SPI_SET_KEYBOARD_CUES:
		core.WriteUInt8(sp.keyboardCues, b)

	case SPI_SET_KEYBOARD_PREF:
		core.WriteUInt8(sp.keyboardPref, b)

	case SPI_SET_MOUSE_BUTTON_SWAP:
		core.WriteUInt8(sp.mouseButtonSwap, b)

	case SPI_SET_WORK_AREA:
		core.WriteUInt16LE(sp.workArea.left, b)
		core.WriteUInt16LE(sp.workArea.top, b)
		core.WriteUInt16LE(sp.workArea.right, b)
		core.WriteUInt16LE(sp.workArea.bottom, b)

	case SPI_DISPLAY_CHANGE:
		core.WriteUInt16LE(sp.displayChange.left, b)
		core.WriteUInt16LE(sp.displayChange.top, b)
		core.WriteUInt16LE(sp.displayChange.right, b)
		core.WriteUInt16LE(sp.displayChange.bottom, b)

	case SPI_TASKBAR_POS:
		core.WriteUInt16LE(sp.taskbarPos.left, b)
		core.WriteUInt16LE(sp.taskbarPos.top, b)
		core.WriteUInt16LE(sp.taskbarPos.right, b)
		core.WriteUInt16LE(sp.taskbarPos.bottom, b)

	case SPI_SET_HIGH_CONTRAST:
		core.WriteUInt32LE(sp.highContrast.flags, b)
		core.WriteUInt32LE(sp.highContrast.colorSchemeLength, b)
		data := core.UnicodeEncode(sp.highContrast.colorScheme)
		core.WriteBytes(data, b)

	case SPI_SETFILTERKEYS:
		core.WriteUInt32LE(sp.filterKeys.Flags, b)
		core.WriteUInt32LE(sp.filterKeys.WaitTime, b)
		core.WriteUInt32LE(sp.filterKeys.DelayTime, b)
		core.WriteUInt32LE(sp.filterKeys.RepeatTime, b)
		core.WriteUInt32LE(sp.filterKeys.BounceTime, b)

	case SPI_SETSTICKYKEYS:
		core.WriteUInt32LE(sp.stickyKeys, b)

	case SPI_SETCARETWIDTH:
		core.WriteUInt32LE(sp.caretWidth, b)

	case SPI_SETTOGGLEKEYS:
		core.WriteUInt32LE(sp.toggleKeys, b)

	case SPI_MASK_SET_SET_SCREEN_SAVE_SECURE:
		core.WriteUInt8(sp.setScreenSaveSecure, b)

	case SPI_MASK_SET_SCREEN_SAVE_ACTIVE:
		core.WriteUInt8(sp.setScreenSaveActive, b)

	default:
		glog.Error("ERROR_BAD_ARGUMENTS")
		return
	}

	c.sendData(TS_RAIL_ORDER_SYSPARAM, length+b.Len(), b.Bytes())
}

type RailExecOrder struct {
	flags                       uint16
	RemoteApplicationProgram    string
	RemoteApplicationWorkingDir string
	RemoteApplicationArguments  string
}

func (c *RailClient) sendClientExecute() {
	glog.Info("Send Client Execute")
	var exec RailExecOrder
	//exec.flags = TS_RAIL_EXEC_FLAG_EXPAND_ARGUMENTS
	exec.RemoteApplicationProgram = c.RemoteApplicationProgram
	exec.RemoteApplicationWorkingDir = c.ShellWorkingDirectory
	exec.RemoteApplicationArguments = c.RemoteApplicationCmdLine

	program := core.UnicodeEncode(exec.RemoteApplicationProgram)
	workdir := core.UnicodeEncode(exec.RemoteApplicationWorkingDir)
	arguments := core.UnicodeEncode(exec.RemoteApplicationArguments)

	length := 4
	b := &bytes.Buffer{}
	core.WriteUInt16LE(exec.flags, b)
	core.WriteUInt16LE(uint16(len(program)), b)
	core.WriteUInt16LE(uint16(len(workdir)), b)
	core.WriteUInt16LE(uint16(len(arguments)), b)
	core.WriteBytes(program, b)
	core.WriteBytes(workdir, b)
	core.WriteBytes(arguments, b)
	length += b.Len()

	c.sendData(TS_RAIL_ORDER_EXEC, length, b.Bytes())

}

func (c *RailClient) processOrderSysparam(b []byte) {
	r := bytes.NewReader(b)
	systemParam, _ := core.ReadUInt32LE(r)
	body, _ := core.ReadUInt8(r)
	glog.Infof("systemParam:0x%x, body:%d", systemParam, body)
}

const (
	//The Client Execute request was successful and the requested application or file has been launched.
	RAIL_EXEC_S_OK = 0x0000
	//The Client Execute request could not be satisfied because the server is not monitoring the current input desktop.
	RAIL_EXEC_E_HOOK_NOT_LOADED = 0x0001
	//The Execute request could not be satisfied because the request PDU was malformed.
	RAIL_EXEC_E_DECODE_FAILED = 0x0002
	//The Client Execute request could not be satisfied because the requested application was blocked by policy from being launched on the server.
	RAIL_EXEC_E_NOT_IN_ALLOWLIST = 0x0003
	//The Client Execute request could not be satisfied because the application or file path could not be found.
	RAIL_EXEC_E_FILE_NOT_FOUND = 0x0005
	//The Client Execute request could not be satisfied because an unspecified error occurred on the server.
	RAIL_EXEC_E_FAIL = 0x0006
	//The Client Execute request could not be satisfied because the remote session is locked.
	RAIL_EXEC_E_SESSION_LOCKED = 0x0007
)

func (c *RailClient) processExecResult(b []byte) {
	r := bytes.NewReader(b)
	flags, _ := core.ReadUint16LE(r)
	execResult, _ := core.ReadUint16LE(r)
	rawResult, _ := core.ReadUInt32LE(r)
	core.ReadUint16LE(r)
	exeOrFileLength, _ := core.ReadUint16LE(r)
	exeOrFile, _ := core.ReadBytes(r.Len(), r)
	glog.Info("flags:", flags, "execResult:", execResult, "rawResult:", rawResult)
	glog.Info("length:", exeOrFileLength, "file:", core.UnicodeDecode(exeOrFile))
}
