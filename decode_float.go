package json

import (
	"strconv"
	"unsafe"
)

type floatDecoder struct {
	op func(unsafe.Pointer, float64)
}

func newFloatDecoder(op func(unsafe.Pointer, float64)) *floatDecoder {
	return &floatDecoder{op: op}
}

var (
	floatTable = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'.': true,
		'e': true,
		'E': true,
		'+': true,
		'-': true,
	}

	validEndNumberChar = [256]bool{
		nul:  true,
		' ':  true,
		'\t': true,
		'\r': true,
		'\n': true,
		',':  true,
		':':  true,
		'}':  true,
		']':  true,
	}
)

func floatBytes(s *stream) []byte {
	start := s.cursor
	for {
		s.cursor++
		if floatTable[s.char()] {
			continue
		} else if s.char() == nul {
			if s.read() {
				s.cursor-- // for retry current character
				continue
			}
		}
		break
	}
	return s.buf[start:s.cursor]
}

func (d *floatDecoder) decodeStreamByte(s *stream) ([]byte, error) {
	for {
		switch s.char() {
		case ' ', '\n', '\t', '\r':
			s.cursor++
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return floatBytes(s), nil
		case nul:
			if s.read() {
				continue
			}
			goto ERROR
		default:
			goto ERROR
		}
	}
ERROR:
	return nil, errUnexpectedEndOfJSON("float", s.totalOffset())
}

func (d *floatDecoder) decodeByte(buf []byte, cursor int64) ([]byte, int64, error) {
	buflen := int64(len(buf))
	for ; cursor < buflen; cursor++ {
		switch buf[cursor] {
		case ' ', '\n', '\t', '\r':
			continue
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			start := cursor
			cursor++
			for ; cursor < buflen; cursor++ {
				if floatTable[buf[cursor]] {
					continue
				}
				break
			}
			num := buf[start:cursor]
			return num, cursor, nil
		default:
			return nil, 0, errUnexpectedEndOfJSON("float", cursor)
		}
	}
	return nil, 0, errUnexpectedEndOfJSON("float", cursor)
}

func (d *floatDecoder) decodeStream(s *stream, p unsafe.Pointer) error {
	bytes, err := d.decodeStreamByte(s)
	if err != nil {
		return err
	}
	if !validEndNumberChar[s.char()] {
		return errUnexpectedEndOfJSON("float", s.totalOffset())
	}
	str := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	d.op(p, f64)
	return nil
}

func (d *floatDecoder) decode(buf []byte, cursor int64, p unsafe.Pointer) (int64, error) {
	bytes, c, err := d.decodeByte(buf, cursor)
	if err != nil {
		return 0, err
	}
	cursor = c
	if !validEndNumberChar[buf[cursor]] {
		return 0, errUnexpectedEndOfJSON("float", cursor)
	}
	s := *(*string)(unsafe.Pointer(&bytes))
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	d.op(p, f64)
	return cursor, nil
}
