package astiav

//#include <libavutil/opt.h>
//#include "option.h"
import "C"
import (
	"github.com/asticode/go-astikit"
	"strings"
)

type OptionFlag int64

const (
	OptionFlagEncodingParam  = OptionFlag(C.AV_OPT_FLAG_ENCODING_PARAM)
	OptionFlagDecodingParam  = OptionFlag(C.AV_OPT_FLAG_DECODING_PARAM)
	OptionFlagAudioParam     = OptionFlag(C.AV_OPT_FLAG_AUDIO_PARAM)
	OptionFlagVideoParam     = OptionFlag(C.AV_OPT_FLAG_VIDEO_PARAM)
	OptionFlagSubtitleParam  = OptionFlag(C.AV_OPT_FLAG_SUBTITLE_PARAM)
	OptionFlagExport         = OptionFlag(C.AV_OPT_FLAG_EXPORT)
	OptionFlagReadonly       = OptionFlag(C.AV_OPT_FLAG_READONLY)
	OptionFlagBsfParam       = OptionFlag(C.AV_OPT_FLAG_BSF_PARAM)
	OptionFlagRuntimeParam   = OptionFlag(C.AV_OPT_FLAG_RUNTIME_PARAM)
	OptionFlagFilteringParam = OptionFlag(C.AV_OPT_FLAG_FILTERING_PARAM)
	OptionFlagDeprecated     = OptionFlag(C.AV_OPT_FLAG_DEPRECATED)
	OptionFlagChildConsts    = OptionFlag(C.AV_OPT_FLAG_CHILD_CONSTS)
)

func (f OptionFlag) BitSet() astikit.BitFlags {
	return astikit.BitFlags(f)
}

var optionFlagNames = map[OptionFlag]string{
	OptionFlagEncodingParam:  "EncodingParam",
	OptionFlagDecodingParam:  "DecodingParam",
	OptionFlagAudioParam:     "AudioParam",
	OptionFlagVideoParam:     "VideoParam",
	OptionFlagSubtitleParam:  "SubtitleParam",
	OptionFlagExport:         "Export",
	OptionFlagReadonly:       "Readonly",
	OptionFlagBsfParam:       "BsfParam",
	OptionFlagRuntimeParam:   "RuntimeParam",
	OptionFlagFilteringParam: "FilteringParam",
	OptionFlagDeprecated:     "Deprecated",
	OptionFlagChildConsts:    "ChildConsts",
}

func (f OptionFlag) String() string {
	var names []string
	for k, v := range optionFlagNames {
		if f.BitSet().Has(uint64(k)) {
			names = append(names, v)
		}
	}
	return strings.Join(names, "|")
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__dict.html#gad9cbc53cec515b72ae7caa2e194c6bc0
type OptionType int64

func (t OptionType) String() string {
	return optionTypeNameMap[t]
}

const (
	OptionTypeFlags     = OptionType(C.AV_OPT_TYPE_FLAGS)
	OptionTypeInt       = OptionType(C.AV_OPT_TYPE_INT)
	OptionTypeInt64     = OptionType(C.AV_OPT_TYPE_INT64)
	OptionTypeDouble    = OptionType(C.AV_OPT_TYPE_DOUBLE)
	OptionTypeFloat     = OptionType(C.AV_OPT_TYPE_FLOAT)
	OptionTypeString    = OptionType(C.AV_OPT_TYPE_STRING)
	OptionTypeRational  = OptionType(C.AV_OPT_TYPE_RATIONAL)
	OptionTypeBinary    = OptionType(C.AV_OPT_TYPE_BINARY) // offset must point to a pointer immediately followed by an int for the length
	OptionTypeDict      = OptionType(C.AV_OPT_TYPE_DICT)
	OptionTypeUint64    = OptionType(C.AV_OPT_TYPE_UINT64)
	OptionTypeConst     = OptionType(C.AV_OPT_TYPE_CONST)
	OptionTypeImageSize = OptionType(C.AV_OPT_TYPE_IMAGE_SIZE) // offset must point to two consecutive integers
	OptionTypePixelFmt  = OptionType(C.AV_OPT_TYPE_PIXEL_FMT)
	OptionTypeSampleFmt = OptionType(C.AV_OPT_TYPE_SAMPLE_FMT)
	OptionTypeVideoRate = OptionType(C.AV_OPT_TYPE_VIDEO_RATE) // offset must point to AVRational
	OptionTypeDuration  = OptionType(C.AV_OPT_TYPE_DURATION)
	OptionTypeColor     = OptionType(C.AV_OPT_TYPE_COLOR)
	OptionTypeBool      = OptionType(C.AV_OPT_TYPE_BOOL)
	OptionTypeChlayout  = OptionType(C.AV_OPT_TYPE_CHLAYOUT)
	OptionTypeFlagArray = OptionType(C.AV_OPT_TYPE_FLAG_ARRAY)
)

var optionTypeNameMap = map[OptionType]string{
	OptionTypeFlags:     "flags",
	OptionTypeFlagArray: "flag_array",
	OptionTypeBool:      "bool",
	OptionTypeConst:     "const",
	OptionTypeInt:       "int",
	OptionTypeInt64:     "int64",
	OptionTypeUint64:    "uint64",
	OptionTypeDouble:    "double",
	OptionTypeFloat:     "float",
	OptionTypeString:    "string",
	OptionTypeRational:  "Rational",
	OptionTypeBinary:    "binary",
	OptionTypeDict:      "Dictionary",
	OptionTypeImageSize: "image_size",
	OptionTypePixelFmt:  "PixelFormat",
	OptionTypeSampleFmt: "sample_fmt",
	OptionTypeVideoRate: "video_rate",
	OptionTypeDuration:  "duration",
	OptionTypeColor:     "color",
	OptionTypeChlayout:  "channel_layout",
}
