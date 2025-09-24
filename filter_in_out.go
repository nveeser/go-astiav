package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html
type FilterInOut struct {
	c *C.AVFilterInOut
}

func newFilterInOutFromC(c *C.AVFilterInOut) *FilterInOut {
	if c == nil {
		return nil
	}
	return &FilterInOut{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga6e1c2931e15eb4283c59c6ccc8b83919
func AllocFilterInOut() *FilterInOut {
	return newFilterInOutFromC(C.avfilter_inout_alloc())
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga294500a9856260eb1552354fd9d9a6c4
func (i *FilterInOut) Free() {
	if i.c != nil {
		C.avfilter_inout_free(&i.c)
	}
}
func (i *FilterInOut) String() string {
	return fmt.Sprintf("%s(%s)@%d", i.Name(), i.FilterContext().Name(), i.PadIdx())
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a88afecac258f51aab7e9a9db9e7a4d58
func (i *FilterInOut) SetName(n string) {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	i.c.name = C.av_strdup(cn)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a88afecac258f51aab7e9a9db9e7a4d58
func (i *FilterInOut) Name() string {
	return C.GoString(i.c.name)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a3227857d0b955b639f4950d13e4e6f40
func (i *FilterInOut) SetFilterContext(c *FilterContext) {
	i.c.filter_ctx = (*C.AVFilterContext)(c.c)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a3227857d0b955b639f4950d13e4e6f40
func (i *FilterInOut) FilterContext() *FilterContext {
	return &FilterContext{c: i.c.filter_ctx}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a386ff90d40aa22f5612dd5eca734ed48
func (i *FilterInOut) SetPadIdx(idx int) {
	i.c.pad_idx = C.int(idx)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#a386ff90d40aa22f5612dd5eca734ed48
func (i *FilterInOut) PadIdx() int {
	return int(i.c.pad_idx)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterInOut.html#af8c8cf9ffb650974d19e791f5bb7cf33
func (i *FilterInOut) SetNext(n *FilterInOut) {
	var nc *C.AVFilterInOut
	if n != nil {
		nc = n.c
	}
	i.c.next = nc
}

func (i *FilterInOut) Next() *FilterInOut {
	if i.c.next == nil {
		return nil
	}
	return newFilterInOutFromC(i.c.next)
}
