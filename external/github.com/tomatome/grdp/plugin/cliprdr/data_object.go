// +build windows

// dataobject.go
package cliprdr

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/win"
)

const (
	S_OK                 = 0x00000000
	E_UNEXPECTED         = 0x8000FFFF
	E_NOTIMPL            = 0x80004001
	E_OUTOFMEMORY        = 0x8007000E
	E_INVALIDARG         = 0x80070057
	E_NOINTERFACE        = 0x80004002
	E_POINTER            = 0x80004003
	E_HANDLE             = 0x80070006
	E_ABORT              = 0x80004004
	E_FAIL               = 0x80004005
	E_ACCESSDENIED       = 0x80070005
	E_PENDING            = 0x8000000A
	E_ADVISENOTSUPPORTED = 0x80040003
	E_FORMATETC          = 0x80040064
	CO_E_CLASSSTRING     = 0x800401F3
)

type HRESULT uintptr

func (hr HRESULT) Error() string {
	switch uint32(hr) {
	case 0x80040064:
		return "DV_E_FORMATETC (0x80040064)"
	case 0x800401D3:
		return "CLIPBRD_E_BAD_DATA (0x800401D3)"
	case 0x80004005:
		return "E_FAIL (0x80004005)"
	case 0x00000001:
		return "S_FALSE (0x00000001)"
	}
	return fmt.Sprintf("%d", hr)
}

const (
	TYMED_NULL     = 0x0
	TYMED_HGLOBAL  = 0x1
	TYMED_FILE     = 0x2
	TYMED_ISTREAM  = 0x4
	TYMED_ISTORAGE = 0x8
	TYMED_GDI      = 0x10
	TYMED_MFPICT   = 0x20
	TYMED_ENHMF    = 0x40
)

type iUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}
type iDataObjectVtbl struct {
	iUnknownVtbl
	GetData               uintptr
	GetDataHere           uintptr
	QueryGetData          uintptr
	GetCanonicalFormatEtc uintptr
	SetData               uintptr
	EnumFormatEtc         uintptr
	DAdvise               uintptr
	DUnadvise             uintptr
	EnumDAdvise           uintptr
}

type FORMATETC struct {
	CFormat        uint32
	DvTargetDevice uintptr
	Aspect         uint32
	Index          int32
	Tymed          uint32
}
type STGMEDIUM struct {
	Tymed          uint32
	UnionMember    uintptr
	PUnkForRelease *IUnknown
}
type ISequentialStreamVtbl struct {
	iUnknownVtbl
	Read  uintptr
	Write uintptr
}

type IUnknown struct {
	vtbl iUnknownVtbl
}

func (obj *IUnknown) Release() error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}

func (m STGMEDIUM) Release() {
	if m.PUnkForRelease == nil {
		win.ReleaseStgMedium((*win.STGMEDIUM)(unsafe.Pointer(&m)))
	}
}

func (m STGMEDIUM) Stream() (*IStream, error) {
	if m.Tymed != TYMED_ISTREAM {
		return nil, fmt.Errorf("invalid Tymed")
	}
	return (*IStream)(unsafe.Pointer(m.UnionMember)), nil
}

func (m STGMEDIUM) Bytes() ([]byte, error) {
	if m.Tymed != TYMED_HGLOBAL {
		return nil, fmt.Errorf("invalid Tymed")
	}
	size := GlobalSize(m.UnionMember)
	l := GlobalLock(m.UnionMember)
	defer GlobalUnlock(m.UnionMember)

	result := make([]byte, size)
	win.RtlCopyMemory(uintptr(unsafe.Pointer(&result[0])), uintptr(unsafe.Pointer(l)), size)

	return result, nil
}

type IDataObject struct {
	vtbl *iDataObjectVtbl
}

func (obj *IDataObject) Release() error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}

func (obj *IDataObject) GetData(formatEtc *FORMATETC, medium *STGMEDIUM) error {
	s2 := unsafe.Sizeof(*medium)
	_ = s2
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetData,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(formatEtc)),
		uintptr(unsafe.Pointer(medium)),
	)

	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}

func (obj *IDataObject) GetDataHere(formatEtc *FORMATETC, medium *STGMEDIUM) error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetDataHere,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(formatEtc)),
		uintptr(unsafe.Pointer(medium)),
	)

	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}
func (obj *IDataObject) QueryGetData(formatEtc *FORMATETC) error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.QueryGetData,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(formatEtc)),
		0,
	)

	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}
