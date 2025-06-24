package astiav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleAllOutputFormats() {
	for _, x := range AllOutputFormats() {
		fmt.Printf("OutputFormat: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}
}

func TestOutputFormat(t *testing.T) {
	formatName := "rawvideo"
	outputFormat := FindOutputFormat(formatName)
	require.NotNil(t, outputFormat)
	require.Equal(t, formatName, outputFormat.Name())
	require.Equal(t, formatName, outputFormat.String())
	require.Equal(t, "raw video", outputFormat.LongName())
	t.Run("AllOutputFormats()", func(t *testing.T) {
		require.GreaterOrEqual(t, len(AllOutputFormats()), 10)
	})
}
