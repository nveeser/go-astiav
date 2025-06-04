package astiav

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	for _, x := range AllOutputFormats() {
		fmt.Printf("OutputFormat: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}
	for _, x := range AllInputFormats() {
		fmt.Printf("InputFormat: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}
	for _, x := range Codecs() {
		fmt.Printf("Codec: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}

	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	require.NoError(t, err)

	classes := map[string]*Class{}
	for _, child := range FindClasses(fc) {
		classes[child.Name()] = child
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