func (obj *IDataObject) EnumFormatEtc(direction uint32, pIEnumFORMATETC **IEnumFORMATETC) error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.EnumFormatEtc,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(direction),
		uintptr(unsafe.Pointer(pIEnumFORMATETC)),
	)

	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}

type iEnumFORMATETCVtbl struct {
	iUnknownVtbl
	Next  uintptr
	Skip  uintptr
	Reset uintptr
	Clone uintptr
}

type IEnumFORMATETC struct {
	vtbl *iEnumFORMATETCVtbl
}

func (obj *IEnumFORMATETC) Release() error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}

func (obj *IEnumFORMATETC) Next(formatEtc []FORMATETC, celtFetched *uint32) error {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Next,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(len(formatEtc)),
		uintptr(unsafe.Pointer(&formatEtc[0])),
		uintptr(unsafe.Pointer(celtFetched)),
		0,
		0,
	)
	if ret != 1 {
		return HRESULT(ret)
	}
	return nil
}

type EnumInstance struct {
	iEnumFORMATETC IEnumFORMATETC
	refCount       int32
	index          int
	formatEtc      []FORMATETC
}

var (
	IID_IDataObject    = win.GUID{0x10e, 0, 0, [8]byte{0xc0, 0, 0, 0, 0, 0, 0, 0x46}}
	IID_IUnknown       = win.GUID{0x000, 0, 0, [8]byte{0xc0, 0, 0, 0, 0, 0, 0, 0x46}}
	IID_IEnumFORMATETC = win.GUID{0x103, 0, 0, [8]byte{0xc0, 0, 0, 0, 0, 0, 0, 0x46}}
	IID_IStream        = win.GUID{0x00c, 0, 0, [8]byte{0xc0, 0, 0, 0, 0, 0, 0, 0x46}}
)

func newEnumInstance(fs []FORMATETC) *EnumInstance {
	var instance EnumInstance
	ie := &instance.iEnumFORMATETC
	ie.vtbl = new(iEnumFORMATETCVtbl)
	ie.vtbl.QueryInterface = syscall.NewCallback((*EnumInstance).QueryInterface)
	ie.vtbl.AddRef = syscall.NewCallback((*EnumInstance).AddRef)
	ie.vtbl.Release = syscall.NewCallback((*EnumInstance).Release)
	ie.vtbl.Next = syscall.NewCallback((*EnumInstance).Next)
	ie.vtbl.Skip = syscall.NewCallback((*EnumInstance).Skip)
	ie.vtbl.Reset = syscall.NewCallback((*EnumInstance).Reset)
	ie.vtbl.Clone = syscall.NewCallback((*EnumInstance).Clone)

	instance.refCount = 1
	if len(fs) > 0 {
		instance.formatEtc = make([]FORMATETC, len(fs))
		copy(instance.formatEtc, fs)
	}

	return &instance
}

func (i *EnumInstance) QueryInterface(riid win.REFGUID, ppvObject *uintptr) uintptr {
	if win.IsEqualGUID(riid, &IID_IEnumFORMATETC) ||
		win.IsEqualGUID(riid, &IID_IUnknown) {
		*ppvObject = uintptr(unsafe.Pointer(i))
		i.AddRef()
		return 0
	}
	*ppvObject = 0
	return E_NOINTERFACE
}

func (i *EnumInstance) AddRef() uintptr {
	n := atomic.AddInt32(&i.refCount, 1)
	return uintptr(n)
}

func (i *EnumInstance) Release() uintptr {
	n := atomic.AddInt32(&i.refCount, -1)
	if n == 0 {
		i = nil
		return 0
	}
	return uintptr(n)
}

