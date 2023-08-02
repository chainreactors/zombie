package cliprdr

import (
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"unicode/utf16"

	"github.com/lunixbochs/struc"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/plugin"
)

/**
 *                                    Initialization Sequence\n
 *     Client                                                                    Server\n
 *        |                                                                         |\n
 *        |<----------------------Server Clipboard Capabilities PDU-----------------|\n
 *        |<-----------------------------Monitor Ready PDU--------------------------|\n
 *        |-----------------------Client Clipboard Capabilities PDU---------------->|\n
 *        |---------------------------Temporary Directory PDU---------------------->|\n
 *        |-------------------------------Format List PDU-------------------------->|\n
 *        |<--------------------------Format List Response PDU----------------------|\n
 *
 */

/**
 *                                    Data Transfer Sequences\n
 *     Shared                                                                     Local\n
 *  Clipboard Owner                                                           Clipboard Owner\n
 *        |                                                                         |\n
 *        |-------------------------------------------------------------------------|\n _
 *        |-------------------------------Format List PDU-------------------------->|\n  |
 *        |<--------------------------Format List Response PDU----------------------|\n _| Copy
 * Sequence
 *        |<---------------------Lock Clipboard Data PDU (Optional)-----------------|\n
 *        |-------------------------------------------------------------------------|\n
 *        |-------------------------------------------------------------------------|\n _
 *        |<--------------------------Format Data Request PDU-----------------------|\n  | Paste
 * Sequence Palette,
 *        |---------------------------Format Data Response PDU--------------------->|\n _| Metafile,
 * File List Data
 *        |-------------------------------------------------------------------------|\n
 *        |-------------------------------------------------------------------------|\n _
 *        |<------------------------Format Contents Request PDU---------------------|\n  | Paste
 * Sequence
 *        |-------------------------Format Contents Response PDU------------------->|\n _| File
 * Stream Data
 *        |<---------------------Lock Clipboard Data PDU (Optional)-----------------|\n
 *        |-------------------------------------------------------------------------|\n
 *
 */

const (
	ChannelName   = plugin.CLIPRDR_SVC_CHANNEL_NAME
	ChannelOption = plugin.CHANNEL_OPTION_INITIALIZED | plugin.CHANNEL_OPTION_ENCRYPT_RDP |
		plugin.CHANNEL_OPTION_COMPRESS_RDP | plugin.CHANNEL_OPTION_SHOW_PROTOCOL
)

type MsgType uint16

const (
	CB_MONITOR_READY         = 0x0001
	CB_FORMAT_LIST           = 0x0002
	CB_FORMAT_LIST_RESPONSE  = 0x0003
	CB_FORMAT_DATA_REQUEST   = 0x0004
	CB_FORMAT_DATA_RESPONSE  = 0x0005
	CB_TEMP_DIRECTORY        = 0x0006
	CB_CLIP_CAPS             = 0x0007
	CB_FILECONTENTS_REQUEST  = 0x0008
	CB_FILECONTENTS_RESPONSE = 0x0009
	CB_LOCK_CLIPDATA         = 0x000A
	CB_UNLOCK_CLIPDATA       = 0x000B
)

type MsgFlags uint16

const (
	CB_RESPONSE_OK   = 0x0001
	CB_RESPONSE_FAIL = 0x0002
	CB_ASCII_NAMES   = 0x0004
)

type DwFlags uint32

const (
	FILECONTENTS_SIZE  = 0x00000001
	FILECONTENTS_RANGE = 0x00000002
)

type CliprdrPDUHeader struct {
	MsgType  uint16 `struc:"little"`
	MsgFlags uint16 `struc:"little"`
	DataLen  uint32 `struc:"little"`
}

func NewCliprdrPDUHeader(mType, flags uint16, ln uint32) *CliprdrPDUHeader {
	return &CliprdrPDUHeader{
		MsgType:  mType,
		MsgFlags: flags,
		DataLen:  ln,
	}
}
func (h *CliprdrPDUHeader) serialize() []byte {
	b := &bytes.Buffer{}
	core.WriteUInt16LE(h.MsgType, b)
	core.WriteUInt16LE(h.MsgFlags, b)
	core.WriteUInt32LE(h.DataLen, b)
	return b.Bytes()
}

