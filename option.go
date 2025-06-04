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