func (i *EnumInstance) Next(celt uint32, rgelt *FORMATETC, pceltFetched *uint32) uintptr {
	r := make([]FORMATETC, celt)
	var idx uint32
	for i.index < len(i.formatEtc) && idx < celt {
		r[idx] = i.formatEtc[i.index]
		i.index++
		idx++
	}

	*rgelt = r[0]
	if pceltFetched != nil {
		*pceltFetched = idx
	}

	if idx != celt {
		return E_FAIL
	}
	return 0
}
func (i *EnumInstance) Skip(celt uint32) uintptr {
	if i.index+int(celt) > len(i.formatEtc) {
		return E_FAIL
	}
	i.index += int(celt)
	return 0
}
func (i *EnumInstance) Reset() uintptr {
	i.index = 0
	return 0
}
func (i *EnumInstance) Clone(ppEnum **IEnumFORMATETC) uintptr {
	ins := newEnumInstance(i.formatEtc)
	ins.index = i.index
	(*ppEnum) = (*IEnumFORMATETC)(unsafe.Pointer(ins))

	return 0
}
func CreateDataObject(c *CliprdrClient) *IDataObject {
	fmtetc := make([]FORMATETC, 2)
	stgmeds := make([]STGMEDIUM, 2)

	fmtetc[0].CFormat = RegisterClipboardFormat(CFSTR_FILEDESCRIPTORW)
	fmtetc[0].Aspect = DVASPECT_CONTENT
	fmtetc[0].Index = 0
	//fmtetc[0].DvTargetDevice = nil
	fmtetc[0].Tymed = TYMED_HGLOBAL
	stgmeds[0].Tymed = TYMED_HGLOBAL
	//stgmeds[0].UnionMember = nil
	stgmeds[0].PUnkForRelease = nil
	fmtetc[1].CFormat = RegisterClipboardFormat(CFSTR_FILECONTENTS)
	fmtetc[1].Aspect = DVASPECT_CONTENT
	fmtetc[1].Index = 0
	//fmtetc[1].DvTargetDevice = nil
	fmtetc[1].Tymed = TYMED_ISTREAM
	stgmeds[1].Tymed = TYMED_ISTREAM
	//stgmeds[1].UnionMember = nil
	stgmeds[1].PUnkForRelease = nil

	instance := newDataInstance()
	instance.refCount = 1
	instance.formatEtc = fmtetc
	instance.stgMedium = stgmeds
	instance.data = c

	return (*IDataObject)(unsafe.Pointer(instance))
}

type DataInstance struct {
	iDataObject IDataObject
	refCount    int32
	idx         int32
	formatEtc   []FORMATETC
	stgMedium   []STGMEDIUM
	streams     []*StreamInstance
	data        interface{}
}

func newDataInstance() *DataInstance {
	var instance DataInstance
	obj := &instance.iDataObject
	obj.vtbl = new(iDataObjectVtbl)
	obj.vtbl.QueryInterface = syscall.NewCallback((*DataInstance).QueryInterface)
	obj.vtbl.AddRef = syscall.NewCallback((*DataInstance).AddRef)
	obj.vtbl.Release = syscall.NewCallback((*DataInstance).Release)
	obj.vtbl.GetData = syscall.NewCallback((*DataInstance).GetData)
	obj.vtbl.GetDataHere = syscall.NewCallback((*DataInstance).GetDataHere)
	obj.vtbl.QueryGetData = syscall.NewCallback((*DataInstance).QueryGetData)
	obj.vtbl.GetCanonicalFormatEtc = syscall.NewCallback((*DataInstance).GetCanonicalFormatEtc)
	obj.vtbl.SetData = syscall.NewCallback((*DataInstance).SetData)
	obj.vtbl.EnumFormatEtc = syscall.NewCallback((*DataInstance).EnumFormatEtc)
	obj.vtbl.DAdvise = syscall.NewCallback((*DataInstance).DAdvise)
	obj.vtbl.DUnadvise = syscall.NewCallback((*DataInstance).DUnadvise)
	obj.vtbl.EnumDAdvise = syscall.NewCallback((*DataInstance).EnumDAdvise)

	return &instance
}

func (i *DataInstance) QueryInterface(riid win.REFGUID, ppvObject *uintptr) uintptr {
	if win.IsEqualGUID(riid, &IID_IDataObject) ||
		win.IsEqualGUID(riid, &IID_IUnknown) {
		*ppvObject = uintptr(unsafe.Pointer(i))
		i.AddRef()
		return 0
	}

	*ppvObject = 0
	return E_NOINTERFACE
}

func (i *DataInstance) AddRef() uintptr {
	n := atomic.AddInt32(&i.refCount, 1)
	return uintptr(n)
}

func (i *DataInstance) Release() uintptr {
	n := atomic.AddInt32(&i.refCount, -1)
	if n == 0 {
		i = nil
		return 0
	}
	return uintptr(n)
}

