package win

import (
	"unicode/utf16"
	"unsafe"
)

func unicode16FromString(s string) []uint16 {
	r := make([]rune, 0)
	for _, c := range s {
		r = append(r, c)
	}
	b := utf16.Encode(r)
	return append(b, uint16(0))
}

func stringFromUnicode16(s *uint16) string {
	if s == nil {
		return ""
	}
	buffer := []uint16{}
	ptr := uintptr(unsafe.Pointer(s))
	for true {
		ch := *(*uint16)(unsafe.Pointer(ptr))
		if ch == 0 {
			break
		}
		buffer = append(buffer, ch)
		ptr += unsafe.Sizeof(ch)
	}
	return string(utf16.Decode(buffer))
}
