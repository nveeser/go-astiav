package astiav

//#include "class.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVClass.html
type Class struct {
	c   *C.AVClass
	ptr unsafe.Pointer
}

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

// https://www.ffmpeg.org/doxygen/7.0/group__opt__mng.html#gabc75970cd87d1bf47a4ff449470e9225
func (c *Class) Options() (list []*Option) {
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
func (c *Class) SetOption(name, value string, f OptionSearchFlags) error {
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
func (c *Class) GetOption(name string, f OptionSearchFlags) (string, error) {
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

// FindClasses returns all the child classes of the specified
// Classer object.
func FindClasses(c Classer) (out []*Class) {
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

type Classer interface {
	// Class returns the wrapped AVClass instance of this Context object.
	Class() *Class

	// LogMode sets the mode for handling logs passed to the Classer. Messages
	// consumed by the classer are used to build an error value using newError().
	LogMode(mode LogMode)

	// resetLog will clear the message buffer before calling downstream system.
	resetLog()

	// Returns true if the msg should also be passed to the log handler
	handleLog(l LogLevel, msg string) bool
	// newError constructs an error value from the return value and
	// any log messages that have been captured since the last call
	// to resetLog() or newError().
	newError(ret C.int) error
}

var defaultLogMode = LogModeIgnore

func SetDefaultClasserLogMode(m LogMode) {
	defaultLogMode = m
}

// LogMode defines how log messages are handled when passed to a specific Classer
// instance. All log message from libav* use av_log which takes an optional
// context containing an AVClass instance.
type LogMode int

const (
	LogModeDefault LogMode = iota

	// LogModeConsume specifies that this Classer will handle all log messages and
	// none will be passed to the installed Log handler. Messages are included in the
	// context of any errors constructed using newError()
	LogModeConsume

	// LogModeIgnore specifies that this Classer will ignore log messages and all
	// will be passed to the log handler.
	LogModeIgnore
)

type classerHandler struct {
	messages []string
	mode     LogMode
}

func (h *classerHandler) resetLog()            { h.messages = nil }
func (h *classerHandler) LogMode(mode LogMode) { h.mode = mode }

func (h *classerHandler) handleLog(l LogLevel, msg string) (handled bool) {
	mode := h.mode
	if mode == LogModeDefault {
		mode = defaultLogMode
	}
	switch h.mode {
	case LogModeDefault:
		fallthrough
	case LogModeIgnore:
		return false
	case LogModeConsume:
		h.messages = append(h.messages, msg)
		return true
	default:
		return false
	}
}

func (h *classerHandler) newError(ret C.int) error {
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
	classerHandler
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
	classerHandler
	c *Class
}

func newClonedClasser(c Classer) *ClonedClasser {
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return &ClonedClasser{c: newClassFromC(cl.ptr)}
}

func (c *ClonedClasser) Class() *Class { return c.c }

var classers = newClasserPool()

type classerPool struct {
	m sync.Mutex
	p map[unsafe.Pointer]Classer
}

func newClasserPool() *classerPool {
	return &classerPool{p: make(map[unsafe.Pointer]Classer)}
}

func (p *classerPool) unsafePointer(c Classer) unsafe.Pointer {
	if c == nil {
		return nil
	}
	cl := c.Class()
	if cl == nil {
		return nil
	}
	return cl.ptr
}

func (p *classerPool) set(c Classer) {
	p.m.Lock()
	defer p.m.Unlock()
	if ptr := p.unsafePointer(c); ptr != nil {
		p.p[ptr] = c
	}
}

func (p *classerPool) del(c Classer) {
	p.m.Lock()
	defer p.m.Unlock()
	if ptr := p.unsafePointer(c); ptr != nil {
		delete(p.p, ptr)
	}
}

func (p *classerPool) get(ptr unsafe.Pointer) (Classer, bool) {
	p.m.Lock()
	defer p.m.Unlock()
	c, ok := p.p[ptr]
	return c, ok
}
