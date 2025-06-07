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

	classes := map[string]*Class{}
	for _, class := range FindClasses(fc) {
		fmt.Printf("  Class: %s Item: %s\n", class.Name(), class.ItemName())
		classes[class.Name()] = class
	}
	require.Len(t, classes, 2)
	require.Contains(t, slices.Collect(maps.Keys(classes)), "AVFormatContext", "mov/mp4/tgp/psp/tg2/ipod/ismv/f4v muxer")
	os := classes["mov/mp4/tgp/psp/tg2/ipod/ismv/f4v muxer"]
	require.NotNil(t, os)
	l := os.Options()
	require.Len(t, l, 56)
	const name = "brand"
	o := l[0]
	require.Equal(t, name, o.Name())
	_, err = os.GetOption("invalid", NewOptionSearchFlags())
	require.Error(t, err)
	v, err := os.GetOption(name, NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, "", v)
	require.Error(t, os.SetOption("invalid", "", NewOptionSearchFlags()))
	const value = "test"
	require.NoError(t, os.SetOption(name, value, NewOptionSearchFlags()))
	v, err = os.GetOption(name, NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, value, v)
}
