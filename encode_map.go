package json

import (
	"sync"
	"unsafe"
)

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
)

// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

// mapextra holds fields that are not present on all maps.
type mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

// hiter is the hash iteration structure.
// If you modify hiter, also change cmd/compile/internal/gc/reflect.go to indicate
// the layout of this structure.
type hiter struct {
	key         unsafe.Pointer
	elem        unsafe.Pointer // _    [5]unsafe.Pointer // key, elem, t, h, buckets
	_           unsafe.Pointer // t
	h           *hmap
	buckets     unsafe.Pointer // bucket ptr at hash_iter initialization time
	bptr        *bmap
	overflow    *[]*bmap
	oldoverflow *[]*bmap   // _       [2]unsafe.Pointer // overflow, oldoverflow
	_           uintptr    // startBucket
	_           uint8      // offset
	_           bool       // wrapped
	_           [2]uint8   // B, i
	_           [2]uintptr // bucket, checkBucket
}

var (
	hiterPool sync.Pool // *hiter
)

var zeroHiter = &hiter{}

func (h *hiter) put() {
	*h = *zeroHiter
	hiterPool.Put(h)
}

//go:linkname newmapiterinit runtime.mapiterinit
//go:noescape
func newmapiterinit(t *rtype, h *hmap, it *hiter)

func newHiter(t *rtype, h *hmap) *hiter {
	v := hiterPool.Get()
	if v == nil {
		return mapiterinit(t, h)
	}
	it := v.(*hiter)
	newmapiterinit(t, h, it)
	return it
}