type CliprdrGeneralCapabilitySet struct {
	CapabilitySetType   uint16 `struc:"little"`
	CapabilitySetLength uint16 `struc:"little"`
	Version             uint32 `struc:"little"`
	GeneralFlags        uint32 `struc:"little"`
}

const (
	CB_CAPSTYPE_GENERAL = 0x0001
)

type CliprdrCapabilitySets struct {
	CapabilitySetType uint16 `struc:"little"`
	LengthCapability  uint16 `struc:"little"`
	Version           uint32 `struc:"little"`
	GeneralFlags      uint32 `struc:"little"`
	//CapabilityData    []byte `struc:"little"`
}
type CliprdrCapabilitiesPDU struct {
	CCapabilitiesSets uint16                        `struc:"little,sizeof=CapabilitySets"`
	Pad1              uint16                        `struc:"little"`
	CapabilitySets    []CliprdrGeneralCapabilitySet `struc:"little"`
}

type CliprdrMonitorReady struct {
}

type GeneralFlags uint32

const (
	/* CLIPRDR_GENERAL_CAPABILITY.generalFlags */
	CB_USE_LONG_FORMAT_NAMES     = 0x00000002
	CB_STREAM_FILECLIP_ENABLED   = 0x00000004
	CB_FILECLIP_NO_FILE_PATHS    = 0x00000008
	CB_CAN_LOCK_CLIPDATA         = 0x00000010
	CB_HUGE_FILE_SUPPORT_ENABLED = 0x00000020
)

const (
	/* CLIPRDR_GENERAL_CAPABILITY.version */
	CB_CAPS_VERSION_1 = 0x00000001
	CB_CAPS_VERSION_2 = 0x00000002
)
const (
	CB_CAPSTYPE_GENERAL_LEN = 12
)

const (
	FD_CLSID      = 0x00000001
	FD_SIZEPOINT  = 0x00000002
	FD_ATTRIBUTES = 0x00000004
	FD_CREATETIME = 0x00000008
	FD_ACCESSTIME = 0x00000010
	FD_WRITESTIME = 0x00000020
	FD_FILESIZE   = 0x00000040
	FD_PROGRESSUI = 0x00004000
	FD_LINKUI     = 0x00008000
)

type FileGroupDescriptor struct {
	CItems uint32           `struc:"little"`
	Fgd    []FileDescriptor `struc:"sizefrom=CItems"`
}
type FileDescriptor struct {
	Flags          uint32   `struc:"little"`
	Clsid          [16]byte `struc:"little"`
	Sizel          [8]byte  `struc:"little"`
	Pointl         [8]byte  `struc:"little"`
	FileAttributes uint32   `struc:"little"`
	CreationTime   [8]byte  `struc:"little"`
	LastAccessTime [8]byte  `struc:"little"`
	LastWriteTime  []byte   `struc:"[8]byte"` //8
	FileSizeHigh   uint32   `struc:"little"`
	FileSizeLow    uint32   `struc:"little"`
	FileName       []byte   `struc:"[512]byte"`
}

func (f *FileGroupDescriptor) Unpack(b []byte) error {
	r := bytes.NewReader(b)
	err := struc.Unpack(r, f)
	if err != nil {
		glog.Error(err)
	}

	return err
}

func (f *FileDescriptor) serialize() []byte {
	b := &bytes.Buffer{}
	core.WriteUInt32LE(f.Flags, b)
	for i := 0; i < 32; i++ {
		core.WriteByte(0, b)
	}
	core.WriteUInt32LE(f.FileAttributes, b)
	for i := 0; i < 16; i++ {
		core.WriteByte(0, b)
	}
	core.WriteBytes(f.LastWriteTime[:], b)
	core.WriteUInt32LE(f.FileSizeHigh, b)
	core.WriteUInt32LE(f.FileSizeLow, b)
	name := make([]byte, 512)
	copy(name, f.FileName)
	core.WriteBytes(name, b)
	return b.Bytes()
}

