package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"fmt"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterLink.html
type FilterLink struct {
	c *C.AVFilterLink
}

func (l *FilterLink) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fmt.Fprintf(f, "%s [%s] -> %s [%s] %s",
			l.Src().Name(), l.SrcPad().Name(),
			l.Dst().Name(), l.DstPad().Name(),
			l.MediaType())
		if f.Flag('+') {
			switch l.MediaType() {
			case MediaTypeVideo:
				fmt.Fprintf(f, " fmt:%s w:%d h:%d tb:%s",
					l.PixelFormat(), l.Width(), l.Height(), l.TimeBase())
			case MediaTypeAudio:
				fmt.Fprintf(f, " fmt:%s sr:%d cl:%s tb:%s",
					l.SampleFormat(), l.SampleRate(), l.ChannelLayout(), l.TimeBase())
			}
		}
	default:
		fmt.Fprintf(f, "%p", l.c)
	}
}

func newFilterLinkC(c *C.AVFilterLink) *FilterLink {
	return &FilterLink{c: c}
}

func (l *FilterLink) MediaType() MediaType {
	return MediaType(l.c._type)
}

func (l *FilterLink) Src() *FilterContext {
	return newFilterContext(l.c.src)
}

func (l *FilterLink) SrcPad() *FilterPad {
	if l == nil || l.c == nil {
		return nil
	}
	return newFilterPad(l.c.srcpad, 0)
}

func (l *FilterLink) Dst() *FilterContext {
	if l == nil || l.c == nil {
		return nil
	}
	return newFilterContext(l.c.dst)
}

func (l *FilterLink) DstPad() *FilterPad {
	return newFilterPad(l.c.dstpad, 0)
}

func (l *FilterLink) PixelFormat() PixelFormat {
	if l.MediaType() != MediaTypeVideo {
		panic("PixelFormat() called for incorrect MediaType: " + l.MediaType().String())
	}
	return PixelFormat(l.c.format)
}

func (l *FilterLink) SampleFormat() SampleFormat {
	if l.MediaType() != MediaTypeAudio {
		panic("SampleFormat() called for incorrect MediaType: " + l.MediaType().String())
	}
	return SampleFormat(l.c.format)
}

// Width returns the width of the video frames.
func (l *FilterLink) Width() int {
	return int(l.c.w)
}

// Height returns the height of the video frames.
func (l *FilterLink) Height() int {
	return int(l.c.h)
}

// TimeBase returns the time base of the video or audio frames.
func (l *FilterLink) TimeBase() Rational {
	return newRationalFromC(l.c.time_base)
}

// ChannelLayout returns the channel layout of the audio frames.
func (l *FilterLink) ChannelLayout() ChannelLayout {
	return newChannelLayoutFromC(&l.c.ch_layout)
}

// SampleRate returns the sample rate of the audio frames.
func (l *FilterLink) SampleRate() int {
	return int(l.c.sample_rate)
}
