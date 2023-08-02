package pdu

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tomatome/grdp/glog"

	"github.com/tomatome/grdp/core"
)

type ControlFlag uint8

const (
	TS_STANDARD             = 0x01
	TS_SECONDARY            = 0x02
	TS_BOUNDS               = 0x04
	TS_TYPE_CHANGE          = 0x08
	TS_DELTA_COORDINATES    = 0x10
	TS_ZERO_BOUNDS_DELTAS   = 0x20
	TS_ZERO_FIELD_BYTE_BIT0 = 0x40
	TS_ZERO_FIELD_BYTE_BIT1 = 0x80
)

type PrimaryOrderType uint8

const (
	ORDER_TYPE_DSTBLT = 0x00 //0
	ORDER_TYPE_PATBLT = 0x01 //1
	ORDER_TYPE_SCRBLT = 0x02 //2
	//ORDER_TYPE_DRAWNINEGRID       = 0x07 //7
	//ORDER_TYPE_MULTI_DRAWNINEGRID = 0x08 //8
	ORDER_TYPE_LINETO     = 0x09 //9
	ORDER_TYPE_OPAQUERECT = 0x0A //10
	ORDER_TYPE_SAVEBITMAP = 0x0B //11
	ORDER_TYPE_MEMBLT     = 0x0D //13
	ORDER_TYPE_MEM3BLT    = 0x0E //14
	//ORDER_TYPE_MULTIDSTBLT        = 0x0F //15
	//ORDER_TYPE_MULTIPATBLT        = 0x10 //16
	//ORDER_TYPE_MULTISCRBLT        = 0x11 //17
	//ORDER_TYPE_MULTIOPAQUERECT    = 0x12 //18
	//ORDER_TYPE_FAST_INDEX         = 0x13 //19
	ORDER_TYPE_POLYGON_SC = 0x14 //20
	ORDER_TYPE_POLYGON_CB = 0x15 //21
	ORDER_TYPE_POLYLINE   = 0x16 //22
	//ORDER_TYPE_FAST_GLYPH         = 0x18 //24
	ORDER_TYPE_ELLIPSE_SC = 0x19 //25
	ORDER_TYPE_ELLIPSE_CB = 0x1A //26
	ORDER_TYPE_TEXT2      = 0x1B //27
)

type SecondaryOrderType uint8

const (
	ORDER_TYPE_BITMAP_UNCOMPRESSED     = 0x00
	ORDER_TYPE_CACHE_COLOR_TABLE       = 0x01
	ORDER_TYPE_CACHE_BITMAP_COMPRESSED = 0x02
	ORDER_TYPE_CACHE_GLYPH             = 0x03
	ORDER_TYPE_BITMAP_UNCOMPRESSED_V2  = 0x04
	ORDER_TYPE_BITMAP_COMPRESSED_V2    = 0x05
	ORDER_TYPE_CACHE_BRUSH             = 0x07
	ORDER_TYPE_BITMAP_COMPRESSED_V3    = 0x08
)

func (s SecondaryOrderType) String() string {
	name := "Unknown"
	switch s {
	case ORDER_TYPE_BITMAP_UNCOMPRESSED:
		name = "Cache Bitmap"
	case ORDER_TYPE_CACHE_COLOR_TABLE:
		name = "Cache Color Table"
	case ORDER_TYPE_CACHE_BITMAP_COMPRESSED:
		name = "Cache Bitmap (Compressed)"
	case ORDER_TYPE_CACHE_GLYPH:
		name = "Cache Glyph"
	case ORDER_TYPE_BITMAP_UNCOMPRESSED_V2:
		name = "Cache Bitmap V2"
	case ORDER_TYPE_BITMAP_COMPRESSED_V2:
		name = "Cache Bitmap V2 (Compressed)"
	case ORDER_TYPE_CACHE_BRUSH:
		name = "Cache Brush"
	case ORDER_TYPE_BITMAP_COMPRESSED_V3:
		name = "Cache Bitmap V3"
	}
	return fmt.Sprintf("[0x%02d] %s", s, name)
}

/* Alternate Secondary Drawing Orders */
const (
	ORDER_TYPE_SWITCH_SURFACE          = 0x00
	ORDER_TYPE_CREATE_OFFSCREEN_BITMAP = 0x01
	ORDER_TYPE_STREAM_BITMAP_FIRST     = 0x02
	ORDER_TYPE_STREAM_BITMAP_NEXT      = 0x03
	ORDER_TYPE_CREATE_NINE_GRID_BITMAP = 0x04
	ORDER_TYPE_GDIPLUS_FIRST           = 0x05
	ORDER_TYPE_GDIPLUS_NEXT            = 0x06
	ORDER_TYPE_GDIPLUS_END             = 0x07
	ORDER_TYPE_GDIPLUS_CACHE_FIRST     = 0x08
	ORDER_TYPE_GDIPLUS_CACHE_NEXT      = 0x09
	ORDER_TYPE_GDIPLUS_CACHE_END       = 0x0A
	ORDER_TYPE_WINDOW                  = 0x0B
	ORDER_TYPE_COMPDESK_FIRST          = 0x0C
	ORDER_TYPE_FRAME_MARKER            = 0x0D
)

