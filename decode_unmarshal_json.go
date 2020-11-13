package json

import (
	"unsafe"
)

type unmarshalJSONDecoder struct {
	typ *rtype
}

func newUnmarshalJSONDecoder(typ *rtype) *unmarshalJSONDecoder {
	return &unmarshalJSONDecoder{typ: typ}
}

func (d *unmarshalJSONDecoder) decodeStream(s *stream, p unsafe.Pointer) error {
	s.skipWhiteSpace()
	start := s.cursor
	if err := s.skipValue(); err != nil {
		return err
	}
	src := s.buf[start:s.cursor]
	v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
		typ: d.typ,
		ptr: p,
	}))
	if err := v.(Unmarshaler).UnmarshalJSON(src); err != nil {
		return err
	}
	return nil
}

func (d *unmarshalJSONDecoder) decode(buf []byte, cursor int64, p unsafe.Pointer) (int64, error) {
	cursor = skipWhiteSpace(buf, cursor)
	start := cursor
	end, err := skipValue(buf, cursor)
	if err != nil {
		return 0, err
	}
	src := buf[start:end]
	v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
		typ: d.typ,
		ptr: p,
	}))
	if err := v.(Unmarshaler).UnmarshalJSON(src); err != nil {
		return 0, err
	}
	return end, nil
}
