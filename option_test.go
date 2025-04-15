package astiav

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleOptions() {
	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	opts := fc.Class().Options()
	for _, oo := range opts.List() {
		if oo.Type() != OptionTypeConst {
			val, err := opts.Get(oo.Name(), NewOptionSearchFlags())
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			fmt.Printf("Option: %s => %s\n", oo, val)
		} else {
			fmt.Printf("Option: %s\n", oo)
		}
	}
}

func TestOption(t *testing.T) {
	const childClassName = "mov/mp4/tgp/psp/tg2/ipod/ismv/f4v muxer"
	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	require.NoError(t, err)

	containsAllNames := func(t *testing.T, ops []*Option, want ...string) {
		var got []string
		for _, o := range ops {
			got = append(got, o.Name())
		}
		for _, name := range want {
			require.Contains(t, want, name)
		}
	}
	testGetSet := func(t *testing.T, os *Options) {
		require.NotNil(t, os)
		l := os.List()
		require.Greater(t, len(l), 50)

		require.Error(t, os.Set("invalid", "", NewOptionSearchFlags()))
		const name = "brand"
		const value = "test"
		require.NoError(t, os.Set(name, value, NewOptionSearchFlags()))
		v, err := os.Get(name, NewOptionSearchFlags())
		require.NoError(t, err)
		require.Equal(t, value, v)
	}
	t.Run("Options", func(t *testing.T) {
		opts := fc.Class().Options()
		require.NotNil(t, opts)
		ol := opts.List()
		containsAllNames(t, ol, "flush_packets", "ignidx", "genpts", "nofillin")
	})
	t.Run("ChildOptions", func(t *testing.T) {
		var childClass *Class
		for _, class := range FindChildClasses(fc) {
			if class.Name() == childClassName {
				childClass = class
				break
			}
		}
		require.NotNil(t, childClass)
		require.Equal(t, childClass.Name(), childClassName)
		containsAllNames(t, childClass.Options().List(), "brand", "frag_size", "moov_size")
		testGetSet(t, childClass.Options())
	})
	t.Run("PrivateData", func(t *testing.T) {
		pd := fc.PrivateData()
		require.NotNil(t, pd)
		require.Equal(t, pd.Class().Name(), childClassName)
		containsAllNames(t, pd.Class().Options().List(), "brand")
		testGetSet(t, pd.Options())
	})
}