const (
	GLYPH_FRAGMENT_NOP = 0x00
	GLYPH_FRAGMENT_USE = 0xFE
	GLYPH_FRAGMENT_ADD = 0xFF

	CBR2_HEIGHT_SAME_AS_WIDTH      = 0x01
	CBR2_PERSISTENT_KEY_PRESENT    = 0x02
	CBR2_NO_BITMAP_COMPRESSION_HDR = 0x08
	CBR2_DO_NOT_CACHE              = 0x10
)

const (
	ORDER_PRIMARY = iota
	ORDER_SECONDARY
	ORDER_ALTSEC
)

type OrderPdu struct {
	ControlFlags uint8
	Type         int
	Altsec       *Altsec
	Primary      *Primary
	Secondary    *Secondary
}

func (o *OrderPdu) HasBounds() bool {
	return o.ControlFlags&TS_BOUNDS != 0
}

type Altsec struct {
}

type Secondary struct {
}

type Primary struct {
	Bounds Bounds
	Data   PrimaryOrder
}

type FastPathOrdersPDU struct {
	NumberOrders uint16
	OrderPdus    []OrderPdu
}

func (*FastPathOrdersPDU) FastPathUpdateType() uint8 {
	return FASTPATH_UPDATETYPE_ORDERS
}

func (f *FastPathOrdersPDU) Unpack(r io.Reader) error {
	f.NumberOrders, _ = core.ReadUint16LE(r)
	//glog.Info("NumberOrders:", f.NumberOrders)
	for i := 0; i < int(f.NumberOrders); i++ {
		var o OrderPdu
		o.ControlFlags, _ = core.ReadUInt8(r)
		if o.ControlFlags&TS_STANDARD == 0 {
			//glog.Info("Altsec order")
			o.processAltsecOrder(r)
			o.Type = ORDER_ALTSEC
			//return errors.New("Not support")
		} else if o.ControlFlags&TS_SECONDARY != 0 {
			//glog.Info("Secondary order")
			o.processSecondaryOrder(r)
			o.Type = ORDER_SECONDARY
		} else {
			//glog.Info("Primary order")
			o.processPrimaryOrder(r)
			o.Type = ORDER_PRIMARY
		}

		if f.OrderPdus == nil {
			f.OrderPdus = make([]OrderPdu, 0, f.NumberOrders)
		}
		f.OrderPdus = append(f.OrderPdus, o)
	}
	return nil
}
func (o *OrderPdu) processAltsecOrder(r io.Reader) error {
	orderType := o.ControlFlags >> 2
	//glog.Info("Altsec:", orderType)
	switch orderType {
	case ORDER_TYPE_SWITCH_SURFACE:
	case ORDER_TYPE_CREATE_OFFSCREEN_BITMAP:
	case ORDER_TYPE_STREAM_BITMAP_FIRST:
	case ORDER_TYPE_STREAM_BITMAP_NEXT:
	case ORDER_TYPE_CREATE_NINE_GRID_BITMAP:
	case ORDER_TYPE_GDIPLUS_FIRST:
	case ORDER_TYPE_GDIPLUS_NEXT:
	case ORDER_TYPE_GDIPLUS_END:
	case ORDER_TYPE_GDIPLUS_CACHE_FIRST:
	case ORDER_TYPE_GDIPLUS_CACHE_NEXT:
	case ORDER_TYPE_GDIPLUS_CACHE_END:
	case ORDER_TYPE_WINDOW:
	case ORDER_TYPE_COMPDESK_FIRST:
	case ORDER_TYPE_FRAME_MARKER:
		core.ReadUInt32LE(r)
	}

	return nil
}
func (o *OrderPdu) processSecondaryOrder(r io.Reader) error {
	var sec Secondary
	length, _ := core.ReadUint16LE(r)
	flags, _ := core.ReadUint16LE(r)
	orderType, _ := core.ReadUInt8(r)

	glog.Info("Secondary:", SecondaryOrderType(orderType))

	b, _ := core.ReadBytes(int(length)+13-6, r)
	r0 := bytes.NewReader(b)

	switch orderType {
	case ORDER_TYPE_BITMAP_UNCOMPRESSED:
		fallthrough
	case ORDER_TYPE_CACHE_BITMAP_COMPRESSED:
		compressed := (orderType == ORDER_TYPE_CACHE_BITMAP_COMPRESSED)
		sec.updateCacheBitmapOrder(r0, compressed, flags)
	case ORDER_TYPE_BITMAP_UNCOMPRESSED_V2:
		fallthrough
	case ORDER_TYPE_BITMAP_COMPRESSED_V2:
		compressed := (orderType == ORDER_TYPE_BITMAP_COMPRESSED_V2)
		sec.updateCacheBitmapV2Order(r0, compressed, flags)
	case ORDER_TYPE_BITMAP_COMPRESSED_V3:
		sec.updateCacheBitmapV3Order(r0, flags)
	case ORDER_TYPE_CACHE_COLOR_TABLE:
		sec.updateCacheColorTableOrder(r0, flags)
	case ORDER_TYPE_CACHE_GLYPH:
		sec.updateCacheGlyphOrder(r0, flags)
	case ORDER_TYPE_CACHE_BRUSH:
		sec.updateCacheBrushOrder(r0, flags)
	default:
		glog.Debugf("Unsupport order type 0x%x", orderType)
	}

	return nil
}
func (b *Bounds) updateBounds(r io.Reader) {
	present, _ := core.ReadUInt8(r)

	if present&1 != 0 {
		readOrderCoord(r, &b.left, false)
	} else if present&16 != 0 {
		readOrderCoord(r, &b.left, true)
	}

	if present&2 != 0 {
		readOrderCoord(r, &b.top, false)
	} else if present&32 != 0 {
		readOrderCoord(r, &b.top, true)
	}

	if present&4 != 0 {
		readOrderCoord(r, &b.right, false)
	} else if present&64 != 0 {
		readOrderCoord(r, &b.right, true)
	}
	if present&8 != 0 {
		readOrderCoord(r, &b.bottom, false)
	} else if present&128 != 0 {
		readOrderCoord(r, &b.bottom, true)
	}
}

