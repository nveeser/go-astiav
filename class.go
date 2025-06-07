package astiav

//#include "class.h"
//#include "libavutil/opt.h"
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVClass.html
type Class struct {
	ptr unsafe.Pointer
	c   *C.AVClass
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

func (c *Class) OptionNames() []string {
	var n []string
	for _, o := range c.Options() {
		n = append(n, o.Name())
	}
	return n
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga5fd4b92bdf4f392a2847f711676a7537
func (c *Class) SetOption(name, value string, f OptionSearchFlags) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	classer, done := classers.get(c.ptr)
	defer done()
	classer.resetLog()
	return classer.newError(C.av_opt_set(c.ptr, cname, cvalue, C.int(f)))
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__get__funcs.html#gaf31144e60f9ce89dbe8cbea57a0b232c
func (c *Class) GetOption(name string, f OptionSearchFlags) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var ctemp *C.uint8_t = nil
	classer, done := classers.get(c.ptr)
	defer done()
	classer.resetLog()
	if err := classer.newError(C.av_opt_get(c.ptr, cname, C.int(f), &ctemp)); err != nil {
		return "", err
	}
	cvalue := (*C.char)(ctemp)
	if cvalue == nil {
		return "", nil
	}
	defer C.av_freep(unsafe.Pointer(&cvalue))
	return C.GoString(cvalue), nil
}

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

type classerHandler struct {
	messages []string
}

func (h *classerHandler) handleLog(l LogLevel, msg string) {
	if 0 <= l && l <= LogLevelError {
		h.messages = append(h.messages, msg)
	}
}

func (h *classerHandler) resetLog() {
	h.messages = nil
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

type Classer interface {
	Class() *Class
	resetLog()
	handleLog(l LogLevel, msg string)
	newError(ret C.int) error
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
	pm sync.Map
}

func newClasserPool() *classerPool {
	return &classerPool{}
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
	if ptr := p.unsafePointer(c); ptr != nil {
		p.pm.Store(ptr, c)
	}
}

func (p *classerPool) del(c Classer) {
	if ptr := p.unsafePointer(c); ptr != nil {
		p.pm.Delete(ptr)
	}
}

func (p *classerPool) find(ptr unsafe.Pointer) (Classer, bool) {
	if c, ok := p.pm.Load(ptr); ok {
		return c.(Classer), ok
	}
	return nil, false
}

func (p *classerPool) get(ptr unsafe.Pointer) (c Classer, done func()) {
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
