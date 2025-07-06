package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/avformat_8h.html#ab3a5958310f614671f5030ed10753ba9
type StreamEventFlag int64

//go:generate go run internal/cmd/enums/enumflags.go -flag_types StreamEventFlag

const (
	StreamEventFlagMetadataUpdated = StreamEventFlag(C.AVSTREAM_EVENT_FLAG_METADATA_UPDATED)
)