type PrimaryOrder interface {
	Type() int
	Unpack(io.Reader, uint32, bool) error
}

var (
	orderType uint8
	bounds    Bounds
)

func (o *OrderPdu) processPrimaryOrder(r io.Reader) error {
	o.Primary = &Primary{}
	if o.ControlFlags&TS_TYPE_CHANGE != 0 {
		orderType, _ = core.ReadUInt8(r)
	}
	size := 1
	switch orderType {
	case ORDER_TYPE_MEM3BLT, ORDER_TYPE_TEXT2:
		size = 3

	case ORDER_TYPE_PATBLT, ORDER_TYPE_MEMBLT, ORDER_TYPE_LINETO, ORDER_TYPE_POLYGON_CB, ORDER_TYPE_ELLIPSE_CB:
		size = 2
	}

	if o.ControlFlags&TS_ZERO_FIELD_BYTE_BIT0 != 0 {
		size--
	}
	if o.ControlFlags&TS_ZERO_FIELD_BYTE_BIT1 != 0 {
		if size < 2 {
			size = 0
		} else {
			size -= 2
		}
	}
	var present uint32
	for i := 0; i < size; i++ {
		bits, _ := core.ReadUInt8(r)
		present |= uint32(bits) << uint32(i*8)
	}

	if o.ControlFlags&TS_BOUNDS != 0 {
		if o.ControlFlags&TS_ZERO_BOUNDS_DELTAS == 0 {
			bounds.updateBounds(r)
		}
		//glog.Infof("updateBounds")
		o.Primary.Bounds = bounds
	}

	delta := o.ControlFlags&TS_DELTA_COORDINATES != 0

	//glog.Infof("present=%d,delta=%v", present, delta)

	var p PrimaryOrder
	switch orderType {
	case ORDER_TYPE_DSTBLT:
		p = &Dstblt{}

	case ORDER_TYPE_PATBLT:
		p = &Patblt{}

	case ORDER_TYPE_SCRBLT:
		p = &Scrblt{}

	//case ORDER_TYPE_DRAWNINEGRID:

	//case ORDER_TYPE_MULTI_DRAWNINEGRID:

	case ORDER_TYPE_LINETO:
		p = &LineTo{}

	case ORDER_TYPE_OPAQUERECT:
		p = &OpaqueRect{}

	case ORDER_TYPE_SAVEBITMAP:
		p = &SaveBitmap{}

	case ORDER_TYPE_MEMBLT:
		p = &Memblt{}

	case ORDER_TYPE_MEM3BLT:
		p = &Mem3blt{}

	//case ORDER_TYPE_MULTIDSTBLT:

	//case ORDER_TYPE_MULTIPATBLT:

	//case ORDER_TYPE_MULTISCRBLT:

	//case ORDER_TYPE_MULTIOPAQUERECT:

	//case ORDER_TYPE_FAST_INDEX:

	case ORDER_TYPE_POLYGON_SC:
		p = &PolygonSc{}

	case ORDER_TYPE_POLYGON_CB:
		p = &PolygonCb{}

	case ORDER_TYPE_POLYLINE:
		p = &Polyline{}

	//case ORDER_TYPE_FAST_GLYPH:

	case ORDER_TYPE_ELLIPSE_SC:
		p = &EllipeSc{}

	case ORDER_TYPE_ELLIPSE_CB:
		p = &EllipeCb{}

	case ORDER_TYPE_TEXT2:
		p = &GlayphIndex{}
	default:
		glog.Error("Not Support order type:", orderType)
		return errors.New("Not Support order type")
	}
	if p != nil {
		if err := p.Unpack(r, present, delta); err != nil {
			return err
		}
	}

	o.Primary.Data = p
	return nil
}
func readOrderCoord(r io.Reader, coord *int32, delta bool) {
	if delta {
		change, _ := core.ReadUInt8(r)
		*coord += int32(int8(change))
	} else {
		change, _ := core.ReadUint16LE(r)
		*coord = int32(int16(change))
	}
}

