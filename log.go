package astiav

//#include "log.h"
//#include <libavutil/log.h>
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/group__lavu__log__constants.html#ga11e329935b59b83ca722b66674f37fd4
type LogLevel int

const (
	LogLevelQuiet   = LogLevel(C.AV_LOG_QUIET)
	LogLevelPanic   = LogLevel(C.AV_LOG_PANIC)
	LogLevelFatal   = LogLevel(C.AV_LOG_FATAL)
	LogLevelError   = LogLevel(C.AV_LOG_ERROR)
	LogLevelWarning = LogLevel(C.AV_LOG_WARNING)
	LogLevelInfo    = LogLevel(C.AV_LOG_INFO)
	LogLevelVerbose = LogLevel(C.AV_LOG_VERBOSE)
	LogLevelDebug   = LogLevel(C.AV_LOG_DEBUG)
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelPanic:
		return "panic"
	case LogLevelFatal:
		return "fatal"
	case LogLevelError:
		return "error"
	case LogLevelWarning:
		return "warning"
	case LogLevelInfo:
		return "info"
	case LogLevelVerbose:
		return "verbose"
	case LogLevelDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__log.html#ga1fd32c74db581e3e2e7f35d277bb1e24
func SetLogLevel(l LogLevel) {
	C.av_log_set_level(C.int(l))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__log.html#gae8ada5cc5722548d8698650b05207904
func GetLogLevel() LogLevel {
	return LogLevel(C.av_log_get_level())
}

type LogCallback func(c Classer, l LogLevel, fmt, msg string)

var logCallback LogCallback

// https://ffmpeg.org/doxygen/7.0/group__lavu__log.html#ga14034761faf581a8b9ed6ef19b313708
func SetLogCallback(c LogCallback) {
	logCallback = c
	C.astiavSetLogCallback()
}

//export goAstiavLogCallback
func goAstiavLogCallback(logCtx unsafe.Pointer, level C.int, msg *C.char) {
	cls, done := findLoggingClasser(logCtx)
	defer done()
	handleLog(cls, LogLevel(level), "", C.GoString(msg))
}

func findLoggingClasser(ptr unsafe.Pointer) (Classer, func()) {
	for ptr != nil {
		if cls, ok := classers.get(ptr); ok {
			return cls, func() {}
		}
		parent := newClassFromC(ptr).Parent()
		if parent == nil {
			break
		}
		ptr = newClassFromC(ptr).Parent().ptr
	}
	return classers.ensure(ptr)
}

func handleLog(c Classer, l LogLevel, format, msg string) {
	if c != nil {
		if c.handleLog(l, msg) {
			return
		}
	}
	if logCallback != nil {
		logCallback(c, l, format, msg)
	}
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__log.html#ga5bd132d2e4ac6f9843ef6d8e3c05050a
func ResetLogCallback() {
	C.astiavResetLogCallback()
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__log.html#gabd386ffd4b27637cf34e98d5d1a6e8ae
func Log(c Classer, l LogLevel, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msgc := C.CString(msg)
	defer C.free(unsafe.Pointer(msgc))

	var ptr unsafe.Pointer
	if c != nil {
		if cl := c.Class(); cl != nil {
			ptr = cl.ptr
		}
	}
	C.astiavLog(ptr, C.int(l), msgc)
}