func (f *FileDescriptor) isDir() bool {
	if f.Flags&FD_ATTRIBUTES != 0 {
		return f.FileAttributes&FILE_ATTRIBUTE_DIRECTORY != 0
	}

	return false
}

func (f *FileDescriptor) hasFileSize() bool {
	return f.Flags&FD_FILESIZE != 0
}

// temp dir
type CliprdrTempDirectory struct {
	SzTempDir []byte `struc:"[260]byte"`
}

// format list
type CliprdrFormat struct {
	FormatId   uint32
	FormatName string
}
type CliprdrFormatList struct {
	NumFormats uint32
	Formats    []CliprdrFormat
}
type ClipboardFormats uint16

const (
	CB_FORMAT_HTML             = 0xD010
	CB_FORMAT_PNG              = 0xD011
	CB_FORMAT_JPEG             = 0xD012
	CB_FORMAT_GIF              = 0xD013
	CB_FORMAT_TEXTURILIST      = 0xD014
	CB_FORMAT_GNOMECOPIEDFILES = 0xD015
	CB_FORMAT_MATECOPIEDFILES  = 0xD016
)

// lock or unlock
type CliprdrCtrlClipboardData struct {
	ClipDataId uint32
}

// format data
type CliprdrFormatDataRequest struct {
	RequestedFormatId uint32
}
type CliprdrFormatDataResponse struct {
	RequestedFormatData []byte
}

// file contents
type CliprdrFileContentsRequest struct {
	StreamId      uint32 `struc:"little"`
	Lindex        uint32 `struc:"little"`
	DwFlags       uint32 `struc:"little"`
	NPositionLow  uint32 `struc:"little"`
	NPositionHigh uint32 `struc:"little"`
	CbRequested   uint32 `struc:"little"`
	ClipDataId    uint32 `struc:"little"`
}

func FileContentsSizeRequest(i uint32) *CliprdrFileContentsRequest {
	return &CliprdrFileContentsRequest{
		StreamId:      1,
		Lindex:        i,
		DwFlags:       FILECONTENTS_SIZE,
		NPositionLow:  0,
		NPositionHigh: 0,
		CbRequested:   65535,
		ClipDataId:    0,
	}
}

type CliprdrFileContentsResponse struct {
	StreamId      uint32
	CbRequested   uint32
	RequestedData []byte
}

func (resp *CliprdrFileContentsResponse) Unpack(b []byte) {
	r := bytes.NewReader(b)
	resp.StreamId, _ = core.ReadUInt32LE(r)
	resp.CbRequested = uint32(r.Len())
	resp.RequestedData, _ = core.ReadBytes(int(resp.CbRequested), r)
}

type CliprdrClient struct {
	w                     core.ChannelSender
	useLongFormatNames    bool
	streamFileClipEnabled bool
	fileClipNoFilePaths   bool
	canLockClipData       bool
	hasHugeFileSupport    bool
	formatIdMap           map[uint32]uint32
	Files                 []FileDescriptor
	reply                 chan []byte
	Control
}

func NewCliprdrClient() *CliprdrClient {
	c := &CliprdrClient{
		formatIdMap: make(map[uint32]uint32, 20),
		Files:       make([]FileDescriptor, 0, 20),
		reply:       make(chan []byte, 100),
	}

	go ClipWatcher(c)

	return c
}

func (c *CliprdrClient) Send(s []byte) (int, error) {
	glog.Debug("len:", len(s), "data:", hex.EncodeToString(s))
	name, _ := c.GetType()
	return c.w.SendToChannel(name, s)
}
func (c *CliprdrClient) Sender(f core.ChannelSender) {
	c.w = f
}
func (c *CliprdrClient) GetType() (string, uint32) {
	return ChannelName, ChannelOption
}