type Dstblt struct {
	x      int32
	y      int32
	cx     int32
	cy     int32
	opcode uint8
}

func (d *Dstblt) Type() int {
	return ORDER_TYPE_DSTBLT
}
func (d *Dstblt) Unpack(r io.Reader, present uint32, delta bool) error {
	glog.Infof("Dstblt Order")
	if present&0x01 != 0 {
		readOrderCoord(r, &d.x, delta)
	}
	if present&0x02 != 0 {
		readOrderCoord(r, &d.y, delta)
	}
	if present&0x04 != 0 {
		readOrderCoord(r, &d.cx, delta)
	}
	if present&0x08 != 0 {
		readOrderCoord(r, &d.cy, delta)
	}
	if present&0x10 != 0 {
		d.opcode, _ = core.ReadUInt8(r)
	}
	return nil
}

type Patblt struct {
	x        int32
	y        int32
	cx       int32
	cy       int32
	opcode   uint8
	bgcolour [4]uint8
	fgcolour [4]uint8
	brush    Brush
}

func (d *Patblt) Type() int {
	return ORDER_TYPE_PATBLT
}
func (d *Patblt) Unpack(r io.Reader, present uint32, delta bool) error {
	glog.Infof("Patblt Order")
	if present&0x01 != 0 {
		readOrderCoord(r, &d.x, delta)
	}
	if present&0x02 != 0 {
		readOrderCoord(r, &d.y, delta)
	}
	if present&0x04 != 0 {
		readOrderCoord(r, &d.cx, delta)
	}
	if present&0x08 != 0 {
		readOrderCoord(r, &d.cy, delta)
	}
	if present&0x10 != 0 {
		d.opcode, _ = core.ReadUInt8(r)
	}
	if present&0x0020 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.bgcolour[0], d.bgcolour[1], d.bgcolour[2], d.bgcolour[3] = b, g, r, a
	}
	if present&0x0040 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.fgcolour[0], d.fgcolour[1], d.fgcolour[2], d.fgcolour[3] = b, g, r, a
	}
	d.brush.updateBrush(r, present>>7)

	return nil
}

type Brush struct {
	X     uint8
	Y     uint8
	Style uint8
	Hatch uint8
	Data  []byte
}

func (b *Brush) updateBrush(r io.Reader, present uint32) {
	if present&1 != 0 {
		b.X, _ = core.ReadUInt8(r)
	}

	if present&2 != 0 {
		b.Y, _ = core.ReadUInt8(r)
	}

	if present&4 != 0 {
		b.Style, _ = core.ReadUInt8(r)
	}

	if present&8 != 0 {
		b.Hatch, _ = core.ReadUInt8(r)
	}

	if present&16 != 0 {
		data, _ := core.ReadBytes(7, r)
		b.Data = make([]byte, 0, 8)
		b.Data = append(b.Data, b.Hatch)
		b.Data = append(b.Data, data...)
	}
}

type Scrblt struct {
	X      int32
	Y      int32
	Cx     int32
	Cy     int32
	Opcode uint8
	Srcx   int32
	Srcy   int32
}

func (d *Scrblt) Type() int {
	return ORDER_TYPE_SCRBLT
}

var d Scrblt

func (d1 *Scrblt) Unpack(r io.Reader, present uint32, delta bool) error {
	glog.Infof("Scrblt Order")
	if present&0x0001 != 0 {
		readOrderCoord(r, &d.X, delta)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.Y, delta)
	}
	if present&0x0004 != 0 {
		readOrderCoord(r, &d.Cx, delta)
	}
	if present&0x0008 != 0 {
		readOrderCoord(r, &d.Cy, delta)
	}
	if present&0x0010 != 0 {
		d.Opcode, _ = core.ReadUInt8(r)
	}
	if present&0x0020 != 0 {
		readOrderCoord(r, &d.Srcx, delta)
	}
	if present&0x0040 != 0 {
		readOrderCoord(r, &d.Srcy, delta)
	}
	*d1 = d
	return nil
}

type LineTo struct {
	Mixmode  uint16
	Startx   int32
	Starty   int32
	Endx     int32
	Endy     int32
	Bgcolour [4]uint8
	Opcode   uint8
	Pen      Pen
}

