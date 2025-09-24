package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraphSegment.html
type FilterGraphSegment struct {
	g *FilterGraph
	c *C.AVFilterGraphSegment
}

func newFilterGraphSegmentFromC(g *FilterGraph, c *C.AVFilterGraphSegment) *FilterGraphSegment {
	if c == nil {
		return nil
	}
	return &FilterGraphSegment{g: g, c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga51283edd8f3685e1f33239f360e14ae8
func (fgs *FilterGraphSegment) Free() {
	if fgs.c != nil {
		C.avfilter_graph_segment_free(&fgs.c)
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraphSegment.html#ad5a2779af221d1520490fe2719f9e39a
func (fgs *FilterGraphSegment) Chains() (cs []*FilterChain) {
	ccs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterChain)(nil))](*C.AVFilterChain))(unsafe.Pointer(fgs.c.chains))
	for i := 0; i < fgs.NbChains(); i++ {
		cs = append(cs, newFilterChainFromC(ccs[i]))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraphSegment.html#ab7563eca151d89e693f6258de5ce0214
func (fgs *FilterGraphSegment) NbChains() int {
	return int(fgs.c.nb_chains)
}

func (fgs *FilterGraphSegment) ApplyX() (inputs, outputs *FilterInOut, err error) {
	var ic *C.AVFilterInOut
	var oc *C.AVFilterInOut
	fgs.g.resetLog()
	if err := fgs.g.newError(C.avfilter_graph_segment_apply(fgs.c, 0, &ic, &oc)); err != nil {
		return nil, nil, err
	}
	return newFilterInOutFromC(ic), newFilterInOutFromC(oc), nil
}

func (fgs *FilterGraphSegment) Apply(inputs, outputs *FilterInOut) error {
	var ic **C.AVFilterInOut
	if inputs != nil {
		ic = &inputs.c
	}
	var oc **C.AVFilterInOut
	if outputs != nil {
		oc = &outputs.c
	}
	fgs.g.resetLog()
	if err := fgs.g.newError(C.avfilter_graph_segment_apply(fgs.c, 0, ic, oc)); err != nil {
		return err
	}
	return nil
}

func (fgs *FilterGraphSegment) CreateFilters() error {
	fgs.g.resetLog()
	if err := fgs.g.newError(C.avfilter_graph_segment_create_filters(fgs.c, 0)); err != nil {
		return err
	}
	return nil
}

func (fgs *FilterGraphSegment) ApplyOptions() error {
	fgs.g.resetLog()
	if err := fgs.g.newError(C.avfilter_graph_segment_apply_opts(fgs.c, 0)); err != nil {
		return err
	}
	return nil
}