func (c *CliprdrClient) Process(s []byte) {
	glog.Debug("recv:", hex.EncodeToString(s))
	r := bytes.NewReader(s)

	msgType, _ := core.ReadUint16LE(r)
	flag, _ := core.ReadUint16LE(r)
	length, _ := core.ReadUInt32LE(r)
	glog.Debugf("cliprdr: type=0x%x flag=%d length=%d, all=%d", msgType, flag, length, r.Len())

	b, _ := core.ReadBytes(int(length), r)

	switch msgType {
	case CB_CLIP_CAPS:
		glog.Info("CB_CLIP_CAPS")
		c.processClipCaps(b)

	case CB_MONITOR_READY:
		glog.Info("CB_MONITOR_READY")
		c.processMonitorReady(b)

	case CB_FORMAT_LIST:
		glog.Info("CB_FORMAT_LIST")
		c.processFormatList(b)

	case CB_FORMAT_LIST_RESPONSE:
		glog.Info("CB_FORMAT_LIST_RESPONSE")
		c.processFormatListResponse(flag, b)

	case CB_FORMAT_DATA_REQUEST:
		glog.Info("CB_FORMAT_DATA_REQUEST")
		c.processFormatDataRequest(b)

	case CB_FORMAT_DATA_RESPONSE:
		glog.Info("CB_FORMAT_DATA_RESPONSE")
		c.processFormatDataResponse(flag, b)

	case CB_FILECONTENTS_REQUEST:
		glog.Info("CB_FILECONTENTS_REQUEST")
		c.processFileContentsRequest(b)

	case CB_FILECONTENTS_RESPONSE:
		glog.Info("CB_FILECONTENTS_RESPONSE")
		c.processFileContentsResponse(flag, b)

	case CB_LOCK_CLIPDATA:
		glog.Info("CB_LOCK_CLIPDATA")
		c.processLockClipData(b)

	case CB_UNLOCK_CLIPDATA:
		glog.Info("CB_UNLOCK_CLIPDATA")
		c.processUnlockClipData(b)

	default:
		glog.Errorf("type 0x%x not supported", msgType)
	}
}
func (c *CliprdrClient) processClipCaps(b []byte) {
	r := bytes.NewReader(b)
	var cp CliprdrCapabilitiesPDU
	err := struc.Unpack(r, &cp)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Debugf("Capabilities:%+v", cp)
	c.useLongFormatNames = cp.CapabilitySets[0].GeneralFlags&CB_USE_LONG_FORMAT_NAMES != 0
	c.streamFileClipEnabled = cp.CapabilitySets[0].GeneralFlags&CB_STREAM_FILECLIP_ENABLED != 0
	c.fileClipNoFilePaths = cp.CapabilitySets[0].GeneralFlags&CB_FILECLIP_NO_FILE_PATHS != 0
	c.canLockClipData = cp.CapabilitySets[0].GeneralFlags&CB_CAN_LOCK_CLIPDATA != 0
	c.hasHugeFileSupport = cp.CapabilitySets[0].GeneralFlags&CB_HUGE_FILE_SUPPORT_ENABLED != 0
	glog.Info("UseLongFormatNames:", c.useLongFormatNames)
	glog.Info("StreamFileClipEnabled:", c.streamFileClipEnabled)
	glog.Info("FileClipNoFilePaths:", c.fileClipNoFilePaths)
	glog.Info("CanLockClipData:", c.canLockClipData)
	glog.Info("HasHugeFileSupport:", c.hasHugeFileSupport)
}