func (d *LineTo) Type() int {
	return ORDER_TYPE_LINETO
}
func (d *LineTo) Unpack(r io.Reader, present uint32, delta bool) error {
	glog.Infof("LineTo Order")
	if present&0x0001 != 0 {
		d.Mixmode, _ = core.ReadUint16LE(r)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.Startx, delta)
	}
	if present&0x0004 != 0 {
		readOrderCoord(r, &d.Starty, delta)
	}
	if present&0x008 != 0 {
		readOrderCoord(r, &d.Endx, delta)
	}
	if present&0x0010 != 0 {
		readOrderCoord(r, &d.Endy, delta)
	}
	if present&0x0020 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.Bgcolour[0], d.Bgcolour[1], d.Bgcolour[2], d.Bgcolour[3] = b, g, r, a
	}
	if present&0x0040 != 0 {
		d.Opcode, _ = core.ReadUInt8(r)
	}

	d.Pen.updatePen(r, present>>7)

	return nil
}

type Pen struct {
	Style  uint8
	Width  uint8
	Colour [4]uint8
}

func (d *Pen) updatePen(r io.Reader, present uint32) {
	if present&1 != 0 {
		d.Style, _ = core.ReadUInt8(r)
	}

	if present&2 != 0 {
		d.Width, _ = core.ReadUInt8(r)
	}

	if present&4 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.Colour[0], d.Colour[1], d.Colour[2], d.Colour[3] = b, g, r, a
	}
}

type OpaqueRect struct {
	X      int32
	Y      int32
	Cx     int32
	Cy     int32
	Colour [4]uint8
}

func (d *OpaqueRect) Type() int {
	return ORDER_TYPE_OPAQUERECT
}
func (d *OpaqueRect) Unpack(r io.Reader, present uint32, delta bool) error {
	glog.Infof("OpaqueRect Order")
	if present&0x0001 != 0 {
		readOrderCoord(r, &d.X, delta)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.Y, delta)
	}
	if present&0x0004 != 0 {
		readOrderCoord(r, &d.Cx, delta)
	}
	if present&0x0008 != 0 {
		readOrderCoord(r, &d.Cy, delta)
	}
	if present&0x0010 != 0 {
		i, _ := core.ReadUInt8(r)
		d.Colour[0] = i
	}
	if present&0x0020 != 0 {
		i, _ := core.ReadUInt8(r)
		d.Colour[1] = i
	}
	if present&0x0040 != 0 {
		i, _ := core.ReadUInt8(r)
		d.Colour[2] = i
	}
	return nil
}

type SaveBitmap struct {
	Offset uint32
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
	action uint8
}

func (d *SaveBitmap) Type() int {
	return ORDER_TYPE_SAVEBITMAP
}
func (d *SaveBitmap) Unpack(r io.Reader, present uint32, delta bool) error {
	if present&0x0001 != 0 {
		d.Offset, _ = core.ReadUInt32LE(r)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.Left, delta)
	}
	if present&0x0004 != 0 {
		readOrderCoord(r, &d.Top, delta)
	}
	if present&0x0008 != 0 {
		readOrderCoord(r, &d.Right, delta)
	}
	if present&0x0010 != 0 {
		readOrderCoord(r, &d.Bottom, delta)
	}
	if present&0x0020 != 0 {
		d.action, _ = core.ReadUInt8(r)
	}
	return nil
}

type Memblt struct {
	ColourTable uint8
	CacheId     uint8
	X           int32
	Y           int32
	Cx          int32
	Cy          int32
	Opcode      uint8
	Srcx        int32
	Srcy        int32
	CacheIdx    uint16
}

func (d *Memblt) Type() int {
	return ORDER_TYPE_MEMBLT
}
func (d *Memblt) Unpack(r io.Reader, present uint32, delta bool) error {
	if present&0x0001 != 0 {
		d.CacheId, _ = core.ReadUInt8(r)
		d.ColourTable, _ = core.ReadUInt8(r)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.X, delta)
	}
	if present&0x0004 != 0 {
		readOrderCoord(r, &d.Y, delta)
	}
	if present&0x0008 != 0 {
		readOrderCoord(r, &d.Cx, delta)
	}
	if present&0x0010 != 0 {
		readOrderCoord(r, &d.Cy, delta)
	}
	if present&0x0020 != 0 {
		d.Opcode, _ = core.ReadUInt8(r)
	}
	if present&0x0040 != 0 {
		readOrderCoord(r, &d.Srcx, delta)
	}
	if present&0x0080 != 0 {
		readOrderCoord(r, &d.Srcy, delta)
	}
	if present&0x0100 != 0 {
		d.CacheIdx, _ = core.ReadUint16LE(r)
	}
	return nil
}

type Mem3blt struct {
	ColourTable uint8
	CacheId     uint8
	X           int32
	Y           int32
	Cx          int32
	Cy          int32
	Opcode      uint8
	Srcx        int32
	Srcy        int32
	Bgcolour    [4]uint8
	Fgcolour    [4]uint8
	Brush       Brush
	CacheIdx    uint16
}

