// cliprdr_windows.go
package cliprdr

import (
	"bytes"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/shirou/w32"
	"github.com/tomatome/grdp/glog"

	"github.com/tomatome/grdp/core"

	"github.com/tomatome/win"
)

const (
	CFSTR_SHELLIDLIST         = "Shell IDList Array"
	CFSTR_SHELLIDLISTOFFSET   = "Shell Object Offsets"
	CFSTR_NETRESOURCES        = "Net Resource"
	CFSTR_FILECONTENTS        = "FileContents"
	CFSTR_FILENAMEA           = "FileName"
	CFSTR_FILENAMEMAPA        = "FileNameMap"
	CFSTR_FILEDESCRIPTORA     = "FileGroupDescriptor"
	CFSTR_INETURLA            = "UniformResourceLocator"
	CFSTR_SHELLURL            = CFSTR_INETURLA
	CFSTR_FILENAMEW           = "FileNameW"
	CFSTR_FILENAMEMAPW        = "FileNameMapW"
	CFSTR_FILEDESCRIPTORW     = "FileGroupDescriptorW"
	CFSTR_INETURLW            = "UniformResourceLocatorW"
	CFSTR_PRINTERGROUP        = "PrinterFriendlyName"
	CFSTR_INDRAGLOOP          = "InShellDragLoop"
	CFSTR_PASTESUCCEEDED      = "Paste Succeeded"
	CFSTR_PERFORMEDDROPEFFECT = "Performed DropEffect"
	CFSTR_PREFERREDDROPEFFECT = "Preferred DropEffect"
)
const DVASPECT_CONTENT = 0x1

const (
	CF_TEXT         = 1
	CF_BITMAP       = 2
	CF_METAFILEPICT = 3
	CF_SYLK         = 4
	CF_DIF          = 5
	CF_TIFF         = 6
	CF_OEMTEXT      = 7
	CF_DIB          = 8
	CF_PALETTE      = 9
	CF_PENDATA      = 10
	CF_RIFF         = 11
	CF_WAVE         = 12
	CF_UNICODETEXT  = 13
	CF_ENHMETAFILE  = 14
	CF_HDROP        = 15
	CF_LOCALE       = 16
	CF_DIBV5        = 17
	CF_MAX          = 18
)
const (
	WM_CLIPRDR_MESSAGE = (w32.WM_USER + 156)
	OLE_SETCLIPBOARD   = 1
)

type Control struct {
	hwnd       uintptr
	dataObject *IDataObject
}

