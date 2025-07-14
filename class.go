package astiav

//#include "class.h"
//#include "libavutil/opt.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// Class provides the go-type for an AVClass. It holds the pointer to the AVClass struct along
// with an unsafe.Pointer of an AVOptions-enabled struct (or in some cases a double pointer to an AVClass)
// which is compatible with AVOption functions and av_log().
//
// For Logging see: https://ffmpeg.org/doxygen/2.2/group__lavu__log.html#details
// For AVOptions see: https://ffmpeg.org/doxygen/7.0/group__avoptions.html
//
// https://ffmpeg.org/doxygen/7.0/structAVClass.html
type Class struct {
	c *C.AVClass
	// ptr contains the C pointer to the AVOptions-enabled struct
	// or any double pointer to an AVClass describing it.
	// That is any C pointer p which can be cast to an AVClass struct
	// using `AVClass class = *(const AVClass**)c`.
	// See: av_get_opt() and av_log()
	ptr unsafe.Pointer
}

// Returns a Class instance from a **C.AVClass type.
func newClassFromC(ptr unsafe.Pointer) *Class {
	if ptr == nil {
		return nil
	}
	c := (**C.AVClass)(ptr)
	if c == nil {
		return nil
	}
	return &Class{
		c:   *c,
		ptr: ptr,
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#a5fc161d93a0d65a608819da20b7203ba
func (c *Class) Category() ClassCategory {
	return ClassCategory(C.astiavClassCategory(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#ad763b2e6a0846234a165e74574a550bd
func (c *Class) ItemName() string {
	return C.GoString(C.astiavClassItemName(c.c, c.ptr))
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#aa8883e113a3f2965abd008f7667db7eb
func (c *Class) Name() string {
	return C.GoString(c.c.class_name)
}

// https://ffmpeg.org/doxygen/7.0/structAVClass.html#a88948c8a7c6515181771615a54a808bf
func (c *Class) Parent() *Class {
	return newClassFromC(unsafe.Pointer(C.astiavClassParent(c.c, c.ptr)))
}

func (c *Class) String() string {
	return fmt.Sprintf("%s [%s] @ %p", c.ItemName(), c.Name(), c.ptr)
}

func (c *Class) Options() *Options {
	return &Options{c}
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__mng.html#gabc75970cd87d1bf47a4ff449470e9225
func (c *Class) List() (list []*Option) {
	var prev *C.AVOption
	for {
		o := C.av_opt_next(c.ptr, prev)
		if o == nil {
			return
		}
		list = append(list, newOptionFromC(o))
		prev = o
	}
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga5fd4b92bdf4f392a2847f711676a7537
func (c *Class) Set(name, value string, f OptionSearchFlags) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	classer, done := classers.ensure(c.ptr)
	defer done()
	classer.resetLog()
	return classer.newError(C.av_opt_set(c.ptr, cname, cvalue, C.int(f)))
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__get__funcs.html#gaf31144e60f9ce89dbe8cbea57a0b232c
func (c *Class) Get(name string, f OptionSearchFlags) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var ctemp *C.uint8_t = nil
	classer, done := classers.ensure(c.ptr)
	defer done()
	classer.resetLog()
	if err := classer.newError(C.av_opt_get(c.ptr, cname, C.int(f), &ctemp)); err != nil {
		return "", err
	}
	cvalue := (*C.char)(unsafe.Pointer(ctemp))
	if cvalue == nil {
		return "", nil
	}
	defer C.av_freep(unsafe.Pointer(&cvalue))
	return C.GoString(cvalue), nil
}

// FindChildClasses returns all the child classes of the specified
// Classer object.
func FindChildClasses(c Classer) (out []*Class) {
	cls := c.Class()
	if cls == nil {
		return nil
	}
	out = append(out, cls)
	findChildClasses(cls, true, func(child *Class) bool {
		out = append(out, child)
		return true
	})
	return out
}

func findChildClasses(c *Class, recurse bool, f func(c *Class) bool) {
	if c == nil || c.ptr == nil {
		panic("invalid")
	}
	var childPtr unsafe.Pointer
	for {
		childPtr = C.av_opt_child_next(c.ptr, childPtr)
		if childPtr == nil {
			break
		}
		child := newClassFromC(childPtr)
		if !f(child) {
			break
		}
		if recurse {
			findChildClasses(child, recurse, f)
		}
	}
}

// Classer is any go type that wraps a native type that supports
// AVOptions and AV Logging.
//
// Note this is any C struct of which the first field is a pointer to an AVClass
// struct (e.g. AVCodecContext, AVFormatContext etc). This allows functions like
// av_opt_get to retrieve the AVClass by casting the struct as a double pointer
// of AVClass. For some this is constructed by simply taking the address of the
// AVClass member variable.
type Classer interface {
	// Class returns the wrapped AVClass instance of this Context object.
	Class() *Class

	// ErrMode sets the mode for handling logs passed to the Classer. Messages
	// consumed by the classer are used to build an error value using newError().
	ErrMode(m ErrMode)

	// resetLog will clear the message buffer before calling downstream system.
	resetLog()

	// newError constructs an error value from the return value and
	// any log messages that have been captured since the last call
	// to resetLog() or newError().
	newError(ret C.int) error

	// Returns true if the message was consumed and should not be passed
	// to the general log handler.
	handleLog(l LogLevel, msg string) bool
}

// SetDefaultClasserErrMode sets the default ErrMode
// for any classer object when
func SetDefaultClasserErrMode(mode ErrMode) {
	errModeGlobalDefault = mode
}

type classerState struct {
	messages []string
	mode     ErrMode
}

// ErrMode defines how a Classer builds a error value from a return code,
// optionally using captured log messages to build context from the error code.
type ErrMode struct {
	Level   LogLevel
	Consume bool
}

func (e ErrMode) String() string {
	return fmt.Sprintf("ErrMode{@%s,always=%t}", e.Level, e.Consume)
}
func (e ErrMode) IsZero() bool {
	return e.Level == 0 && !e.Consume
}

var (
	// ErrModeNoLogs uses no logs when building error value from a return code.
	ErrModeNoLogs = ErrMode{
		Level:   LogLevelQuiet,
		Consume: false,
	}
	// ErrModeConsumeLog consumes log message of ERROR or lower when building and error
	// value. Captured log message are not passed to any error handler.
	ErrModeConsumeLog = ErrMode{
		Level:   LogLevelError,
		Consume: true,
	}
	// ErrModeUseLog collects log message of ERROR or lower when building and error
	// value. Matching error messages are also passed to any error handler.
	ErrModeUseLog = ErrMode{
		Level:   LogLevelError,
		Consume: false,
	}
	errModeGlobalDefault = ErrModeNoLogs
)

func (h *classerState) resetLog()         { h.messages = nil }
func (h *classerState) ErrMode(m ErrMode) { h.mode = m }
func (h *classerState) handleLog(l LogLevel, msg string) (handled bool) {
	mode := h.mode
	if mode.IsZero() {
		mode = errModeGlobalDefault
	}
	if mode.Level < l {
		return false
	}
	h.messages = append(h.messages, msg)
	return mode.Consume
}

func (h *classerState) newError(ret C.int) error {
	i := int(ret)
	if i >= 0 {
		return nil
	}
	msg := h.messages
	h.messages = nil
	return &loggedError{Error(ret), msg}
}

var _ Classer = (*UnknownClasser)(nil)

type UnknownClasser struct {
	classerState
	c *Class
}

func newUnknownClasser(ptr unsafe.Pointer) *UnknownClasser {
	return &UnknownClasser{c: newClassFromC(ptr)}
}

func (c *UnknownClasser) Class() *Class {
	return c.c
}

var _ Classer = (*ClonedClasser)(nil)

type ClonedClasser struct {
	classerState
	c *Class
}

func newClonedClasser(c Classer) *ClonedClasser {
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return &ClonedClasser{c: newClassFromC(cl.ptr)}
}

func (c *ClonedClasser) Class() *Class {
	return c.c
}

var classers = newClasserPool()

type classerPool struct {
	pm sync.Map
}

func newClasserPool() *classerPool {
	return &classerPool{}
}

func (p *classerPool) unsafePointer(c Classer) (unsafe.Pointer, bool) {
	if c == nil {
		return nil, false
	}
	cl := c.Class()
	if cl == nil {
		return nil, false
	}
	return cl.ptr, true
}

func (p *classerPool) set(c Classer) {
	if ptr, ok := p.unsafePointer(c); ok {
		p.pm.Store(ptr, c)
	}
}

func (p *classerPool) del(c Classer) {
	if ptr, ok := p.unsafePointer(c); ok {
		p.pm.Delete(ptr)
	}
}

// get returns a Classer and true if a value was found for the specified ptr
// value, otherwise it returns nil and false.
func (p *classerPool) get(ptr unsafe.Pointer) (Classer, bool) {
	if c, ok := p.pm.Load(ptr); ok {
		return c.(Classer), ok
	}
	return nil, false
}

// ensure returns the classer for the specified pointer or allocates an
// UnknownClasser for the ptr and returns that with a function to remove it from
// the pool when complete.
func (p *classerPool) ensure(ptr unsafe.Pointer) (c Classer, done func()) {
	done = func() {}
	val, exists := p.pm.LoadOrStore(ptr, newUnknownClasser(ptr))
	if !exists {
		done = func() { p.pm.Delete(ptr) }
	}
	return val.(Classer), done
}

func (p *classerPool) size() int {
	var i int
	p.pm.Range(func(key, value interface{}) bool {
		i++
		return true
	})
	return i
}
