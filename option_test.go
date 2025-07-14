package astiav

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	require.NoError(t, err)

	t.Run("FindClasses", func(t *testing.T) {
		classes := map[string]*Class{}
		for _, class := range FindClasses(fc) {
			fmt.Printf("  Class: %s Item: %s\n", class.Name(), class.ItemName())
			classes[class.Name()] = class
		}
		require.Len(t, classes, 2)
		require.Contains(t, slices.Collect(maps.Keys(classes)), "AVFormatContext", "mov/mp4/tgp/psp/tg2/ipod/ismv/f4v muxer")
	})
	t.Run("InnerClass", func(t *testing.T) {
		var muxerClass *Class
		for _, class := range FindClasses(fc) {
			if class.Name() == "mov/mp4/tgp/psp/tg2/ipod/ismv/f4v muxer" {
				muxerClass = class
				break
			}
		}
		require.NotNil(t, muxerClass)
		l := muxerClass.Options()
		require.Len(t, l, 56)
		const name = "brand"
		o := l[0]
		require.Equal(t, name, o.Name())
		_, err = muxerClass.GetOption("invalid", NewOptionSearchFlags())
		require.Error(t, err)
		v, err := muxerClass.GetOption(name, NewOptionSearchFlags())
		require.NoError(t, err)
		require.Equal(t, "", v)
		require.Error(t, muxerClass.SetOption("invalid", "", NewOptionSearchFlags()))
		const value = "test"
		require.NoError(t, muxerClass.SetOption(name, value, NewOptionSearchFlags()))
		v, err = muxerClass.GetOption(name, NewOptionSearchFlags())
		require.NoError(t, err)
		require.Equal(t, value, v)
	})
}
