package json

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"unsafe"
)

const startDetectingCyclesAfter = 1000

func load(base uintptr, idx uintptr) uintptr {
	addr := base + idx
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func store(base uintptr, idx uintptr, p uintptr) {
	addr := base + idx
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func errUnsupportedValue(code *opcode, ptr uintptr) *UnsupportedValueError {
	v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
		typ: code.typ,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
	}))
	return &UnsupportedValueError{
		Value: reflect.ValueOf(v),
		Str:   fmt.Sprintf("encountered a cycle via %s", code.typ),
	}
}

func (e *Encoder) run(ctx *encodeRuntimeContext, code *opcode) error {
	recursiveLevel := 0
	seenPtr := map[uintptr]struct{}{}
	ptrOffset := uintptr(0)
	ctxptr := ctx.ptr()

	for {
		switch code.op {
		default:
			return fmt.Errorf("failed to handle opcode. doesn't implement %s", code.op)
		case opPtr, opPtrIndent:
			ptr := load(ctxptr, code.idx)
			code = code.next
			store(ctxptr, code.idx, e.ptrToPtr(ptr))
		case opInt:
			e.encodeInt(e.ptrToInt(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opIntIndent:
			e.encodeInt(e.ptrToInt(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opInt8:
			e.encodeInt8(e.ptrToInt8(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opInt8Indent:
			e.encodeInt8(e.ptrToInt8(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opInt16:
			e.encodeInt16(e.ptrToInt16(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opInt16Indent:
			e.encodeInt16(e.ptrToInt16(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opInt32:
			e.encodeInt32(e.ptrToInt32(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opInt32Indent:
			e.encodeInt32(e.ptrToInt32(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opInt64:
			e.encodeInt64(e.ptrToInt64(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opInt64Indent:
			e.encodeInt64(e.ptrToInt64(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opUint:
			e.encodeUint(e.ptrToUint(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opUintIndent:
			e.encodeUint(e.ptrToUint(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opUint8:
			e.encodeUint8(e.ptrToUint8(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opUint8Indent:
			e.encodeUint8(e.ptrToUint8(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opUint16:
			e.encodeUint16(e.ptrToUint16(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opUint16Indent:
			e.encodeUint16(e.ptrToUint16(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opUint32:
			e.encodeUint32(e.ptrToUint32(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opUint32Indent:
			e.encodeUint32(e.ptrToUint32(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opUint64:
			e.encodeUint64(e.ptrToUint64(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opUint64Indent:
			e.encodeUint64(e.ptrToUint64(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opFloat32:
			e.encodeFloat32(e.ptrToFloat32(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opFloat32Indent:
			e.encodeFloat32(e.ptrToFloat32(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opFloat64:
			v := e.ptrToFloat64(load(ctxptr, code.idx))
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeFloat64(v)
			e.encodeByte(',')
			code = code.next
		case opFloat64Indent:
			v := e.ptrToFloat64(load(ctxptr, code.idx))
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeFloat64(v)
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opString:
			e.encodeString(e.ptrToString(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opStringIndent:
			e.encodeString(e.ptrToString(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opBool:
			e.encodeBool(e.ptrToBool(load(ctxptr, code.idx)))
			e.encodeByte(',')
			code = code.next
		case opBoolIndent:
			e.encodeBool(e.ptrToBool(load(ctxptr, code.idx)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opBytes:
			ptr := load(ctxptr, code.idx)
			header := *(**sliceHeader)(unsafe.Pointer(&ptr))
			if ptr == 0 || uintptr(header.data) == 0 {
				e.encodeNull()
			} else {
				e.encodeByteSlice(e.ptrToBytes(ptr))
			}
			e.encodeByte(',')
			code = code.next
		case opBytesIndent:
			ptr := load(ctxptr, code.idx)
			header := *(**sliceHeader)(unsafe.Pointer(&ptr))
			if ptr == 0 || uintptr(header.data) == 0 {
				e.encodeNull()
			} else {
				e.encodeByteSlice(e.ptrToBytes(ptr))
			}
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opInterface:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.next
				break
			}
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
			}))
			if _, exists := seenPtr[ptr]; exists {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   fmt.Sprintf("encountered a cycle via %s", code.typ),
				}
			}
			seenPtr[ptr] = struct{}{}
			rv := reflect.ValueOf(v)
			if rv.IsNil() {
				e.encodeNull()
				e.encodeByte(',')
				code = code.next
				break
			}
			vv := rv.Interface()
			header := (*interfaceHeader)(unsafe.Pointer(&vv))
			typ := header.typ
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			var c *opcode
			if typ.Kind() == reflect.Map {
				code, err := e.compileMap(&encodeCompileContext{
					typ:        typ,
					root:       code.root,
					withIndent: e.enabledIndent,
					indent:     code.indent,
				}, false)
				if err != nil {
					return err
				}
				c = code
			} else {
				code, err := e.compile(&encodeCompileContext{
					typ:        typ,
					root:       code.root,
					withIndent: e.enabledIndent,
					indent:     code.indent,
				})
				if err != nil {
					return err
				}
				c = code
			}

			beforeLastCode := c.beforeLastCode()
			lastCode := beforeLastCode.next
			lastCode.idx = beforeLastCode.idx + uintptrSize
			totalLength := uintptr(code.totalLength())
			nextTotalLength := uintptr(c.totalLength())
			curlen := uintptr(len(ctx.ptrs))
			offsetNum := ptrOffset / uintptrSize
			oldOffset := ptrOffset
			ptrOffset += totalLength * uintptrSize

			newLen := offsetNum + totalLength + nextTotalLength
			if curlen < newLen {
				ctx.ptrs = append(ctx.ptrs, make([]uintptr, newLen-curlen)...)
			}
			ctxptr = ctx.ptr() + ptrOffset // assign new ctxptr

			store(ctxptr, 0, uintptr(header.ptr))
			store(ctxptr, lastCode.idx, oldOffset)

			// link lastCode ( opInterfaceEnd ) => code.next
			lastCode.op = opInterfaceEnd
			lastCode.next = code.next

			code = c
			recursiveLevel++
		case opInterfaceIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
				break
			}
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
			}))
			if _, exists := seenPtr[ptr]; exists {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   fmt.Sprintf("encountered a cycle via %s", code.typ),
				}
			}
			seenPtr[ptr] = struct{}{}
			rv := reflect.ValueOf(v)
			if rv.IsNil() {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
				break
			}
			vv := rv.Interface()
			header := (*interfaceHeader)(unsafe.Pointer(&vv))
			typ := header.typ
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			var c *opcode
			if typ.Kind() == reflect.Map {
				code, err := e.compileMap(&encodeCompileContext{
					typ:        typ,
					root:       code.root,
					withIndent: e.enabledIndent,
					indent:     code.indent,
				}, false)
				if err != nil {
					return err
				}
				c = code
			} else {
				code, err := e.compile(&encodeCompileContext{
					typ:        typ,
					root:       code.root,
					withIndent: e.enabledIndent,
					indent:     code.indent,
				})
				if err != nil {
					return err
				}
				c = code
			}

			beforeLastCode := c.beforeLastCode()
			lastCode := beforeLastCode.next
			lastCode.idx = beforeLastCode.idx + uintptrSize
			totalLength := uintptr(code.totalLength())
			nextTotalLength := uintptr(c.totalLength())
			curlen := uintptr(len(ctx.ptrs))
			offsetNum := ptrOffset / uintptrSize
			oldOffset := ptrOffset
			ptrOffset += totalLength * uintptrSize

			newLen := offsetNum + totalLength + nextTotalLength
			if curlen < newLen {
				ctx.ptrs = append(ctx.ptrs, make([]uintptr, newLen-curlen)...)
			}
			ctxptr = ctx.ptr() + ptrOffset // assign new ctxptr

			store(ctxptr, 0, uintptr(header.ptr))
			store(ctxptr, lastCode.idx, oldOffset)

			// link lastCode ( opInterfaceEnd ) => code.next
			lastCode.op = opInterfaceEnd
			lastCode.next = code.next

			code = c
			recursiveLevel++
		case opInterfaceEnd, opInterfaceEndIndent:
			recursiveLevel--
			// restore ctxptr
			offset := load(ctxptr, code.idx)
			ctxptr = ctx.ptr() + offset
			ptrOffset = offset
			code = code.next
		case opMarshalJSON:
			ptr := load(ctxptr, code.idx)
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			if len(b) == 0 {
				return errUnexpectedEndOfJSON(
					fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
					0,
				)
			}
			var buf bytes.Buffer
			if e.enabledIndent {
				if err := encodeWithIndent(
					&buf,
					b,
					string(e.prefix)+string(bytes.Repeat(e.indentStr, code.indent)),
					string(e.indentStr),
				); err != nil {
					return err
				}
			} else {
				if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
					return err
				}
			}
			e.encodeBytes(buf.Bytes())
			e.encodeByte(',')
			code = code.next
		case opMarshalJSONIndent:
			ptr := load(ctxptr, code.idx)
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			if len(b) == 0 {
				return errUnexpectedEndOfJSON(
					fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
					0,
				)
			}
			var buf bytes.Buffer
			if e.enabledIndent {
				if err := encodeWithIndent(
					&buf,
					b,
					string(e.prefix)+string(bytes.Repeat(e.indentStr, code.indent)),
					string(e.indentStr),
				); err != nil {
					return err
				}
			} else {
				if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
					return err
				}
			}
			e.encodeBytes(buf.Bytes())
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opMarshalText:
			ptr := load(ctxptr, code.idx)
			isPtr := code.typ.Kind() == reflect.Ptr
			p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
			if p == nil {
				e.encodeNull()
				e.encodeByte(',')
			} else if isPtr && *(*unsafe.Pointer)(p) == nil {
				e.encodeBytes([]byte{'"', '"'})
				e.encodeByte(',')
			} else {
				if isPtr && code.typ.Elem().Implements(marshalTextType) {
					p = *(*unsafe.Pointer)(p)
				}
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: p,
				}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
			}
			code = code.next
		case opMarshalTextIndent:
			ptr := load(ctxptr, code.idx)
			isPtr := code.typ.Kind() == reflect.Ptr
			p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
			if p == nil {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
			} else if isPtr && *(*unsafe.Pointer)(p) == nil {
				e.encodeBytes([]byte{'"', '"'})
				e.encodeBytes([]byte{',', '\n'})
			} else {
				if isPtr && code.typ.Elem().Implements(marshalTextType) {
					p = *(*unsafe.Pointer)(p)
				}
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: p,
				}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opSliceHead:
			p := load(ctxptr, code.idx)
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				store(ctxptr, code.elemIdx, 0)
				store(ctxptr, code.length, uintptr(header.len))
				store(ctxptr, code.idx, uintptr(header.data))
				if header.len > 0 {
					e.encodeByte('[')
					code = code.next
					store(ctxptr, code.idx, uintptr(header.data))
				} else {
					e.encodeBytes([]byte{'[', ']', ','})
					code = code.end.next
				}
			}
		case opSliceElem:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if idx < length {
				store(ctxptr, code.elemIdx, idx)
				data := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, data+idx*size)
			} else {
				last := len(e.buf) - 1
				e.buf[last] = ']'
				e.encodeByte(',')
				code = code.end.next
			}
		case opSliceHeadIndent:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				header := *(**sliceHeader)(unsafe.Pointer(&p))
				store(ctxptr, code.elemIdx, 0)
				store(ctxptr, code.length, uintptr(header.len))
				store(ctxptr, code.idx, uintptr(header.data))
				if header.len > 0 {
					e.encodeBytes([]byte{'[', '\n'})
					e.encodeIndent(code.indent + 1)
					code = code.next
					store(ctxptr, code.idx, uintptr(header.data))
				} else {
					e.encodeIndent(code.indent)
					e.encodeBytes([]byte{'[', ']', '\n'})
					code = code.end.next
				}
			}
		case opRootSliceHeadIndent:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				header := *(**sliceHeader)(unsafe.Pointer(&p))
				store(ctxptr, code.elemIdx, 0)
				store(ctxptr, code.length, uintptr(header.len))
				store(ctxptr, code.idx, uintptr(header.data))
				if header.len > 0 {
					e.encodeBytes([]byte{'[', '\n'})
					e.encodeIndent(code.indent + 1)
					code = code.next
					store(ctxptr, code.idx, uintptr(header.data))
				} else {
					e.encodeIndent(code.indent)
					e.encodeBytes([]byte{'[', ']', ',', '\n'})
					code = code.end.next
				}
			}
		case opSliceElemIndent:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if idx < length {
				e.encodeIndent(code.indent + 1)
				store(ctxptr, code.elemIdx, idx)
				data := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, data+idx*size)
			} else {
				e.buf = e.buf[:len(e.buf)-2]
				e.encodeByte('\n')
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{']', ',', '\n'})
				code = code.end.next
			}
		case opRootSliceElemIndent:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if idx < length {
				e.encodeIndent(code.indent + 1)
				store(ctxptr, code.elemIdx, idx)
				code = code.next
				data := load(ctxptr, code.headIdx)
				store(ctxptr, code.idx, data+idx*code.size)
			} else {
				e.encodeByte('\n')
				e.encodeIndent(code.indent)
				e.encodeByte(']')
				code = code.end.next
			}
		case opArrayHead:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				if code.length > 0 {
					e.encodeByte('[')
					store(ctxptr, code.elemIdx, 0)
					code = code.next
					store(ctxptr, code.idx, p)
				} else {
					e.encodeBytes([]byte{'[', ']', ','})
					code = code.end.next
				}
			}
		case opArrayElem:
			idx := load(ctxptr, code.elemIdx)
			idx++
			if idx < code.length {
				store(ctxptr, code.elemIdx, idx)
				p := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, p+idx*size)
			} else {
				last := len(e.buf) - 1
				e.buf[last] = ']'
				e.encodeByte(',')
				code = code.end.next
			}
		case opArrayHeadIndent:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				if code.length > 0 {
					e.encodeBytes([]byte{'[', '\n'})
					e.encodeIndent(code.indent + 1)
					store(ctxptr, code.elemIdx, 0)
					code = code.next
					store(ctxptr, code.idx, p)
				} else {
					e.encodeIndent(code.indent)
					e.encodeBytes([]byte{'[', ']', ',', '\n'})
					code = code.end.next
				}
			}
		case opArrayElemIndent:
			idx := load(ctxptr, code.elemIdx)
			idx++
			if idx < code.length {
				e.encodeIndent(code.indent + 1)
				store(ctxptr, code.elemIdx, idx)
				p := load(ctxptr, code.headIdx)
				size := code.size
				code = code.next
				store(ctxptr, code.idx, p+idx*size)
			} else {
				e.buf = e.buf[:len(e.buf)-2]
				e.encodeByte('\n')
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{']', ',', '\n'})
				code = code.end.next
			}
		case opMapHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				uptr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				mlen := maplen(uptr)
				if mlen > 0 {
					e.encodeByte('{')
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))
					if !e.unorderedMap {
						pos := make([]int, 0, mlen)
						pos = append(pos, len(e.buf))
						posPtr := unsafe.Pointer(&pos)
						ctx.keepRefs = append(ctx.keepRefs, posPtr)
						store(ctxptr, code.end.mapPos, uintptr(posPtr))
					}
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
					code = code.end.next
				}
			}
		case opMapHeadLoad:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				// load pointer
				ptr = uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&ptr)))
				uptr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				if ptr == 0 {
					e.encodeNull()
					e.encodeByte(',')
					code = code.end.next
					break
				}
				mlen := maplen(uptr)
				if mlen > 0 {
					e.encodeByte('{')
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					if !e.unorderedMap {
						pos := make([]int, 0, mlen)
						pos = append(pos, len(e.buf))
						posPtr := unsafe.Pointer(&pos)
						ctx.keepRefs = append(ctx.keepRefs, posPtr)
						store(ctxptr, code.end.mapPos, uintptr(posPtr))
					}
					code = code.next
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
					code = code.end.next
				}
			}
		case opMapKey:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if e.unorderedMap {
				if idx < length {
					ptr := load(ctxptr, code.mapIter)
					iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
					store(ctxptr, code.elemIdx, idx)
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					last := len(e.buf) - 1
					e.buf[last] = '}'
					e.encodeByte(',')
					code = code.end.next
				}
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				posPtr := (*[]int)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
				*posPtr = append(*posPtr, len(e.buf))
				if idx < length {
					ptr := load(ctxptr, code.mapIter)
					iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
					store(ctxptr, code.elemIdx, idx)
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					code = code.end
				}
			}
		case opMapValue:
			if e.unorderedMap {
				last := len(e.buf) - 1
				e.buf[last] = ':'
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				posPtr := (*[]int)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
				*posPtr = append(*posPtr, len(e.buf))
			}
			ptr := load(ctxptr, code.mapIter)
			iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
			value := mapitervalue(iter)
			store(ctxptr, code.next.idx, uintptr(value))
			mapiternext(iter)
			code = code.next
		case opMapEnd:
			// this operation only used by sorted map.
			length := int(load(ctxptr, code.length))
			type mapKV struct {
				key   string
				value string
			}
			kvs := make([]mapKV, 0, length)
			ptr := load(ctxptr, code.mapPos)
			posPtr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
			pos := *(*[]int)(posPtr)
			for i := 0; i < length; i++ {
				startKey := pos[i*2]
				startValue := pos[i*2+1]
				var endValue int
				if i+1 < length {
					endValue = pos[i*2+2]
				} else {
					endValue = len(e.buf)
				}
				kvs = append(kvs, mapKV{
					key:   string(e.buf[startKey:startValue]),
					value: string(e.buf[startValue:endValue]),
				})
			}
			sort.Slice(kvs, func(i, j int) bool {
				return kvs[i].key < kvs[j].key
			})
			buf := e.buf[pos[0]:]
			buf = buf[:0]
			for _, kv := range kvs {
				buf = append(buf, []byte(kv.key)...)
				buf[len(buf)-1] = ':'
				buf = append(buf, []byte(kv.value)...)
			}
			buf[len(buf)-1] = '}'
			buf = append(buf, ',')
			e.buf = e.buf[:pos[0]]
			e.buf = append(e.buf, buf...)
			code = code.next
		case opMapHeadIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				uptr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				mlen := maplen(uptr)
				if mlen > 0 {
					e.encodeBytes([]byte{'{', '\n'})
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))

					if !e.unorderedMap {
						pos := make([]int, 0, mlen)
						pos = append(pos, len(e.buf))
						posPtr := unsafe.Pointer(&pos)
						ctx.keepRefs = append(ctx.keepRefs, posPtr)
						store(ctxptr, code.end.mapPos, uintptr(posPtr))
					} else {
						e.encodeIndent(code.next.indent)
					}

					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					e.encodeIndent(code.indent)
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					code = code.end.next
				}
			}
		case opMapHeadLoadIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				code = code.end.next
			} else {
				// load pointer
				ptr = uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&ptr)))
				uptr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				if uintptr(uptr) == 0 {
					e.encodeIndent(code.indent)
					e.encodeNull()
					e.encodeBytes([]byte{',', '\n'})
					code = code.end.next
					break
				}
				mlen := maplen(uptr)
				if mlen > 0 {
					e.encodeBytes([]byte{'{', '\n'})
					iter := mapiterinit(code.typ, uptr)
					ctx.keepRefs = append(ctx.keepRefs, iter)
					store(ctxptr, code.elemIdx, 0)
					store(ctxptr, code.length, uintptr(mlen))
					store(ctxptr, code.mapIter, uintptr(iter))
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))

					if !e.unorderedMap {
						pos := make([]int, 0, mlen)
						pos = append(pos, len(e.buf))
						posPtr := unsafe.Pointer(&pos)
						ctx.keepRefs = append(ctx.keepRefs, posPtr)
						store(ctxptr, code.end.mapPos, uintptr(posPtr))
					} else {
						e.encodeIndent(code.next.indent)
					}

					code = code.next
				} else {
					e.encodeIndent(code.indent)
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					code = code.end.next
				}
			}
		case opMapKeyIndent:
			idx := load(ctxptr, code.elemIdx)
			length := load(ctxptr, code.length)
			idx++
			if e.unorderedMap {
				if idx < length {
					e.encodeIndent(code.indent)
					store(ctxptr, code.elemIdx, idx)
					ptr := load(ctxptr, code.mapIter)
					iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					last := len(e.buf) - 1
					e.buf[last] = '\n'
					e.encodeIndent(code.indent - 1)
					e.encodeBytes([]byte{'}', ',', '\n'})
					code = code.end.next
				}
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				posPtr := (*[]int)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
				*posPtr = append(*posPtr, len(e.buf))
				if idx < length {
					ptr := load(ctxptr, code.mapIter)
					iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
					store(ctxptr, code.elemIdx, idx)
					key := mapiterkey(iter)
					store(ctxptr, code.next.idx, uintptr(key))
					code = code.next
				} else {
					code = code.end
				}
			}
		case opMapValueIndent:
			if e.unorderedMap {
				e.encodeBytes([]byte{':', ' '})
			} else {
				ptr := load(ctxptr, code.end.mapPos)
				posPtr := (*[]int)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
				*posPtr = append(*posPtr, len(e.buf))
			}
			ptr := load(ctxptr, code.mapIter)
			iter := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
			value := mapitervalue(iter)
			store(ctxptr, code.next.idx, uintptr(value))
			mapiternext(iter)
			code = code.next
		case opMapEndIndent:
			// this operation only used by sorted map
			length := int(load(ctxptr, code.length))
			type mapKV struct {
				key   string
				value string
			}
			kvs := make([]mapKV, 0, length)
			ptr := load(ctxptr, code.mapPos)
			pos := *(*[]int)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
			for i := 0; i < length; i++ {
				startKey := pos[i*2]
				startValue := pos[i*2+1]
				var endValue int
				if i+1 < length {
					endValue = pos[i*2+2]
				} else {
					endValue = len(e.buf)
				}
				kvs = append(kvs, mapKV{
					key:   string(e.buf[startKey:startValue]),
					value: string(e.buf[startValue:endValue]),
				})
			}
			sort.Slice(kvs, func(i, j int) bool {
				return kvs[i].key < kvs[j].key
			})
			buf := e.buf[pos[0]:]
			buf = buf[:0]
			for _, kv := range kvs {
				buf = append(buf, e.prefix...)
				buf = append(buf, bytes.Repeat(e.indentStr, code.indent+1)...)

				buf = append(buf, []byte(kv.key)...)
				buf[len(buf)-2] = ':'
				buf[len(buf)-1] = ' '
				buf = append(buf, []byte(kv.value)...)
			}
			buf = buf[:len(buf)-2]
			buf = append(buf, '\n')
			buf = append(buf, e.prefix...)
			buf = append(buf, bytes.Repeat(e.indentStr, code.indent)...)
			buf = append(buf, '}', ',', '\n')
			e.buf = e.buf[:pos[0]]
			e.buf = append(e.buf, buf...)
			code = code.next
		case opStructFieldPtrAnonymousHeadRecursive:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadRecursive:
			fallthrough
		case opStructFieldRecursive:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				if recursiveLevel > startDetectingCyclesAfter {
					if _, exists := seenPtr[ptr]; exists {
						return errUnsupportedValue(code, ptr)
					}
				}
			}
			seenPtr[ptr] = struct{}{}
			c := code.jmp.code
			c.end.next = newEndOp(&encodeCompileContext{})
			c.op = c.op.ptrHeadToHead()

			beforeLastCode := c.end
			lastCode := beforeLastCode.next

			lastCode.idx = beforeLastCode.idx + uintptrSize
			lastCode.elemIdx = lastCode.idx + uintptrSize

			// extend length to alloc slot for elemIdx
			totalLength := uintptr(code.totalLength() + 1)
			nextTotalLength := uintptr(c.totalLength() + 1)

			curlen := uintptr(len(ctx.ptrs))
			offsetNum := ptrOffset / uintptrSize
			oldOffset := ptrOffset
			ptrOffset += totalLength * uintptrSize

			newLen := offsetNum + totalLength + nextTotalLength
			if curlen < newLen {
				ctx.ptrs = append(ctx.ptrs, make([]uintptr, newLen-curlen)...)
			}
			ctxptr = ctx.ptr() + ptrOffset // assign new ctxptr

			store(ctxptr, c.idx, ptr)
			store(ctxptr, lastCode.idx, oldOffset)
			store(ctxptr, lastCode.elemIdx, uintptr(unsafe.Pointer(code.next)))

			// link lastCode ( opStructFieldRecursiveEnd ) => code.next
			lastCode.op = opStructFieldRecursiveEnd
			code = c
			recursiveLevel++
		case opStructFieldRecursiveEnd:
			recursiveLevel--

			// restore ctxptr
			offset := load(ctxptr, code.idx)
			ptr := load(ctxptr, code.elemIdx)
			code = (*opcode)(*(*unsafe.Pointer)(unsafe.Pointer(&ptr)))
			ctxptr = ctx.ptr() + offset
			ptrOffset = offset
		case opStructFieldPtrHead:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHead {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				if !code.anonymousKey {
					e.encodeKey(code)
				}
				p := ptr + code.offset
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldAnonymousHead:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				code = code.next
				store(ctxptr, code.idx, ptr)
			}
		case opStructFieldPtrHeadInt:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadInt {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeInt(e.ptrToInt(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadInt:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeInt(e.ptrToInt(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadInt8:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadInt8 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeInt8(e.ptrToInt8(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadInt8:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeInt8(e.ptrToInt8(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadInt16:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadInt16 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeInt16(e.ptrToInt16(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadInt16:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeInt16(e.ptrToInt16(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadInt32:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadInt32 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeInt32(e.ptrToInt32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadInt32:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeInt32(e.ptrToInt32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadInt64:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadInt64 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeInt64(e.ptrToInt64(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadInt64:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeInt64(e.ptrToInt64(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadUint:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadUint {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeUint(e.ptrToUint(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUint:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeUint(e.ptrToUint(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadUint8:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadUint8 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeUint8(e.ptrToUint8(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUint8:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeUint8(e.ptrToUint8(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadUint16:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadUint16 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeUint16(e.ptrToUint16(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUint16:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeUint16(e.ptrToUint16(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadUint32:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadUint32 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeUint32(e.ptrToUint32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUint32:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeUint32(e.ptrToUint32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadUint64:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadUint64 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeUint64(e.ptrToUint64(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadUint64:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeUint64(e.ptrToUint64(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadFloat32:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadFloat32 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeFloat32(e.ptrToFloat32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat32:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeFloat32(e.ptrToFloat32(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadFloat64:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadFloat64 {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				v := e.ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeFloat64(v)
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadFloat64:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeKey(code)
				e.encodeFloat64(v)
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadString:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadString {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(e.ptrToString(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadString:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(e.ptrToString(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadBool:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadBool {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeBool(e.ptrToBool(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadBool:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeBool(e.ptrToBool(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadBytes:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadBytes {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'{', '}', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeByteSlice(e.ptrToBytes(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadBytes:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeByteSlice(e.ptrToBytes(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadArray:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadArray:
			ptr := load(ctxptr, code.idx) + code.offset
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadArray {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'[', ']', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				if !code.anonymousKey {
					e.encodeKey(code)
				}
				code = code.next
				store(ctxptr, code.idx, ptr)
			}
		case opStructFieldPtrAnonymousHeadArray:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadArray:
			ptr := load(ctxptr, code.idx) + code.offset
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				store(ctxptr, code.idx, ptr)
				code = code.next
			}
		case opStructFieldPtrHeadSlice:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadSlice:
			ptr := load(ctxptr, code.idx)
			p := ptr + code.offset
			if p == 0 {
				if code.op == opStructFieldPtrHeadSlice {
					e.encodeNull()
					e.encodeByte(',')
				} else {
					e.encodeBytes([]byte{'[', ']', ','})
				}
				code = code.end.next
			} else {
				e.encodeByte('{')
				if !code.anonymousKey {
					e.encodeKey(code)
				}
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrAnonymousHeadSlice:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadSlice:
			ptr := load(ctxptr, code.idx)
			p := ptr + code.offset
			if p == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				store(ctxptr, code.idx, p)
				code = code.next
			}
		case opStructFieldPtrHeadMarshalJSON:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				ptr += code.offset
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
				}))
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					e.encodeNull()
					code = code.end
					break
				}
				b, err := rv.Interface().(Marshaler).MarshalJSON()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(b) == 0 {
					return errUnexpectedEndOfJSON(
						fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
						0,
					)
				}
				var buf bytes.Buffer
				if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
					return err
				}
				e.encodeBytes(buf.Bytes())
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadMarshalJSON:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				ptr += code.offset
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
				}))
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					e.encodeNull()
					code = code.end.next
					break
				}
				b, err := rv.Interface().(Marshaler).MarshalJSON()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(b) == 0 {
					return errUnexpectedEndOfJSON(
						fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
						0,
					)
				}
				var buf bytes.Buffer
				if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
					return err
				}
				e.encodeBytes(buf.Bytes())
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadMarshalText:
			p := load(ctxptr, code.idx)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
				break
			}
			store(ctxptr, code.idx, e.ptrToPtr(p))
			fallthrough
		case opStructFieldHeadMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				ptr += code.offset
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
				}))
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					e.encodeNull()
					e.encodeByte(',')
					code = code.end
					break
				}
				bytes, err := rv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadMarshalText:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldAnonymousHeadMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				ptr += code.offset
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: *(*unsafe.Pointer)(unsafe.Pointer(&ptr)),
				}))
				rv := reflect.ValueOf(v)
				if rv.Type().Kind() == reflect.Interface && rv.IsNil() {
					e.encodeNull()
					e.encodeByte(',')
					code = code.end.next
					break
				}
				bytes, err := rv.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else if code.next == code.end {
				// not exists fields
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '}', ',', '\n'})
				code = code.end.next
				store(ctxptr, code.idx, ptr)
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				code = code.next
				store(ctxptr, code.idx, ptr)
			}
		case opStructFieldPtrHeadIntIndent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadIntIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				if code.op == opStructFieldPtrHeadIntIndent {
					e.encodeIndent(code.indent)
					e.encodeNull()
					e.encodeBytes([]byte{',', '\n'})
				} else {
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
				}
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt(e.ptrToInt(ptr + code.offset))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadInt8Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadInt8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt8(e.ptrToInt8(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadInt16Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadInt16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt16(e.ptrToInt16(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadInt32Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadInt32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt32(e.ptrToInt32(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadInt64Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadInt64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt64(e.ptrToInt64(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadUintIndent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUintIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint(e.ptrToUint(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadUint8Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUint8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint8(e.ptrToUint8(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadUint16Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUint16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint16(e.ptrToUint16(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadUint32Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUint32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint32(e.ptrToUint32(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadUint64Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadUint64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint64(e.ptrToUint64(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadFloat32Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeFloat32(e.ptrToFloat32(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadFloat64Indent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadFloat64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				v := e.ptrToFloat64(ptr)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeFloat64(v)
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringIndent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadStringIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(e.ptrToString(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadBoolIndent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadBoolIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeBool(e.ptrToBool(ptr))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadBytesIndent:
			store(ctxptr, code.idx, e.ptrToPtr(load(ctxptr, code.idx)))
			fallthrough
		case opStructFieldHeadBytesIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				s := base64.StdEncoding.EncodeToString(e.ptrToBytes(ptr))
				e.encodeByte('"')
				e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
				e.encodeByte('"')
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				p := ptr + code.offset
				if p == 0 || *(*uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&p))) == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmpty:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				p := ptr + code.offset
				if p == 0 || *(*uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&p))) == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToInt(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToInt(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToInt8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt8(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToInt8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt8(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToInt16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt16(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToInt16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt16(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToInt32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToInt32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToInt64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyInt64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToInt64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeInt64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToUint(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToUint(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToUint8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint8(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToUint8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint8(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToUint16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint16(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToUint16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint16(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToUint32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToUint32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToUint64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyUint64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToUint64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeUint64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeFloat32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeFloat32(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return &UnsupportedValueError{
							Value: reflect.ValueOf(v),
							Str:   strconv.FormatFloat(v, 'g', -1, 64),
						}
					}
					e.encodeKey(code)
					e.encodeFloat64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return &UnsupportedValueError{
							Value: reflect.ValueOf(v),
							Str:   strconv.FormatFloat(v, 'g', -1, 64),
						}
					}
					e.encodeKey(code)
					e.encodeFloat64(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeString(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeString(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToBool(ptr + code.offset)
				if !v {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeBool(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToBool(ptr + code.offset)
				if !v {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeBool(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToBytes(ptr + code.offset)
				if len(v) == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeByteSlice(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToBytes(ptr + code.offset)
				if len(v) == 0 {
					code = code.nextField
				} else {
					e.encodeKey(code)
					e.encodeByteSlice(v)
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					b, err := v.(Marshaler).MarshalJSON()
					if err != nil {
						return &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					if len(b) == 0 {
						if isPtr {
							return errUnexpectedEndOfJSON(
								fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
								0,
							)
						}
						code = code.nextField
					} else {
						var buf bytes.Buffer
						if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
							return err
						}
						e.encodeKey(code)
						e.encodeBytes(buf.Bytes())
						e.encodeByte(',')
						code = code.next
					}
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					b, err := v.(Marshaler).MarshalJSON()
					if err != nil {
						return &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					if len(b) == 0 {
						if isPtr {
							return errUnexpectedEndOfJSON(
								fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
								0,
							)
						}
						code = code.nextField
					} else {
						var buf bytes.Buffer
						if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
							return err
						}
						e.encodeKey(code)
						e.encodeBytes(buf.Bytes())
						e.encodeByte(',')
						code = code.next
					}
				}
			}
		case opStructFieldPtrHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bytes, err := v.(encoding.TextMarshaler).MarshalText()
					if err != nil {
						return &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					e.encodeKey(code)
					e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadOmitEmptyMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				if p == nil || (!isPtr && *(*unsafe.Pointer)(p) == nil) {
					code = code.nextField
				} else {
					v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
					bytes, err := v.(encoding.TextMarshaler).MarshalText()
					if err != nil {
						return &MarshalerError{
							Type: rtype2type(code.typ),
							Err:  err,
						}
					}
					e.encodeKey(code)
					e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				p := ptr + code.offset
				if p == 0 || *(*uintptr)(*(*unsafe.Pointer)(unsafe.Pointer(&p))) == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldPtrHeadOmitEmptyIntIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyIntIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToInt(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeInt(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToInt8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeInt8(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToInt16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeInt16(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToInt32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeInt32(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyInt64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyInt64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToInt64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeInt64(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUintIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUintIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToUint(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeUint(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToUint8(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeUint8(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToUint16(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeUint16(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToUint32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeUint32(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyUint64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyUint64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToUint64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeUint64(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyFloat32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToFloat32(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeFloat32(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyFloat64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyFloat64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToFloat64(ptr + code.offset)
				if v == 0 {
					code = code.nextField
				} else {
					if math.IsInf(v, 0) || math.IsNaN(v) {
						return &UnsupportedValueError{
							Value: reflect.ValueOf(v),
							Str:   strconv.FormatFloat(v, 'g', -1, 64),
						}
					}
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeFloat64(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyStringIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyStringIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToString(ptr + code.offset)
				if v == "" {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeString(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyBoolIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBoolIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToBool(ptr + code.offset)
				if !v {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					e.encodeBool(v)
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadOmitEmptyBytesIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadOmitEmptyBytesIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeIndent(code.indent)
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToBytes(ptr + code.offset)
				if len(v) == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent + 1)
					e.encodeKey(code)
					e.encodeByte(' ')
					s := base64.StdEncoding.EncodeToString(v)
					e.encodeByte('"')
					e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
					e.encodeByte('"')
					e.encodeBytes([]byte{',', '\n'})
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				p := ptr + code.offset
				e.encodeKey(code)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrAnonymousHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTag:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				code = code.next
				store(ctxptr, code.idx, ptr+code.offset)
			}
		case opStructFieldPtrHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt8(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagInt8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt8(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt16(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagInt16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt16(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagInt32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt64(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagInt64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagInt64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToInt64(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint8(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUint8:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint8:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint8(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint16(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUint16:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint16:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint16(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUint32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint64(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagUint64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagUint64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToUint64(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToFloat32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat32:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToFloat32(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				v := e.ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(v))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagFloat64:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				v := e.ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(v))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				s := e.ptrToString(ptr + code.offset)
				if e.enabledHTMLEscape {
					enc.encodeEscapedString(s)
				} else {
					enc.encodeNoEscapedString(s)
				}
				e.encodeString(string(enc.buf))
				e.encodeByte(',')
				enc.release()
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagString:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(strconv.Quote(e.ptrToString(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToBool(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagBool:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeString(fmt.Sprint(e.ptrToBool(ptr + code.offset)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				e.encodeKey(code)
				e.encodeByteSlice(e.ptrToBytes(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagBytes:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				e.encodeKey(code)
				e.encodeByteSlice(e.ptrToBytes(ptr + code.offset))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				b, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(b) == 0 {
					if isPtr {
						return errUnexpectedEndOfJSON(
							fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
							0,
						)
					}
					e.encodeKey(code)
					e.encodeBytes([]byte{'"', '"'})
					e.encodeByte(',')
					code = code.nextField
				} else {
					var buf bytes.Buffer
					if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
						return err
					}
					e.encodeString(buf.String())
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrAnonymousHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagMarshalJSON:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				isPtr := code.typ.Kind() == reflect.Ptr
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				b, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				if len(b) == 0 {
					if isPtr {
						return errUnexpectedEndOfJSON(
							fmt.Sprintf("error calling MarshalJSON for type %s", code.typ),
							0,
						)
					}
					e.encodeKey(code)
					e.encodeBytes([]byte{'"', '"'})
					e.encodeByte(',')
					code = code.nextField
				} else {
					var buf bytes.Buffer
					if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
						return err
					}
					e.encodeKey(code)
					e.encodeString(buf.String())
					e.encodeByte(',')
					code = code.next
				}
			}
		case opStructFieldPtrHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.end.next
			} else {
				e.encodeByte('{')
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeKey(code)
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrAnonymousHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldAnonymousHeadStringTagMarshalText:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				code = code.end.next
			} else {
				ptr += code.offset
				p := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{typ: code.typ, ptr: p}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeKey(code)
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
				code = code.next
			}
		case opStructFieldPtrHeadStringTagIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				p := ptr + code.offset
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldPtrHeadStringTagIntIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagIntIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToInt(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToInt8(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToInt16(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToInt32(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagInt64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagInt64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToInt64(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUintIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUintIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToUint(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint8Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToUint8(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint16Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToUint16(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToUint32(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagUint64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagUint64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToUint64(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat32Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToFloat32(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagFloat64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagFloat64Indent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				v := e.ptrToFloat64(ptr + code.offset)
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(v))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagStringIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagStringIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				s := e.ptrToString(ptr + code.offset)
				if e.enabledHTMLEscape {
					enc.encodeEscapedString(s)
				} else {
					enc.encodeNoEscapedString(s)
				}
				e.encodeString(string(enc.buf))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagBoolIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBoolIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(fmt.Sprint(e.ptrToBool(ptr + code.offset)))
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructFieldPtrHeadStringTagBytesIndent:
			ptr := load(ctxptr, code.idx)
			if ptr != 0 {
				store(ctxptr, code.idx, e.ptrToPtr(ptr))
			}
			fallthrough
		case opStructFieldHeadStringTagBytesIndent:
			ptr := load(ctxptr, code.idx)
			if ptr == 0 {
				e.encodeIndent(code.indent)
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.end.next
			} else {
				e.encodeBytes([]byte{'{', '\n'})
				e.encodeIndent(code.indent + 1)
				e.encodeKey(code)
				e.encodeByte(' ')
				s := base64.StdEncoding.EncodeToString(
					e.ptrToBytes(ptr + code.offset),
				)
				e.encodeByte('"')
				e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
				e.encodeByte('"')
				e.encodeBytes([]byte{',', '\n'})
				code = code.next
			}
		case opStructField:
			if !code.anonymousKey {
				e.encodeKey(code)
			}
			ptr := load(ctxptr, code.headIdx) + code.offset
			code = code.next
			store(ctxptr, code.idx, ptr)
		case opStructFieldPtrInt:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeInt(e.ptrToInt(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldInt:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeInt(e.ptrToInt(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrInt8:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeInt8(e.ptrToInt8(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldInt8:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeInt8(e.ptrToInt8(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrInt16:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeInt16(e.ptrToInt16(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldInt16:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeInt16(e.ptrToInt16(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrInt32:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeInt32(e.ptrToInt32(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldInt32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeInt32(e.ptrToInt32(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrInt64:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeInt64(e.ptrToInt64(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldInt64:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeInt64(e.ptrToInt64(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrUint:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeUint(e.ptrToUint(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldUint:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeUint(e.ptrToUint(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrUint8:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeUint8(e.ptrToUint8(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldUint8:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeUint8(e.ptrToUint8(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrUint16:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeUint16(e.ptrToUint16(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldUint16:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeUint16(e.ptrToUint16(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrUint32:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeUint32(e.ptrToUint32(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldUint32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeUint32(e.ptrToUint32(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrUint64:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeUint64(e.ptrToUint64(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldUint64:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeUint64(e.ptrToUint64(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrFloat32:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeFloat32(e.ptrToFloat32(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldFloat32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeFloat32(e.ptrToFloat32(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrFloat64:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
				e.encodeByte(',')
				code = code.next
				break
			}
			v := e.ptrToFloat64(p)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeFloat64(v)
			e.encodeByte(',')
			code = code.next
		case opStructFieldFloat64:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			v := e.ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeFloat64(v)
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrString:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeString(e.ptrToString(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldString:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(e.ptrToString(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldPtrBool:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := e.ptrToPtr(ptr + code.offset)
			if p == 0 {
				e.encodeNull()
			} else {
				e.encodeBool(e.ptrToBool(p))
			}
			e.encodeByte(',')
			code = code.next
		case opStructFieldBool:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeBool(e.ptrToBool(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldBytes:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeByteSlice(e.ptrToBytes(ptr + code.offset))
			e.encodeByte(',')
			code = code.next
		case opStructFieldMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			var buf bytes.Buffer
			if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
				return err
			}
			e.encodeBytes(buf.Bytes())
			e.encodeByte(',')
			code = code.next
		case opStructFieldMarshalText:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldArray:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldSlice:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldMap:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldMapLoad:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldStruct:
			e.encodeKey(code)
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldIntIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeInt(e.ptrToInt(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldInt8Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeInt8(e.ptrToInt8(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldInt16Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeInt16(e.ptrToInt16(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldInt32Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeInt32(e.ptrToInt32(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldInt64Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeInt64(e.ptrToInt64(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldUintIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeUint(e.ptrToUint(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldUint8Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeUint8(e.ptrToUint8(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldUint16Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeUint16(e.ptrToUint16(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldUint32Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeUint32(e.ptrToUint32(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldUint64Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeUint64(e.ptrToUint64(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldFloat32Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeFloat32(e.ptrToFloat32(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldFloat64Indent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeFloat64(v)
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeString(e.ptrToString(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldBoolIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			e.encodeBool(e.ptrToBool(ptr + code.offset))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldBytesIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			s := base64.StdEncoding.EncodeToString(e.ptrToBytes(ptr + code.offset))
			e.encodeByte('"')
			e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
			e.encodeByte('"')
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldMarshalJSONIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			var buf bytes.Buffer
			if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
				return err
			}
			e.encodeBytes(buf.Bytes())
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldArrayIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.nextField
			} else {
				code = code.next
			}
		case opStructFieldSliceIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				e.encodeNull()
				e.encodeBytes([]byte{',', '\n'})
				code = code.nextField
			} else {
				code = code.next
			}
		case opStructFieldMapIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				e.encodeNull()
				code = code.nextField
			} else {
				p = uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				mlen := maplen(*(*unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					mapCode := code.next
					code = mapCode.end.next
				} else {
					code = code.next
				}
			}
		case opStructFieldMapLoadIndent:
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				e.encodeNull()
				code = code.nextField
			} else {
				p = uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				mlen := maplen(*(*unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					code = code.nextField
				} else {
					code = code.next
				}
			}
		case opStructFieldStructIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			if p == 0 {
				e.encodeBytes([]byte{'{', '}', ',', '\n'})
				code = code.nextField
			} else {
				headCode := code.next
				if headCode.next == headCode.end {
					// not exists fields
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					code = code.nextField
				} else {
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldOmitEmpty:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 || **(**uintptr)(unsafe.Pointer(&p)) == 0 {
				code = code.nextField
			} else {
				e.encodeKey(code)
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldOmitEmptyInt:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeInt(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyInt8:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt8(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeInt8(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyInt16:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt16(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeInt16(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyInt32:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt32(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeInt32(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyInt64:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt64(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeInt64(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyUint:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeUint(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyUint8:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint8(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeUint8(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyUint16:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint16(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeUint16(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyUint32:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint32(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeUint32(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyUint64:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint64(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeUint64(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyFloat32:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat32(ptr + code.offset)
			if v != 0 {
				e.encodeKey(code)
				e.encodeFloat32(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat64(ptr + code.offset)
			if v != 0 {
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeKey(code)
				e.encodeFloat64(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyString:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToString(ptr + code.offset)
			if v != "" {
				e.encodeKey(code)
				e.encodeString(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyBool:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToBool(ptr + code.offset)
			if v {
				e.encodeKey(code)
				e.encodeBool(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyBytes:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToBytes(ptr + code.offset)
			if len(v) > 0 {
				e.encodeKey(code)
				e.encodeByteSlice(v)
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			if v != nil {
				b, err := v.(Marshaler).MarshalJSON()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				var buf bytes.Buffer
				if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
					return err
				}
				e.encodeBytes(buf.Bytes())
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			if v != nil {
				v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
					typ: code.typ,
					ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
				}))
				bytes, err := v.(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return &MarshalerError{
						Type: rtype2type(code.typ),
						Err:  err,
					}
				}
				e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
				e.encodeByte(',')
			}
			code = code.next
		case opStructFieldOmitEmptyArray:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				code = code.nextField
			} else {
				code = code.next
			}
		case opStructFieldOmitEmptySlice:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				code = code.nextField
			} else {
				code = code.next
			}
		case opStructFieldOmitEmptyMap:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					code = code.next
				}
			}
		case opStructFieldOmitEmptyMapLoad:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					code = code.next
				}
			}
		case opStructFieldOmitEmptyIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 || **(**uintptr)(unsafe.Pointer(&p)) == 0 {
				code = code.nextField
			} else {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				code = code.next
				store(ctxptr, code.idx, p)
			}
		case opStructFieldOmitEmptyIntIndent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyInt8Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt8(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt8(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyInt16Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt16(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt16(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyInt32Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt32(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt32(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyInt64Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToInt64(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeInt64(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyUintIndent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyUint8Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint8(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint8(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyUint16Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint16(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint16(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyUint32Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint32(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint32(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyUint64Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToUint64(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeUint64(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyFloat32Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat32(ptr + code.offset)
			if v != 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeFloat32(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyFloat64Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat64(ptr + code.offset)
			if v != 0 {
				if math.IsInf(v, 0) || math.IsNaN(v) {
					return &UnsupportedValueError{
						Value: reflect.ValueOf(v),
						Str:   strconv.FormatFloat(v, 'g', -1, 64),
					}
				}
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeFloat64(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyStringIndent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToString(ptr + code.offset)
			if v != "" {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeString(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyBoolIndent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToBool(ptr + code.offset)
			if v {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				e.encodeBool(v)
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyBytesIndent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToBytes(ptr + code.offset)
			if len(v) > 0 {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				s := base64.StdEncoding.EncodeToString(v)
				e.encodeByte('"')
				e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
				e.encodeByte('"')
				e.encodeBytes([]byte{',', '\n'})
			}
			code = code.next
		case opStructFieldOmitEmptyArrayIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				code = code.nextField
			} else {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				code = code.next
			}
		case opStructFieldOmitEmptySliceIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			header := *(**sliceHeader)(unsafe.Pointer(&p))
			if p == 0 || uintptr(header.data) == 0 {
				code = code.nextField
			} else {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				code = code.next
			}
		case opStructFieldOmitEmptyMapIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent)
					e.encodeKey(code)
					e.encodeByte(' ')
					code = code.next
				}
			}
		case opStructFieldOmitEmptyMapLoadIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				mlen := maplen(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
				if mlen == 0 {
					code = code.nextField
				} else {
					e.encodeIndent(code.indent)
					e.encodeKey(code)
					e.encodeByte(' ')
					code = code.next
				}
			}
		case opStructFieldOmitEmptyStructIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			if p == 0 {
				code = code.nextField
			} else {
				e.encodeIndent(code.indent)
				e.encodeKey(code)
				e.encodeByte(' ')
				headCode := code.next
				if headCode.next == headCode.end {
					// not exists fields
					e.encodeBytes([]byte{'{', '}', ',', '\n'})
					code = code.nextField
				} else {
					code = code.next
					store(ctxptr, code.idx, p)
				}
			}
		case opStructFieldStringTag:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			e.encodeKey(code)
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldStringTagInt:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToInt(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagInt8:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToInt8(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagInt16:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToInt16(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagInt32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToInt32(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagInt64:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToInt64(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagUint:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToUint(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagUint8:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToUint8(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagUint16:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToUint16(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagUint32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToUint32(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagUint64:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToUint64(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagFloat32:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToFloat32(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagFloat64:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(v))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagString:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			var b bytes.Buffer
			enc := NewEncoder(&b)
			enc.encodeString(e.ptrToString(ptr + code.offset))
			e.encodeString(string(enc.buf))
			code = code.next
		case opStructFieldStringTagBool:
			ptr := load(ctxptr, code.headIdx)
			e.encodeKey(code)
			e.encodeString(fmt.Sprint(e.ptrToBool(ptr + code.offset)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagBytes:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToBytes(ptr + code.offset)
			e.encodeKey(code)
			e.encodeByteSlice(v)
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagMarshalJSON:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			var buf bytes.Buffer
			if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
				return err
			}
			e.encodeString(buf.String())
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagMarshalText:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
			e.encodeByte(',')
			code = code.next
		case opStructFieldStringTagIndent:
			ptr := load(ctxptr, code.headIdx)
			p := ptr + code.offset
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			code = code.next
			store(ctxptr, code.idx, p)
		case opStructFieldStringTagIntIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToInt(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagInt8Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToInt8(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagInt16Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToInt16(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagInt32Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToInt32(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagInt64Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToInt64(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagUintIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToUint(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagUint8Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToUint8(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagUint16Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToUint16(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagUint32Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToUint32(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagUint64Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToUint64(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagFloat32Indent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToFloat32(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagFloat64Indent:
			ptr := load(ctxptr, code.headIdx)
			v := e.ptrToFloat64(ptr + code.offset)
			if math.IsInf(v, 0) || math.IsNaN(v) {
				return &UnsupportedValueError{
					Value: reflect.ValueOf(v),
					Str:   strconv.FormatFloat(v, 'g', -1, 64),
				}
			}
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(v))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagStringIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			var b bytes.Buffer
			enc := NewEncoder(&b)
			enc.encodeString(e.ptrToString(ptr + code.offset))
			e.encodeString(string(enc.buf))
			e.encodeBytes([]byte{',', '\n'})
			enc.release()
			code = code.next
		case opStructFieldStringTagBoolIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			e.encodeString(fmt.Sprint(e.ptrToBool(ptr + code.offset)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagBytesIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			s := base64.StdEncoding.EncodeToString(
				e.ptrToBytes(ptr + code.offset),
			)
			e.encodeByte('"')
			e.encodeBytes(*(*[]byte)(unsafe.Pointer(&s)))
			e.encodeByte('"')
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagMarshalJSONIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			b, err := v.(Marshaler).MarshalJSON()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			var buf bytes.Buffer
			if err := compact(&buf, b, e.enabledHTMLEscape); err != nil {
				return err
			}
			e.encodeString(buf.String())
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructFieldStringTagMarshalTextIndent:
			ptr := load(ctxptr, code.headIdx)
			e.encodeIndent(code.indent)
			e.encodeKey(code)
			e.encodeByte(' ')
			p := ptr + code.offset
			v := *(*interface{})(unsafe.Pointer(&interfaceHeader{
				typ: code.typ,
				ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
			}))
			bytes, err := v.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return &MarshalerError{
					Type: rtype2type(code.typ),
					Err:  err,
				}
			}
			e.encodeString(*(*string)(unsafe.Pointer(&bytes)))
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opStructEnd:
			last := len(e.buf) - 1
			if e.buf[last] == ',' {
				e.buf[last] = '}'
			} else {
				e.encodeByte('}')
			}
			e.encodeByte(',')
			code = code.next
		case opStructAnonymousEnd:
			code = code.next
		case opStructEndIndent:
			last := len(e.buf) - 1
			if e.buf[last] == '\n' {
				// to remove ',' and '\n' characters
				e.buf = e.buf[:len(e.buf)-2]
			}
			e.encodeByte('\n')
			e.encodeIndent(code.indent)
			e.encodeByte('}')
			e.encodeBytes([]byte{',', '\n'})
			code = code.next
		case opEnd:
			goto END
		}
	}
END:
	return nil
}

func (e *Encoder) ptrToPtr(p uintptr) uintptr     { return **(**uintptr)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToInt(p uintptr) int         { return **(**int)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToInt8(p uintptr) int8       { return **(**int8)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToInt16(p uintptr) int16     { return **(**int16)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToInt32(p uintptr) int32     { return **(**int32)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToInt64(p uintptr) int64     { return **(**int64)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToUint(p uintptr) uint       { return **(**uint)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToUint8(p uintptr) uint8     { return **(**uint8)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToUint16(p uintptr) uint16   { return **(**uint16)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToUint32(p uintptr) uint32   { return **(**uint32)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToUint64(p uintptr) uint64   { return **(**uint64)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToFloat32(p uintptr) float32 { return **(**float32)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToFloat64(p uintptr) float64 { return **(**float64)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToBool(p uintptr) bool       { return **(**bool)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToByte(p uintptr) byte       { return **(**byte)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToBytes(p uintptr) []byte    { return **(**[]byte)(unsafe.Pointer(&p)) }
func (e *Encoder) ptrToString(p uintptr) string   { return **(**string)(unsafe.Pointer(&p)) }