func (c *Control) withOpenClipboard(f func()) {
	if OpenClipboard(c.hwnd) {
		f()
		CloseClipboard()
	}
}
func ClipWatcher(c *CliprdrClient) {
	win.OleInitialize(0)
	defer win.OleUninitialize()
	className := syscall.StringToUTF16Ptr("ClipboardHiddenMessageProcessor")
	windowName := syscall.StringToUTF16Ptr("rdpclip")
	wndClassEx := w32.WNDCLASSEX{
		ClassName: className,
		WndProc: syscall.NewCallback(func(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
			switch msg {
			case w32.WM_CLIPBOARDUPDATE:
				glog.Info("info: WM_CLIPBOARDUPDATE wParam:", wParam)
				glog.Debug("IsClipboardOwner:", IsClipboardOwner(win.HWND(c.hwnd)))
				glog.Debug("OleIsCurrentClipboard:", OleIsCurrentClipboard(c.dataObject))
				if !IsClipboardOwner(win.HWND(c.hwnd)) && int(wParam) != 0 &&
					!OleIsCurrentClipboard(c.dataObject) {
					c.sendFormatListPDU()
				}

			case w32.WM_RENDERALLFORMATS:
				glog.Info("info: WM_RENDERALLFORMATS")
				c.withOpenClipboard(func() {
					EmptyClipboard()
				})

			case w32.WM_RENDERFORMAT:
				glog.Info("info: WM_RENDERFORMAT wParam:", wParam)
				formatId := uint32(wParam)
				c.sendFormatDataRequest(formatId)
				b := <-c.reply
				hmem := HmemAlloc(b)
				SetClipboardData(formatId, hmem)

			case WM_CLIPRDR_MESSAGE:
				glog.Info("info: WM_CLIPRDR_MESSAGE wParam:", wParam)
				if wParam == OLE_SETCLIPBOARD {
					if !OleIsCurrentClipboard(c.dataObject) {
						o := CreateDataObject(c)
						OleSetClipboard(o)
						c.dataObject = o
					}
				}
			default:
				return w32.DefWindowProc(hwnd, msg, wParam, lParam)
			}
			return 0
		}),
		Style: w32.CS_OWNDC,
	}
	wndClassEx.Size = uint32(unsafe.Sizeof(wndClassEx))
	w32.RegisterClassEx(&wndClassEx)

	hwnd := w32.CreateWindowEx(w32.WS_EX_LEFT, className, windowName, 0, 0, 0, 0, 0, w32.HWND_MESSAGE, 0, 0, nil)
	c.hwnd = uintptr(hwnd)
	w32.AddClipboardFormatListener(hwnd)
	defer w32.RemoveClipboardFormatListener(hwnd)

	msg := w32.MSG{}
	for w32.GetMessage(&msg, 0, 0, 0) > 0 {
		w32.DispatchMessage(&msg)
	}

}
func OpenClipboard(hwnd uintptr) bool {
	return win.OpenClipboard(win.HWND(hwnd))
}
func CloseClipboard() bool {
	return win.CloseClipboard()
}
func CountClipboardFormats() int32 {
	return win.CountClipboardFormats()
}
func IsClipboardFormatAvailable(id uint32) bool {
	return win.IsClipboardFormatAvailable(win.UINT(id))
}
func EnumClipboardFormats(formatId uint32) uint32 {
	id := win.EnumClipboardFormats(win.UINT(formatId))
	return uint32(id)
}
func GetClipboardFormatName(id uint32) string {
	buf := make([]uint16, 250)
	n := win.GetClipboardFormatName(win.UINT(id), win.LPWSTR(unsafe.Pointer(&buf[0])), int32(len(buf)))
	return string(utf16.Decode(buf[:n]))
}
func EmptyClipboard() bool {
	return win.EmptyClipboard()
}
func RegisterClipboardFormat(format string) uint32 {
	id := win.RegisterClipboardFormat(format)
	return uint32(id)
}
func IsClipboardOwner(h win.HWND) bool {
	hwnd := win.GetClipboardOwner()
	return h == hwnd
}

func HmemAlloc(data []byte) uintptr {
	ln := (len(data))
	h := win.GlobalAlloc(0x0002, win.SIZE_T(ln))
	if h == 0 {
		return uintptr(h)
	}
	if ln == 0 {
		return uintptr(h)
	}
	l := win.GlobalLock(h)
	defer win.GlobalUnlock(h)

	win.RtlCopyMemory(uintptr(unsafe.Pointer(l)), uintptr(unsafe.Pointer(&data[0])), win.SIZE_T(ln))

	return uintptr(h)

}
func SetClipboardData(formatId uint32, hmem uintptr) bool {
	r := win.SetClipboardData(win.UINT(formatId), win.HANDLE(hmem))
	if r == 0 {
		//glog.Error("SetClipboardData failed:", formatId, hmem)
		return false
	}
	return true
}
func GetClipboardData(formatId uint32) string {
	r := win.GetClipboardData(win.UINT(formatId))
	if r == 0 {
		return ""
	}

	h := win.GlobalHandle(uintptr(r))
	size := win.GlobalSize(h)
	l := win.GlobalLock(h)
	defer win.GlobalUnlock(h)

	result := make([]byte, size)
	win.RtlCopyMemory(uintptr(unsafe.Pointer(&result[0])), uintptr(unsafe.Pointer(l)), size)

	return core.UnicodeDecode(result)
}

func GetFormatList(hwnd uintptr) []CliprdrFormat {
	list := make([]CliprdrFormat, 0, 10)
	if OpenClipboard(hwnd) {
		n := CountClipboardFormats()
		if IsClipboardFormatAvailable(CF_HDROP) {
			formatId := RegisterClipboardFormat(CFSTR_FILEDESCRIPTORW)
			var f CliprdrFormat
			f.FormatId = formatId
			f.FormatName = CFSTR_FILEDESCRIPTORW
			list = append(list, f)
			formatId = RegisterClipboardFormat(CFSTR_FILECONTENTS)
			var f1 CliprdrFormat
			f1.FormatId = formatId
			f1.FormatName = CFSTR_FILECONTENTS
			list = append(list, f1)
		} else {
			var id uint32
			for i := 0; i < int(n); i++ {
				id = EnumClipboardFormats(id)
				name := GetClipboardFormatName(id)
				var f CliprdrFormat
				f.FormatId = id
				f.FormatName = name
				list = append(list, f)
			}
		}
		CloseClipboard()
	}
	return list
}

