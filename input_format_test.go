package astiav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleAllInputFormats() {
	for _, x := range AllInputFormats() {
		fmt.Printf("InputFormat: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}
}
func TestProbeInputFormat(t *testing.T) {
	format, score, err := ProbeInputFormat("testdata/video.mp4")
	require.Nil(t, err)
	require.Greater(t, score, 0)
	require.NotNil(t, format)
}

func TestInputFormat(t *testing.T) {
	formatName := "rawvideo"
	inputFormat := FindInputFormat(formatName)
	require.NotNil(t, inputFormat)
	require.Equal(t, formatName, inputFormat.Name())
	require.Equal(t, formatName, inputFormat.String())
	require.Equal(t, "raw video", inputFormat.LongName())
	t.Run("AllInputFormats()", func(t *testing.T) {
		require.GreaterOrEqual(t, len(AllInputFormats()), 10)
	})
}