func (c *CliprdrClient) processMonitorReady(b []byte) {
	//Client Clipboard Capabilities PDU
	c.sendClientCapabilitiesPDU()

	//Temporary Directory PDU
	//c.sendTemporaryDirectoryPDU()

	//Format List PDU
	c.sendFormatListPDU()

}
func (c *CliprdrClient) processFormatList(b []byte) {
	c.withOpenClipboard(func() {
		if !EmptyClipboard() {
			glog.Error("EmptyClipboard failed")
		}
	})
	fl, hasFile := c.readForamtList(b)
	glog.Info("numFormats:", fl.NumFormats)

	if hasFile {
		c.SendCliprdrMessage()
	} else {
		c.withOpenClipboard(func() {
			if !EmptyClipboard() {
				glog.Error("EmptyClipboard failed")
			}
			for i := range c.formatIdMap {
				glog.Debug("i:", i)
				SetClipboardData(i, 0)
			}
		})

	}

	c.sendFormatListResponse(CB_RESPONSE_OK)
}
func (c *CliprdrClient) processFormatListResponse(flag uint16, b []byte) {
	if flag != CB_RESPONSE_OK {
		glog.Error("Format List Response Failed")
		return
	}
	glog.Error("Format List Response OK")
}
func getFilesDescriptor(name string) (FileDescriptor, error) {
	var fd FileDescriptor
	fd.Flags = FD_ATTRIBUTES | FD_FILESIZE | FD_WRITESTIME | FD_PROGRESSUI
	f, e := os.Stat(name)
	if e != nil {
		glog.Error(e.Error())
		return fd, e
	}
	fd.FileAttributes, fd.LastWriteTime,
		fd.FileSizeHigh, fd.FileSizeLow = GetFileInfo(f.Sys())
	fd.FileName = core.UnicodeEncode(name)

	return fd, nil
}
func (c *CliprdrClient) processFormatDataRequest(b []byte) {
	r := bytes.NewReader(b)
	requestId, _ := core.ReadUInt32LE(r)

	buff := &bytes.Buffer{}
	if requestId == RegisterClipboardFormat(CFSTR_FILEDESCRIPTORW) {
		fs := GetFileNames()
		core.WriteUInt32LE(uint32(len(fs)), buff)
		c.Files = c.Files[:0]
		for _, v := range fs {
			glog.Info("Name:", v)
			f, _ := getFilesDescriptor(v)
			buff.Write(f.serialize())
			for i := 0; i < 8; i++ {
				buff.WriteByte(0)
			}
			c.Files = append(c.Files, f)

		}
	} else {
		c.withOpenClipboard(func() {
			data := GetClipboardData(requestId)
			glog.Debug("data:", data)
			buff.Write(core.UnicodeEncode(data))
			buff.Write([]byte{0, 0})
		})
	}

	c.sendFormatDataResponse(buff.Bytes())
}
func (c *CliprdrClient) processFormatDataResponse(flag uint16, b []byte) {
	if flag != CB_RESPONSE_OK {
		glog.Error("Format Data Response Failed")
	}
	c.reply <- b
}

func (c *CliprdrClient) processFileContentsRequest(b []byte) {
	r := bytes.NewReader(b)
	var req CliprdrFileContentsRequest
	struc.Unpack(r, &req)
	if len(c.Files) <= int(req.Lindex) {
		glog.Error("No found file:", req.Lindex)
		c.sendFormatContentsResponse(req.StreamId, []byte{})
		return
	}
	buff := &bytes.Buffer{}
	/*o := OleGetClipboard()
	var format_etc FORMATETC
	var stg_medium STGMEDIUM
	format_etc.CFormat = RegisterClipboardFormat(CFSTR_FILECONTENTS)
	format_etc.Tymed = TYMED_ISTREAM
	format_etc.Aspect = 1
	format_etc.Index = req.Lindex
	o.GetData(&format_etc, &stg_medium)
	s, _ := stg_medium.Stream()*/
	f := c.Files[req.Lindex]
	if req.DwFlags == FILECONTENTS_SIZE {
		core.WriteUInt32LE(f.FileSizeLow, buff)
		core.WriteUInt32LE(f.FileSizeHigh, buff)
		c.sendFormatContentsResponse(req.StreamId, buff.Bytes())
	} else if req.DwFlags == FILECONTENTS_RANGE {
		name := core.UnicodeDecode(f.FileName)
		fi, err := os.Open(name)
		if err != nil {
			glog.Error(err.Error())
			return
		}
		defer fi.Close()
		data := make([]byte, req.CbRequested)
		n, _ := fi.ReadAt(data, int64(f.FileSizeHigh))
		c.sendFormatContentsResponse(req.StreamId, data[:n])
	}
}
func (c *CliprdrClient) processFileContentsResponse(flag uint16, b []byte) {
	if flag != CB_RESPONSE_OK {
		glog.Error("File Contents Response Failed")
	}
	var resp CliprdrFileContentsResponse
	resp.Unpack(b)
	glog.Debug("Get File Contents Response:", resp.StreamId, resp.CbRequested)
	c.reply <- resp.RequestedData
}
func (c *CliprdrClient) processLockClipData(b []byte) {
	r := bytes.NewReader(b)
	var l CliprdrCtrlClipboardData
	l.ClipDataId, _ = core.ReadUInt32LE(r)
}
func (c *CliprdrClient) processUnlockClipData(b []byte) {
	r := bytes.NewReader(b)
	var l CliprdrCtrlClipboardData
	l.ClipDataId, _ = core.ReadUInt32LE(r)

}