func OleGetClipboard() *IDataObject {
	var dataObject *IDataObject
	win.OleGetClipboard((**win.IDataObject)(unsafe.Pointer(&dataObject)))
	return dataObject
}

func OleSetClipboard(dataObject *IDataObject) bool {
	r := win.OleSetClipboard((*win.IDataObject)(unsafe.Pointer(dataObject)))
	if r != 0 {
		glog.Error("OleSetClipboard failed")
		return false
	}
	return true
}

func OleIsCurrentClipboard(dataObject *IDataObject) bool {
	r := win.OleIsCurrentClipboard((*win.IDataObject)(unsafe.Pointer(dataObject)))
	if r != 0 {
		return false
	}
	return true
}
func GlobalSize(hMem uintptr) win.SIZE_T {
	return win.GlobalSize(win.HGLOBAL(hMem))
}
func GlobalLock(hMem uintptr) uintptr {
	r := win.GlobalLock(win.HGLOBAL(hMem))

	return uintptr(r)
}
func GlobalUnlock(hMem uintptr) {
	win.GlobalUnlock(win.HGLOBAL(hMem))
}

func (c *Control) SendCliprdrMessage() {
	win.PostMessage(win.HWND(c.hwnd), WM_CLIPRDR_MESSAGE, OLE_SETCLIPBOARD, 0)
}
func GetFileInfo(sys interface{}) (uint32, []byte, uint32, uint32) {
	f := sys.(*syscall.Win32FileAttributeData)
	b := &bytes.Buffer{}
	core.WriteUInt32LE(f.LastWriteTime.LowDateTime, b)
	core.WriteUInt32LE(f.LastWriteTime.HighDateTime, b)
	return f.FileAttributes, b.Bytes(), f.FileSizeHigh, f.FileSizeLow
}

func GetFileNames() []string {
	o := OleGetClipboard()
	var formatEtc FORMATETC
	var stgMedium STGMEDIUM
	formatEtc.CFormat = CF_HDROP
	formatEtc.Tymed = TYMED_HGLOBAL
	formatEtc.Aspect = DVASPECT_CONTENT
	formatEtc.Index = -1
	o.GetData(&formatEtc, &stgMedium)
	b, _ := stgMedium.Bytes()
	//DROPFILES
	r := bytes.NewReader(b[20:])
	fs := make([]string, 0, 20)
	for r.Len() > 0 {
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

		if name == "" {
			continue
		}
		fs = append(fs, name)
	}

	return fs
}

const (
	/* File attribute flags */
	FILE_SHARE_READ   = 0x00000001
	FILE_SHARE_WRITE  = 0x00000002
	FILE_SHARE_DELETE = 0x00000004

	FILE_ATTRIBUTE_READONLY            = 0x00000001
	FILE_ATTRIBUTE_HIDDEN              = 0x00000002
	FILE_ATTRIBUTE_SYSTEM              = 0x00000004
	FILE_ATTRIBUTE_DIRECTORY           = 0x00000010
	FILE_ATTRIBUTE_ARCHIVE             = 0x00000020
	FILE_ATTRIBUTE_DEVICE              = 0x00000040
	FILE_ATTRIBUTE_NORMAL              = 0x00000080
	FILE_ATTRIBUTE_TEMPORARY           = 0x00000100
	FILE_ATTRIBUTE_SPARSE_FILE         = 0x00000200
	FILE_ATTRIBUTE_REPARSE_POINT       = 0x00000400
	FILE_ATTRIBUTE_COMPRESSED          = 0x00000800
	FILE_ATTRIBUTE_OFFLINE             = 0x00001000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED = 0x00002000
	FILE_ATTRIBUTE_ENCRYPTED           = 0x00004000
	FILE_ATTRIBUTE_INTEGRITY_STREAM    = 0x00008000
	FILE_ATTRIBUTE_VIRTUAL             = 0x00010000
	FILE_ATTRIBUTE_NO_SCRUB_DATA       = 0x00020000
	FILE_ATTRIBUTE_EA                  = 0x00040000
)

type DROPFILES struct {
	pFiles uintptr
	pt     uintptr
	fNC    bool
	fWide  bool
}
