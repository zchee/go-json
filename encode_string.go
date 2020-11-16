package json

import (
	"unicode/utf8"
)

// htmlSafeSet holds the value true if the ASCII character with the given
// array position can be safely represented inside a JSON string, embedded
// inside of HTML <script> tags, without any additional escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), the backslash character ("\"), HTML opening and closing
// tags ("<" and ">"), and the ampersand ("&").
var htmlSafeSet = [256]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
	0x80:     false,
	0x81:     false,
	0x82:     false,
	0x83:     false,
	0x84:     false,
	0x85:     false,
	0x86:     false,
	0x87:     false,
	0x88:     false,
	0x89:     false,
	0x8a:     false,
	0x8b:     false,
	0x8c:     false,
	0x8d:     false,
	0x8e:     false,
	0x8f:     false,
	0x90:     false,
	0x91:     false,
	0x92:     false,
	0x93:     false,
	0x94:     false,
	0x95:     false,
	0x96:     false,
	0x97:     false,
	0x98:     false,
	0x99:     false,
	0x9a:     false,
	0x9b:     false,
	0x9c:     false,
	0x9d:     false,
	0x9e:     false,
	0x9f:     false,
	0xa0:     false,
	0xa1:     false,
	0xa2:     false,
	0xa3:     false,
	0xa4:     false,
	0xa5:     false,
	0xa6:     false,
	0xa7:     false,
	0xa8:     false,
	0xa9:     false,
	0xaa:     false,
	0xab:     false,
	0xac:     false,
	0xad:     false,
	0xae:     false,
	0xaf:     false,
	0xb0:     false,
	0xb1:     false,
	0xb2:     false,
	0xb3:     false,
	0xb4:     false,
	0xb5:     false,
	0xb6:     false,
	0xb7:     false,
	0xb8:     false,
	0xb9:     false,
	0xba:     false,
	0xbb:     false,
	0xbc:     false,
	0xbd:     false,
	0xbe:     false,
	0xbf:     false,
	0xc0:     false,
	0xc1:     false,
	0xc2:     false,
	0xc3:     false,
	0xc4:     false,
	0xc5:     false,
	0xc6:     false,
	0xc7:     false,
	0xc8:     false,
	0xc9:     false,
	0xca:     false,
	0xcb:     false,
	0xcc:     false,
	0xcd:     false,
	0xce:     false,
	0xcf:     false,
	0xd0:     false,
	0xd1:     false,
	0xd2:     false,
	0xd3:     false,
	0xd4:     false,
	0xd5:     false,
	0xd6:     false,
	0xd7:     false,
	0xd8:     false,
	0xd9:     false,
	0xda:     false,
	0xdb:     false,
	0xdc:     false,
	0xdd:     false,
	0xde:     false,
	0xdf:     false,
	0xe0:     false,
	0xe1:     false,
	0xe2:     false,
	0xe3:     false,
	0xe4:     false,
	0xe5:     false,
	0xe6:     false,
	0xe7:     false,
	0xe8:     false,
	0xe9:     false,
	0xea:     false,
	0xeb:     false,
	0xec:     false,
	0xed:     false,
	0xee:     false,
	0xef:     false,
	0xf0:     false,
	0xf1:     false,
	0xf2:     false,
	0xf3:     false,
	0xf4:     false,
	0xf5:     false,
	0xf6:     false,
	0xf7:     false,
	0xf8:     false,
	0xf9:     false,
	0xfa:     false,
	0xfb:     false,
	0xfc:     false,
	0xfd:     false,
	0xfe:     false,
	0xff:     false,
}

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var hex = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func (e *Encoder) encodeEscapedString(s string) {
	valLen := len(s)
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < valLen; i++ {
		if !htmlSafeSet[s[i]] {
			break
		}
	}
	e.buf = append(e.buf, '"')
	if i == valLen {
		e.buf = append(e.buf, s...)
		e.buf = append(e.buf, '"')
		return
	}
	e.buf = append(e.buf, s[:i]...)
	e.writeStringSlowPathWithHTMLEscaped(i, s, valLen)
}

func (e *Encoder) writeStringSlowPathWithHTMLEscaped(i int, s string, valLen int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				e.buf = append(e.buf, s[start:i]...)
			}
			switch b {
			case '\\', '"':
				e.buf = append(e.buf, '\\', b)
			case '\n':
				e.buf = append(e.buf, '\\', 'n')
			case '\r':
				e.buf = append(e.buf, '\\', 'r')
			case '\t':
				e.buf = append(e.buf, '\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.buf = append(e.buf, `\u00`...)
				e.buf = append(e.buf, hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.buf = append(e.buf, s[start:i]...)
			}
			e.buf = append(e.buf, `\ufffd`...)
			i++
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.buf = append(e.buf, s[start:i]...)
			}
			e.buf = append(e.buf, `\u202`...)
			e.buf = append(e.buf, hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.buf = append(e.buf, s[start:]...)
	}
	e.buf = append(e.buf, '"')
}

func (e *Encoder) encodeNoEscapedString(s string) {
	valLen := len(s)

	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < valLen; i++ {
		c := s[i]
		if c <= 31 || c == '"' || c == '\\' {
			break
		}
	}
	e.buf = append(e.buf, '"')
	if i == valLen {
		e.buf = append(e.buf, s...)
		e.buf = append(e.buf, '"')
		return
	}
	e.buf = append(e.buf, s[:i]...)
	e.writeStringSlowPath(i, s, valLen)
}

func (e *Encoder) writeStringSlowPath(i int, s string, valLen int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				e.buf = append(e.buf, s[start:i]...)
			}
			switch b {
			case '\\', '"':
				e.buf = append(e.buf, '\\', b)
			case '\n':
				e.buf = append(e.buf, '\\', 'n')
			case '\r':
				e.buf = append(e.buf, '\\', 'r')
			case '\t':
				e.buf = append(e.buf, '\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.buf = append(e.buf, []byte(`\u00`)...)
				e.buf = append(e.buf, hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		i++
		continue
	}
	if start < len(s) {
		e.buf = append(e.buf, s[start:]...)
	}
	e.buf = append(e.buf, '"')
}