func (d *Mem3blt) Type() int {
	return ORDER_TYPE_MEM3BLT
}
func (d *Mem3blt) Unpack(r io.Reader, present uint32, delta bool) error {
	if present&0x000001 != 0 {
		d.CacheId, _ = core.ReadUInt8(r)
		d.ColourTable, _ = core.ReadUInt8(r)
	}
	if present&0x000002 != 0 {
		readOrderCoord(r, &d.X, delta)
	}
	if present&0x000004 != 0 {
		readOrderCoord(r, &d.Y, delta)
	}
	if present&0x000008 != 0 {
		readOrderCoord(r, &d.Cx, delta)
	}
	if present&0x000010 != 0 {
		readOrderCoord(r, &d.Cy, delta)
	}
	if present&0x000020 != 0 {
		d.Opcode, _ = core.ReadUInt8(r)
	}
	if present&0x000040 != 0 {
		readOrderCoord(r, &d.Srcx, delta)
	}
	if present&0x000080 != 0 {
		readOrderCoord(r, &d.Srcy, delta)
	}
	if present&0x000100 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.Bgcolour[0], d.Bgcolour[1], d.Bgcolour[2], d.Bgcolour[3] = b, g, r, a
	}
	if present&0x000200 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.Fgcolour[0], d.Fgcolour[1], d.Fgcolour[2], d.Fgcolour[3] = b, g, r, a
	}
	d.Brush.updateBrush(r, present>>10)
	if present&0x008000 != 0 {
		d.CacheIdx, _ = core.ReadUint16LE(r)
	}
	if present&0x010000 != 0 {
		core.ReadUint16LE(r)
	}

	return nil
}

type PolygonSc struct {
	X        int32
	Y        int32
	Opcode   uint8
	Fillmode uint8
	Fgcolour [4]uint8
	Npoints  uint8
	Points   []Point
}

type Point struct {
	X int32
	Y int32
}

func (d *PolygonSc) Type() int {
	return ORDER_TYPE_POLYGON_SC
}
func (d *PolygonSc) Unpack(r io.Reader, present uint32, delta bool) error {
	if present&0x0001 != 0 {
		readOrderCoord(r, &d.X, delta)
	}
	if present&0x0002 != 0 {
		readOrderCoord(r, &d.Y, delta)
	}
	if present&0x0004 != 0 {
		d.Opcode, _ = core.ReadUInt8(r)
	}
	if present&0x0008 != 0 {
		d.Fillmode, _ = core.ReadUInt8(r)
	}
	if present&0x0010 != 0 {
		b, g, r, a := updateReadColorRef(r)
		d.Fgcolour[0], d.Fgcolour[1], d.Fgcolour[2], d.Fgcolour[3] = b, g, r, a
	}
	if present&0x0020 != 0 {
		d.Npoints, _ = core.ReadUInt8(r)
		d.Points = make([]Point, 0, d.Npoints+1)
	}
	if present&0x0040 != 0 {
		size, _ := core.ReadUInt8(r)
		data, _ := core.ReadBytes(int(size), r)
		d.Points = append(d.Points, Point{d.X, d.Y})
		var flags uint8
		r = bytes.NewReader(data)
		for i := 1; i <= int(d.Npoints); i++ {
			var p Point
			if (i-1)%4 == 0 {
				flags, _ = core.ReadUInt8(r)
			}
			if (^flags)&0x80 != 0 {
				p.X = parseDelta(r)
			}
			if (^flags)&0x40 != 0 {
				p.Y = parseDelta(r)
			}
			flags <<= 2
		}
	}

	return nil
}

func parseDelta(r io.Reader) (v int32) {
	b, _ := core.ReadUInt8(r)
	if b&0x40 != 0 {
		v = int32(b) | (^0x3F)
	} else {
		v = int32(b & 0x3F)
	}
	if b&0x80 != 0 {
		b, _ := core.ReadUInt8(r)
		v = (v << 8) | int32(b)
	}
	return
}

type PolygonCb struct {
}

func (d *PolygonCb) Type() int {
	return ORDER_TYPE_POLYGON_CB
}
func (d *PolygonCb) Unpack(r io.Reader, present uint32, delta bool) error {
	return nil
}

type Polyline struct {
}

func (d *Polyline) Type() int {
	return ORDER_TYPE_POLYLINE
}
func (d *Polyline) Unpack(r io.Reader, present uint32, delta bool) error {
	return nil
}

type EllipeSc struct {
}

func (d *EllipeSc) Type() int {
	return ORDER_TYPE_ELLIPSE_SC
}
func (d *EllipeSc) Unpack(r io.Reader, present uint32, delta bool) error {
	return nil
}

type EllipeCb struct {
}

func (d *EllipeCb) Type() int {
	return ORDER_TYPE_ELLIPSE_CB
}
func (d *EllipeCb) Unpack(r io.Reader, present uint32, delta bool) error {
	return nil
}

type GlayphIndex struct {
}

func (d *GlayphIndex) Type() int {
	return ORDER_TYPE_TEXT2
}
func (d *GlayphIndex) Unpack(r io.Reader, present uint32, delta bool) error {
	return nil
}

