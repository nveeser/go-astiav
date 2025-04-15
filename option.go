package astiav

//#include <libavutil/opt.h>
//#include "option.h"
import "C"
import (
	"fmt"
)

// https://www.ffmpeg.org/doxygen/7.0/structAVOption.html
type Option struct {
	c *C.AVOption
}

func newOptionFromC(c *C.AVOption) *Option {
	if c == nil {
		return nil
	}
	return &Option{c: c}
}

// https://www.ffmpeg.org/doxygen/7.0/structAVOption.html#a87e81c6e58d6a94d97a98ad15a4e507c
func (o *Option) Name() string     { return C.GoString(o.c.name) }
func (o *Option) Help() string     { return C.GoString(o.c.help) }
func (o *Option) Type() OptionType { return OptionType(o.c._type) }

func (o *Option) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, "%s[%s]", o.Name(), o.Type())
	if verb == 'v' {
		unit := C.GoString(o.c.unit)
		if unit != "" {
			fmt.Fprintf(s, " unit=%s", unit)
		}
		fmt.Fprintf(s, " flags=%s", OptionFlag(o.c.flags))
	}
}

type Options struct {
	class *Class
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__mng.html#gabc75970cd87d1bf47a4ff449470e9225
func (os *Options) List() (list []*Option) {
	return os.class.List()
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__set__funcs.html#ga5fd4b92bdf4f392a2847f711676a7537
func (os *Options) Set(name, value string, f OptionSearchFlags) error {
	return os.class.Set(name, value, f)
}

// https://www.ffmpeg.org/doxygen/7.0/group__opt__get__funcs.html#gaf31144e60f9ce89dbe8cbea57a0b232c
func (os *Options) Get(name string, f OptionSearchFlags) (string, error) {
	return os.class.Get(name, f)
}
