package astiav

//#include <libavfilter/avfilter.h>
import "C"

// Struct attributes are internal but there are C functions to get some of them
type FilterPad struct {
	name      string
	mediaType MediaType
}

func newFilterPad(pads *C.AVFilterPad, idx int) *FilterPad {
	name := C.GoString(C.avfilter_pad_get_name(pads, C.int(idx)))
	mediaType := MediaType(C.avfilter_pad_get_type(pads, C.int(idx)))
	return &FilterPad{name: name, mediaType: mediaType}
}

func (fp *FilterPad) Name() string {
	return fp.name
}

func (fp *FilterPad) MediaType() MediaType {
	return fp.mediaType
}
