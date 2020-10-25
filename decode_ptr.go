package json

import (
	"unsafe"
)

type ptrDecoder struct {
	dec decoder
	typ *rtype
}

func newPtrDecoder(dec decoder, typ *rtype) *ptrDecoder {
	return &ptrDecoder{dec: dec, typ: typ}
}

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(*rtype) unsafe.Pointer

func (d *ptrDecoder) decodeStream(s *stream, p unsafe.Pointer) error {
	newptr := unsafe_New(d.typ)
	if err := d.dec.decodeStream(s, unsafe.Pointer(newptr)); err != nil {
		return err
	}
	*(*uintptr)(p) = uintptr(newptr)
	return nil
}

func (d *ptrDecoder) decode(buf []byte, cursor int64, p unsafe.Pointer) (int64, error) {
	newptr := unsafe_New(d.typ)
	c, err := d.dec.decode(buf, cursor, unsafe.Pointer(newptr))
	if err != nil {
		return 0, err
	}
	cursor = c
	*(*uintptr)(p) = uintptr(newptr)
	return cursor, nil
}
