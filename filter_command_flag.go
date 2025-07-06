package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gace41bae000b621fc8804a93ce9f2d6e9
type FilterCommandFlag int64

//go:generate go run internal/cmd/enums/enumflags.go -flag_types FilterCommandFlag

const (
	FilterCommandFlagOne  = FilterCommandFlag(C.AVFILTER_CMD_FLAG_ONE)
	FilterCommandFlagFast = FilterCommandFlag(C.AVFILTER_CMD_FLAG_FAST)
)
