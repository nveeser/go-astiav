package astiav

//#include <libavformat/avformat.h>
import "C"
import (
	"os"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVInputFormat.html
type InputFormat struct {
	classerHandler
	c *C.AVInputFormat
}

func newInputFormatFromC(c *C.AVInputFormat) *InputFormat {
	if c == nil {
		return nil
	}
	return &InputFormat{c: c}
}

func ProbeInputFormat(filename string) (*InputFormat, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()
	ioContext, err := AllocIOContext(4096, false, f.Read, f.Seek, nil)
	if err != nil {
		return nil, 0, err
	}

	var format *C.AVInputFormat
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	ioContext.resetLog()
	ret := C.av_probe_input_buffer2(ioContext.c, &format, cfilename, nil, 0, 0)
	if err = ioContext.newError(ret); err != nil {
		return nil, 0, err
	}
	return newInputFormatFromC(format), int(ret), nil
}

func AllInputFormats() []*InputFormat {
	var out []*InputFormat
	var iter unsafe.Pointer
	var iformat *C.AVInputFormat
	for {
		iformat = C.av_demuxer_iterate(&iter)
		if iformat == nil {
			break
		}
		out = append(out, newInputFormatFromC(iformat))
	}
	return out
}

func (f *InputFormat) Class() *Class {
	if f.c.priv_class != nil {
		priv_class := f.c.priv_class
		return &Class{unsafe.Pointer(&priv_class), f.c.priv_class}
	}
	return nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#ga40034b6d64d372e1c989e16dde4b459a
func FindInputFormat(name string) *InputFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newInputFormatFromC(C.av_find_input_format(cname))
}

// https://ffmpeg.org/doxygen/7.0/structAVInputFormat.html#a1b30f6647d0c2faf38ba8786d7c3a838
func (f *InputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVInputFormat.html#a850db3eb225e22b64f3304d72134ca0c
func (f *InputFormat) Name() string {
	return C.GoString(f.c.name)
}

// https://ffmpeg.org/doxygen/7.0/structAVInputFormat.html#a1f67064a527941944017f1dfe65d3aa9
func (f *InputFormat) LongName() string {
	return C.GoString(f.c.long_name)
}

func (f *InputFormat) String() string {
	return f.Name()
}