/*Secondary*/
func (s *Secondary) updateCacheBitmapOrder(r io.Reader, compressed bool, flags uint16) {
	var cb CacheBitmapOrder
	cb.cacheId, _ = core.ReadUInt8(r)
	core.ReadUInt8(r)
	cb.bitmapWidth, _ = core.ReadUInt8(r)
	cb.bitmapHeight, _ = core.ReadUInt8(r)
	cb.bitmapBpp, _ = core.ReadUInt8(r)
	bitmapLength, _ := core.ReadUint16LE(r)
	cb.cacheIndex, _ = core.ReadUint16LE(r)
	var bitmapComprHdr []byte
	if compressed {
		if (flags & NO_BITMAP_COMPRESSION_HDR) == 0 {
			bitmapComprHdr, _ = core.ReadBytes(8, r)
			bitmapLength -= 8
		}
	}
	cb.bitmapComprHdr = bitmapComprHdr
	cb.bitmapDataStream, _ = core.ReadBytes(int(bitmapLength), r)
	cb.bitmapLength = bitmapLength

}

type CacheBitmapOrder struct {
	cacheId          uint8
	bitmapBpp        uint8
	bitmapWidth      uint8
	bitmapHeight     uint8
	bitmapLength     uint16
	cacheIndex       uint16
	bitmapComprHdr   []byte
	bitmapDataStream []byte
}

func getCbV2Bpp(bpp uint32) (b uint32) {
	switch bpp {
	case 3:
		b = 8
	case 4:
		b = 16
	case 5:
		b = 24
	case 6:
		b = 32
	default:
		b = 0
	}
	return
}

type CacheBitmapV2Order struct {
	cacheId            uint32
	flags              uint32
	key1               uint32
	key2               uint32
	bitmapBpp          uint32
	bitmapWidth        uint8
	bitmapHeight       uint8
	bitmapLength       uint16
	cacheIndex         uint32
	compressed         bool
	cbCompFirstRowSize uint16
	cbCompMainBodySize uint16
	cbScanWidth        uint16
	cbUncompressedSize uint16
	bitmapDataStream   []byte
}

func (s *Secondary) updateCacheBitmapV2Order(r io.Reader, compressed bool, flags uint16) {
	var cb CacheBitmapV2Order
	cb.cacheId = uint32(flags) & 0x0003
	cb.flags = (uint32(flags) & 0xFF80) >> 7
	bitsPerPixelId := (uint32(flags) & 0x0078) >> 3
	cb.bitmapBpp = getCbV2Bpp(bitsPerPixelId)

	if cb.flags&CBR2_PERSISTENT_KEY_PRESENT != 0 {
		cb.key1, _ = core.ReadUInt32LE(r)
		cb.key2, _ = core.ReadUInt32LE(r)
	}

	if cb.flags&CBR2_HEIGHT_SAME_AS_WIDTH != 0 {
		cb.bitmapWidth, _ = core.ReadUInt8(r)
		cb.bitmapHeight = cb.bitmapWidth
	} else {
		cb.bitmapWidth, _ = core.ReadUInt8(r)
		cb.bitmapHeight, _ = core.ReadUInt8(r)
	}

	bitmapLength, _ := core.ReadUint16LE(r)
	cacheIndex, _ := core.ReadUInt8(r)

	if cb.flags&CBR2_DO_NOT_CACHE != 0 {
		cb.cacheIndex = 0x7FFF
	} else {
		cb.cacheIndex = uint32(cacheIndex)
	}

	if compressed {
		if cb.flags&CBR2_NO_BITMAP_COMPRESSION_HDR == 0 {
			cb.cbCompFirstRowSize, _ = core.ReadUint16LE(r)
			cb.cbCompMainBodySize, _ = core.ReadUint16LE(r)
			cb.cbScanWidth, _ = core.ReadUint16LE(r)
			cb.cbUncompressedSize, _ = core.ReadUint16LE(r)
			bitmapLength = cb.cbCompMainBodySize
		}
	}

	cb.bitmapDataStream, _ = core.ReadBytes(int(bitmapLength), r)
	cb.bitmapLength = bitmapLength
	cb.compressed = compressed

}

type CacheBitmapV3Order struct {
	cacheId    uint32
	bpp        uint32
	flags      uint32
	cacheIndex uint16
	key1       uint32
	key2       uint32
	bitmapData BitmapDataEx
}
type BitmapDataEx struct {
	bpp     uint8
	codecID uint8
	width   uint16
	height  uint16
	length  uint32
	data    []byte
}

func (s *Secondary) updateCacheBitmapV3Order(r io.Reader, flags uint16) {
	var cb CacheBitmapV3Order

	cb.cacheId = uint32(flags) & 0x00000003
	cb.flags = (uint32(flags) & 0x0000FF80) >> 7
	bitsPerPixelId := (uint32(flags) & 0x00000078) >> 3
	cb.bpp = getCbV2Bpp(bitsPerPixelId)

	cacheIndex, _ := core.ReadUint16LE(r)
	cb.cacheIndex = cacheIndex
	cb.key1, _ = core.ReadUInt32LE(r)
	cb.key2, _ = core.ReadUInt32LE(r)

	bitmapData := &cb.bitmapData
	bitmapData.bpp, _ = core.ReadUInt8(r)
	core.ReadUInt8(r)
	core.ReadUInt8(r)
	bitmapData.codecID, _ = core.ReadUInt8(r)
	bitmapData.width, _ = core.ReadUint16LE(r)
	bitmapData.height, _ = core.ReadUint16LE(r)
	new_len, _ := core.ReadUInt32LE(r)

	bitmapData.data, _ = core.ReadBytes(int(new_len), r)
	bitmapData.length = new_len

}