func (i *DataInstance) GetData(formatEtc *FORMATETC, medium *STGMEDIUM) uintptr {
	idx := -1
	for j, f := range i.formatEtc {
		if formatEtc.Tymed&f.Tymed != 0 &&
			formatEtc.CFormat == f.CFormat &&
			formatEtc.Aspect&f.Aspect != 0 {
			idx = j
		}
	}
	if idx == -1 {
		return E_FORMATETC
	}
	glog.Debugf("GetData:%+v, %s", formatEtc.CFormat, GetClipboardFormatName(formatEtc.CFormat))

	medium.Tymed = i.formatEtc[idx].Tymed

	if i.formatEtc[idx].CFormat == RegisterClipboardFormat(CFSTR_FILEDESCRIPTORW) {
		c := i.data.(*CliprdrClient)
		if remoteid, ok := c.formatIdMap[i.formatEtc[idx].CFormat]; ok {
			c.sendFormatDataRequest(remoteid)
			b := <-c.reply
			if len(b) == 0 {
				return E_FAIL
			}
			medium.UnionMember = HmemAlloc(b)
			var dsc FileGroupDescriptor
			dsc.Unpack(b)
			if dsc.CItems > 0 {
				glog.Debug("Items:", dsc.CItems)
				i.streams = make([]*StreamInstance, dsc.CItems)
				var j uint32
				for j = 0; j < dsc.CItems; j++ {
					glog.Debug("FileName:", core.UnicodeDecode(dsc.Fgd[j].FileName))
					s := newStream(j, i.data, &dsc.Fgd[j])
					i.streams[j] = s
				}
			}
		}

	} else if i.formatEtc[idx].CFormat == RegisterClipboardFormat(CFSTR_FILECONTENTS) {
		if formatEtc.Index >= 0 && formatEtc.Index < int32(len(i.streams)) {
			medium.UnionMember = uintptr(unsafe.Pointer(i.streams[formatEtc.Index]))
			i.AddRef()
		} else {
			return E_INVALIDARG
		}
	} else {
		return E_UNEXPECTED
	}

	return 0
}

func (i *DataInstance) GetDataHere(formatEtc *FORMATETC, medium *STGMEDIUM) uintptr {
	return E_NOTIMPL
}

func (i *DataInstance) QueryGetData(formatEtc *FORMATETC) uintptr {
	for _, f := range i.formatEtc {
		if formatEtc.Tymed&f.Tymed != 0 &&
			formatEtc.CFormat == f.CFormat &&
			formatEtc.Aspect&f.Aspect != 0 {
			return 0
		}
	}
	return E_FORMATETC
}

func (i *DataInstance) GetCanonicalFormatEtc(informatEtc, outformatEtc *FORMATETC) uintptr {
	return E_NOTIMPL
}

func (i *DataInstance) SetData(formatEtc *FORMATETC, medium *STGMEDIUM, r bool) uintptr {
	return E_NOTIMPL
}

func (i *DataInstance) EnumFormatEtc(dwDirection uint32, ppenumFormatEtc **IEnumFORMATETC) uintptr {
	if dwDirection == 1 {
		ins := newEnumInstance(i.formatEtc)
		(*ppenumFormatEtc) = (*IEnumFORMATETC)(unsafe.Pointer(ins))
		return 0
	}

	return E_NOTIMPL
}

func (i *DataInstance) DAdvise(formatEtc *FORMATETC, advf uint32, pAdvSink uintptr, pdwConnection *uint32) uintptr {
	return E_ADVISENOTSUPPORTED
}
func (i *DataInstance) DUnadvise(dwDirection uint32) uintptr {
	return E_ADVISENOTSUPPORTED
}
func (i *DataInstance) EnumDAdvise(ppenumAdvise uintptr) uintptr {
	return E_ADVISENOTSUPPORTED
}

const (
	STREAM_SEEK_SET = 0
	STREAM_SEEK_CUR = 1
	STREAM_SEEK_END = 2
)

const (
	STATFLAG_DEFAULT = 0
	STATFLAG_NONAME  = 1
	STATFLAG_NOOPEN  = 2
)

const (
	STG_E_INSUFFICIENTMEMORY = 0x80030008
	STG_E_INVALIDFLAG        = 0x800300FF
)

const (
	STGTY_STORAGE   = 1
	STGTY_STREAM    = 2
	STGTY_LOCKBYTES = 3
	STGTY_PROPERTY  = 4
)

const (
	LOCK_WRITE     = 1
	LOCK_EXCLUSIVE = 2
	LOCK_ONLYONCE  = 4
)

const (
	GENERIC_READ    = 0x80000000
	GENERIC_WRITE   = 0x40000000
	GENERIC_EXECUTE = 0x20000000
)

type IStreamVtbl struct {
	ISequentialStreamVtbl
	Seek         uintptr
	SetSize      uintptr
	CopyTo       uintptr
	Commit       uintptr
	Revert       uintptr
	LockRegion   uintptr
	UnlockRegion uintptr
	Stat         uintptr
	Clone        uintptr
}

