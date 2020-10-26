package json

import (
	"bytes"
	"testing"
	"unsafe"
)

// TestHeader ensures that the headers can successfully mutate
// variables of the corresponding built-in types.
//
// This test is expected to fail under -race (which implicitly enables
// -d=checkptr) if the runtime views the header types as incompatible with the
// underlying built-in types.
func TestHeader(t *testing.T) {
	t.Run("sliceHeader", func(t *testing.T) {
		s := []byte("Hello, checkptr!")[:5]

		var alias []byte
		hdr := (*sliceHeader)(unsafe.Pointer(&alias))
		hdr.data = unsafe.Pointer(&s[0])
		hdr.cap = cap(s)
		hdr.len = len(s)

		if !bytes.Equal(alias, s) {
			t.Fatalf("alias of %T(%q) constructed via sliceHeader = %T(%q)", s, s, alias, alias)
		}
		if cap(alias) != cap(s) {
			t.Fatalf("alias of %T with cap %d has cap %d", s, cap(s), cap(alias))
		}
	})

	t.Run("stringHeader", func(t *testing.T) {
		s := "Hello, checkptr!"

		var alias string
		hdr := (*stringHeader)(unsafe.Pointer(&alias))
		hdr.data = (*stringHeader)(unsafe.Pointer(&s)).data
		hdr.len = len(s)

		if alias != s {
			t.Fatalf("alias of %q constructed via stringHeader = %q", s, alias)
		}
	})
}