func (c *CliprdrClient) sendClientCapabilitiesPDU() {
	glog.Info("Send Client Clipboard Capabilities PDU")
	var cs CliprdrGeneralCapabilitySet
	cs.CapabilitySetLength = 12
	cs.CapabilitySetType = CB_CAPSTYPE_GENERAL
	cs.Version = CB_CAPS_VERSION_2
	cs.GeneralFlags = CB_USE_LONG_FORMAT_NAMES |
		CB_STREAM_FILECLIP_ENABLED |
		CB_FILECLIP_NO_FILE_PATHS
	var cc CliprdrCapabilitiesPDU
	cc.CCapabilitiesSets = 1
	cc.Pad1 = 0
	cc.CapabilitySets = make([]CliprdrGeneralCapabilitySet, 0, 1)
	cc.CapabilitySets = append(cc.CapabilitySets, cs)
	header := NewCliprdrPDUHeader(CB_CLIP_CAPS, 0, 16)

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt16LE(cc.CCapabilitiesSets, buff)
	core.WriteUInt16LE(cc.Pad1, buff)
	for _, v := range cc.CapabilitySets {
		struc.Pack(buff, v)
	}

	c.Send(buff.Bytes())
}

func (c *CliprdrClient) sendTemporaryDirectoryPDU() {
	glog.Info("Send Temporary Directory PDU")
	var t CliprdrTempDirectory
	header := &CliprdrPDUHeader{CB_TEMP_DIRECTORY, 0, 260}
	t.SzTempDir = core.UnicodeEncode(os.TempDir())

	buff := &bytes.Buffer{}
	core.WriteBytes(header.serialize(), buff)
	core.WriteBytes(t.SzTempDir, buff)
	c.Send(buff.Bytes())
}
func (c *CliprdrClient) sendFormatListPDU() {
	glog.Info("Send Format List PDU")
	var f CliprdrFormatList

	f.Formats = GetFormatList(c.hwnd)
	f.NumFormats = uint32(len(f.Formats))

	glog.Info("NumFormats:", f.NumFormats)
	glog.Debug("Formats:", f.Formats)

	b := &bytes.Buffer{}
	for _, v := range f.Formats {
		core.WriteUInt32LE(v.FormatId, b)
		if v.FormatName == "" {
			core.WriteUInt16LE(0, b)
		} else {
			n := core.UnicodeEncode(v.FormatName)
			core.WriteBytes(n, b)
			b.Write([]byte{0, 0})
		}
	}

	header := NewCliprdrPDUHeader(CB_FORMAT_LIST, 0, uint32(b.Len()))

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteBytes(b.Bytes(), buff)

	c.Send(buff.Bytes())
}
func (c *CliprdrClient) readForamtList(b []byte) (*CliprdrFormatList, bool) {
	r := bytes.NewReader(b)
	fs := make([]CliprdrFormat, 0, 20)
	var numFormats uint32 = 0
	hasFile := false
	c.formatIdMap = make(map[uint32]uint32, 0)
	for r.Len() > 0 {
		foramtId, _ := core.ReadUInt32LE(r)
		bs := make([]uint16, 0, 20)
		ln := r.Len()
		for j := 0; j < ln; j++ {
			b, _ := core.ReadUint16LE(r)
			if b == 0 {
				break
			}
			bs = append(bs, b)
		}
		name := string(utf16.Decode(bs))
		if strings.EqualFold(name, CFSTR_FILEDESCRIPTORW) {
			hasFile = true
		}
		glog.Infof("Foramt:%d Name:<%s>", foramtId, name)
		if name != "" {
			localId := RegisterClipboardFormat(name)
			glog.Info("local:", localId, "remote:", foramtId)
			c.formatIdMap[localId] = foramtId
		} else {
			c.formatIdMap[foramtId] = foramtId
		}

		numFormats++
		fs = append(fs, CliprdrFormat{foramtId, name})
	}

	return &CliprdrFormatList{numFormats, fs}, hasFile
}