type IStream struct {
	vtbl *IStreamVtbl
}
type ULARGE_INTEGER struct {
	QuadPart uint64
}
type LARGE_INTEGER struct {
	QuadPart int64
}

func (l *ULARGE_INTEGER) LowPart() *uint32 {
	return (*uint32)(unsafe.Pointer(&l.QuadPart))
}
func (l *ULARGE_INTEGER) HighPart() *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&l.QuadPart)) + uintptr(4)))
}
func (l *LARGE_INTEGER) LowPart() *uint32 {
	return (*uint32)(unsafe.Pointer(&l.QuadPart))
}
func (l *LARGE_INTEGER) HighPart() *int32 {
	return (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&l.QuadPart)) + uintptr(4)))
}

type StreamInstance struct {
	iStream  IStream
	streamId uint32
	refCount int32
	index    uint32
	lSize    ULARGE_INTEGER
	lOffset  ULARGE_INTEGER
	dsc      FileDescriptor
	data     interface{}
}
type STATSTG struct {
	pwcsName          []uint16
	sType             uint32
	cbSize            ULARGE_INTEGER
	mtime             win.FILETIME
	ctime             win.FILETIME
	atime             win.FILETIME
	grfMode           uint32
	grfLocksSupported uint32
	clsid             win.CLSID
	grfStateBits      uint32
	reserved          uint32
}

func (obj *IStream) Read(buffer []byte) (int, error) {
	bufPtr := &buffer[0]
	var read uint32
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Read,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(bufPtr)),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&read)),
		0,
		0,
	)
	if ret == 1 {
		return int(read), io.EOF
	}
	if ret != 0 {
		return int(read), HRESULT(ret)
	}
	return int(read), nil
}
func (obj *IStream) Close() error {
	return obj.Release()
}
func (obj *IStream) Release() error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	if ret != 0 {
		return HRESULT(ret)
	}
	return nil
}
func newStream(index uint32, data interface{}, dsc *FileDescriptor) *StreamInstance {
	var instance StreamInstance
	is := &instance.iStream
	is.vtbl = new(IStreamVtbl)
	is.vtbl.QueryInterface = syscall.NewCallback((*StreamInstance).QueryInterface)
	is.vtbl.AddRef = syscall.NewCallback((*StreamInstance).AddRef)
	is.vtbl.Release = syscall.NewCallback((*StreamInstance).Release)
	is.vtbl.Read = syscall.NewCallback((*StreamInstance).Read)
	is.vtbl.Write = syscall.NewCallback((*StreamInstance).Write)
	is.vtbl.Seek = syscall.NewCallback((*StreamInstance).Seek)
	is.vtbl.SetSize = syscall.NewCallback((*StreamInstance).SetSize)
	is.vtbl.CopyTo = syscall.NewCallback((*StreamInstance).CopyTo)
	is.vtbl.Commit = syscall.NewCallback((*StreamInstance).Commit)
	is.vtbl.Revert = syscall.NewCallback((*StreamInstance).Revert)
	is.vtbl.LockRegion = syscall.NewCallback((*StreamInstance).LockRegion)
	is.vtbl.UnlockRegion = syscall.NewCallback((*StreamInstance).UnlockRegion)
	is.vtbl.Stat = syscall.NewCallback((*StreamInstance).Stat)
	is.vtbl.Clone = syscall.NewCallback((*StreamInstance).Clone)

	instance.streamId = (uint32)(reflect.ValueOf(&instance).Pointer())
	instance.refCount = 1
	instance.dsc = *dsc
	instance.data = data
	instance.index = index
	if !instance.dsc.hasFileSize() && !instance.dsc.isDir() {
		c := data.(*CliprdrClient)
		var r CliprdrFileContentsRequest
		r.StreamId = instance.streamId
		r.Lindex = instance.index
		r.DwFlags = FILECONTENTS_SIZE
		r.CbRequested = 8
		c.sendFormatContentsRequest(r)
		b := <-c.reply
		instance.lSize.QuadPart = core.BytesToUint64(b)
	} else {
		b := &bytes.Buffer{}
		core.WriteUInt32LE(dsc.FileSizeLow, b)
		core.WriteUInt32LE(dsc.FileSizeHigh, b)
		instance.lSize.QuadPart = core.BytesToUint64(b.Bytes())
	}

	return &instance
}

