package astiav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleAllOutputFormats() {
	for _, x := range AllOutputFormats() {
		fmt.Printf("OutputFormat: %s\n", x.Name())
	}
}

func TestOutputFormat(t *testing.T) {
	formatName := "rawvideo"
	outputFormat := FindOutputFormat(formatName)
	require.NotNil(t, outputFormat)
	require.Equal(t, formatName, outputFormat.Name())
	require.Equal(t, formatName, outputFormat.String())
	require.Equal(t, "raw video", outputFormat.LongName())
	require.Equal(t, []string{"yuv", "rgb"}, outputFormat.Extensions())
	t.Run("AllOutputFormats()", func(t *testing.T) {
		require.GreaterOrEqual(t, len(AllOutputFormats()), 10)
	})
	require.Equal(t, []string{"yuv", "rgb"}, outputFormat.Extensions())
}