func (c *CliprdrClient) sendFormatListResponse(flags uint16) {
	glog.Info("Send Format List Response")
	header := NewCliprdrPDUHeader(CB_FORMAT_LIST_RESPONSE, flags, 0)
	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	c.Send(buff.Bytes())
}

func (c *CliprdrClient) sendFormatDataRequest(id uint32) {
	glog.Info("Send Format Data Request")
	var r CliprdrFormatDataRequest
	r.RequestedFormatId = id
	header := NewCliprdrPDUHeader(CB_FORMAT_DATA_REQUEST, 0, 4)

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt32LE(r.RequestedFormatId, buff)

	c.Send(buff.Bytes())
}
func (c *CliprdrClient) sendFormatDataResponse(b []byte) {
	glog.Info("Send Format Data Response")
	var resp CliprdrFormatDataResponse
	resp.RequestedFormatData = b

	header := NewCliprdrPDUHeader(CB_FORMAT_DATA_RESPONSE, CB_RESPONSE_OK, uint32(len(resp.RequestedFormatData)))

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	buff.Write(resp.RequestedFormatData)

	c.Send(buff.Bytes())
}

func (c *CliprdrClient) sendFormatContentsRequest(r CliprdrFileContentsRequest) uint32 {
	glog.Info("Send Format Contents Request")
	glog.Debugf("Format Contents Request:%+v", r)
	header := NewCliprdrPDUHeader(CB_FILECONTENTS_REQUEST, 0, 28)

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt32LE(r.StreamId, buff)
	core.WriteUInt32LE(uint32(r.Lindex), buff)
	core.WriteUInt32LE(r.DwFlags, buff)
	core.WriteUInt32LE(r.NPositionLow, buff)
	core.WriteUInt32LE(r.NPositionHigh, buff)
	core.WriteUInt32LE(r.CbRequested, buff)
	core.WriteUInt32LE(r.ClipDataId, buff)

	c.Send(buff.Bytes())

	return uint32(buff.Len())
}
func (c *CliprdrClient) sendFormatContentsResponse(streamId uint32, b []byte) {
	glog.Info("Send Format Contents Response")
	var r CliprdrFileContentsResponse
	r.StreamId = streamId
	r.RequestedData = b
	r.CbRequested = uint32(len(b))
	header := NewCliprdrPDUHeader(CB_FILECONTENTS_RESPONSE, CB_RESPONSE_OK, uint32(4+r.CbRequested))

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt32LE(r.StreamId, buff)
	core.WriteBytes(r.RequestedData, buff)

	c.Send(buff.Bytes())
}

func (c *CliprdrClient) sendLockClipData() {
	glog.Info("Send Lock Clip Data")
	var r CliprdrCtrlClipboardData
	header := NewCliprdrPDUHeader(CB_LOCK_CLIPDATA, 0, 4)

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt32LE(r.ClipDataId, buff)

	c.Send(buff.Bytes())
}

func (c *CliprdrClient) sendUnlockClipData() {
	glog.Info("Send Unlock Clip Data")
	var r CliprdrCtrlClipboardData
	header := NewCliprdrPDUHeader(CB_UNLOCK_CLIPDATA, 0, 4)

	buff := &bytes.Buffer{}
	buff.Write(header.serialize())
	core.WriteUInt32LE(r.ClipDataId, buff)

	c.Send(buff.Bytes())
}