type CacheColorTableOrder struct {
	cacheIndex   uint8
	numberColors uint16
	colorTable   [256 * 4]uint8
}

func (s *Secondary) updateCacheColorTableOrder(r io.Reader, flags uint16) {
	var cb CacheColorTableOrder
	cb.cacheIndex, _ = core.ReadUInt8(r)
	cb.numberColors, _ = core.ReadUint16LE(r)

	if cb.numberColors != 256 {
		/* This field MUST be set to 256 */
		return
	}

	for i := 0; i < int(cb.numberColors)*4; i++ {
		cb.colorTable[i], cb.colorTable[i+1], cb.colorTable[i+2], cb.colorTable[i+3] = updateReadColorRef(r)
	}
}
func updateReadColorRef(r io.Reader) (uint8, uint8, uint8, uint8) {
	blue, _ := core.ReadUInt8(r)
	green, _ := core.ReadUInt8(r)
	red, _ := core.ReadUInt8(r)
	core.ReadUInt8(r)

	return blue, green, red, 255
}

type CacheGlyphOrder struct {
	cacheId uint8
	nglyphs uint8
	glyphs  []CacheGlyph
}
type CacheGlyph struct {
	character uint16
	offset    uint16
	baseline  uint16
	width     uint16
	height    uint16
	datasize  int
	data      []uint8
}

func (s *Secondary) updateCacheGlyphOrder(r io.Reader, flags uint16) {
	var cb CacheGlyphOrder

	cb.cacheId, _ = core.ReadUInt8(r)
	cb.nglyphs, _ = core.ReadUInt8(r)
	cb.glyphs = make([]CacheGlyph, 0, cb.nglyphs)

	for i := 0; i < int(cb.nglyphs); i++ {
		var c CacheGlyph
		c.character, _ = core.ReadUint16LE(r)
		c.offset, _ = core.ReadUint16LE(r)
		c.baseline, _ = core.ReadUint16LE(r)
		c.width, _ = core.ReadUint16LE(r)
		c.height, _ = core.ReadUint16LE(r)

		c.datasize = int(c.height*((c.width+7)/8)+3) & ^3
		c.data, _ = core.ReadBytes(c.datasize, r)

		cb.glyphs = append(cb.glyphs, c)
	}
}

type CacheBrushOrder struct {
	index  uint8
	bpp    uint8
	cx     uint8
	cy     uint8
	style  uint8
	length uint8
	data   []uint8
}

func (s *Secondary) updateCacheBrushOrder(r io.Reader, flags uint16) {
	var cb CacheBrushOrder
	cb.index, _ = core.ReadUInt8(r)
	cb.bpp, _ = core.ReadUInt8(r)
	cb.cx, _ = core.ReadUInt8(r)
	cb.cy, _ = core.ReadUInt8(r)
	cb.style, _ = core.ReadUInt8(r)
	cb.length, _ = core.ReadUInt8(r)
	if cb.cx == 8 && cb.cy == 8 {
		if cb.bpp == 1 {
			for i := 7; i >= 0; i-- {
				cb.data[i], _ = core.ReadUInt8(r)
			}
		} else {
			bpp := int(cb.bpp) - 2
			if int(cb.length) == 16+4*bpp {
				/* compressed brush */
				data, _ := core.ReadBytes(int(cb.length), r)
				cb.data = update_decompress_brush(data, bpp)
			} else {
				/* uncompressed brush */
				scanline := 8 * 8 * bpp
				cb.data, _ = core.ReadBytes(scanline, r)
			}
		}
	}
}
func update_decompress_brush(in []uint8, bpp int) []uint8 {
	var pal_index, in_index, shift int

	pal := in[16:]
	out := make([]uint8, 8*8*bpp)
	/* read it bottom up */
	for y := 7; y >= 0; y-- {
		/* 2 bytes per row */
		x := 0
		for do2 := 0; do2 < 2; do2++ {
			/* 4 pixels per byte */
			shift = 6
			for shift >= 0 {
				pal_index = int((in[in_index] >> uint(shift)) & 3)
				/* size of palette entries depends on bpp */
				for i := 0; i < bpp; i++ {
					out[(y*8+x)*bpp+i] = pal[pal_index*bpp+i]
				}
				x++
				shift -= 2
			}
			in_index++
		}
	}

	return out
}

/*Primary*/
type Bounds struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}
type OrderInfo struct {
	controlFlags     uint32
	orderType        uint32
	fieldFlags       uint32
	boundsFlags      uint32
	bounds           Bounds
	deltaCoordinates bool
}