func (i *StreamInstance) QueryInterface(riid win.REFGUID, ppvObject *uintptr) uintptr {
	if win.IsEqualGUID(riid, &IID_IStream) ||
		win.IsEqualGUID(riid, &IID_IUnknown) {
		*ppvObject = uintptr(unsafe.Pointer(i))
		i.AddRef()
		return 0
	}
	*ppvObject = 0
	return E_NOINTERFACE
}

func (i *StreamInstance) AddRef() uintptr {
	n := atomic.AddInt32(&i.refCount, 1)
	return uintptr(n)
}

func (i *StreamInstance) Release() uintptr {
	n := atomic.AddInt32(&i.refCount, -1)
	if n == 0 {
		i = nil
		return 0
	}
	return uintptr(n)
}

func (i *StreamInstance) Read(pv uintptr, cb uint32, cbRead *uint32) uintptr {
	glog.Debug("StreamInstance Read:", i.lOffset.QuadPart, i.lSize.QuadPart)
	if i.lOffset.QuadPart >= i.lSize.QuadPart {
		return 1
	}

	c := i.data.(*CliprdrClient)
	*cbRead = 0
	var r CliprdrFileContentsRequest
	r.StreamId = i.streamId
	r.Lindex = i.index
	r.DwFlags = FILECONTENTS_RANGE
	r.NPositionHigh = *(i.lOffset.HighPart())
	r.NPositionLow = *(i.lOffset.LowPart())
	r.CbRequested = cb
	c.sendFormatContentsRequest(r)
	b := <-c.reply
	if len(b) == 0 {
		return E_FAIL
	}
	win.RtlCopyMemory(pv, uintptr(unsafe.Pointer(&b[0])), win.SIZE_T(len(b)))
	*cbRead = uint32(len(b))
	i.lOffset.QuadPart += uint64(len(b))
	glog.Debug("StreamInstance Read:", *cbRead, cb)
	if *cbRead < cb {
		return 1
	}
	return 0
}

func (i *StreamInstance) Write(pv uintptr, cb uint32, cbWritten *uint32) uintptr {
	return E_ACCESSDENIED
}

func (i *StreamInstance) Seek(dlibMove LARGE_INTEGER, dwOrigin uint32, plibNewPosition *ULARGE_INTEGER) uintptr {
	glog.Debug("StreamInstance Seek:", dwOrigin, dlibMove, plibNewPosition)
	var newoffset uint64 = i.lOffset.QuadPart
	switch dwOrigin {
	case STREAM_SEEK_SET:
		newoffset = uint64(dlibMove.QuadPart)
		break

	case STREAM_SEEK_CUR:
		newoffset += uint64(dlibMove.QuadPart)
		break

	case STREAM_SEEK_END:
		newoffset = i.lSize.QuadPart + uint64(dlibMove.QuadPart)
		break

	default:
		return E_INVALIDARG
	}
	glog.Debug("StreamInstance Seek:", newoffset, i.lSize.QuadPart)
	if newoffset < 0 || newoffset >= i.lSize.QuadPart {
		return 1
	}
	i.lOffset.QuadPart = newoffset

	if plibNewPosition != nil {
		(*plibNewPosition).QuadPart = i.lOffset.QuadPart
	}

	return 0
}

func (i *StreamInstance) SetSize(libNewSize ULARGE_INTEGER) uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) CopyTo(pstm *IStream, cb ULARGE_INTEGER, cbRead, cbWritten *ULARGE_INTEGER) uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) Commit(grfCommitFlags uint32) uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) Revert() uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) LockRegion(libOffset, cb ULARGE_INTEGER, dwLockType uint32) uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) UnlockRegion(libOffset, cb ULARGE_INTEGER, dwLockType uint32) uintptr {
	return E_NOTIMPL
}

func (i *StreamInstance) Stat(pstatstg *STATSTG, grfStatFlag uint32) uintptr {
	switch grfStatFlag {
	case STATFLAG_DEFAULT:
		return STG_E_INSUFFICIENTMEMORY

	case STATFLAG_NONAME:
		pstatstg.cbSize.QuadPart = i.lSize.QuadPart
		pstatstg.grfLocksSupported = LOCK_EXCLUSIVE
		pstatstg.grfMode = GENERIC_READ
		pstatstg.grfStateBits = 0
		pstatstg.sType = STGTY_STREAM
		break

	case STATFLAG_NOOPEN:
		return STG_E_INVALIDFLAG

	default:
		return STG_E_INVALIDFLAG
	}

	return 0
}

func (i *StreamInstance) Clone(ppstm **IStream) uintptr {
	return E_NOTIMPL
}
