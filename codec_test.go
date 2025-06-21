package astiav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleCodecs() {
	for _, x := range Codecs() {
		fmt.Printf("Codec: %s %v\n", x.Name(), x.Class())
		for _, child := range FindClasses(x) {
			fmt.Printf("  Class: %s Item: %s\n", child.Name(), child.ItemName())
		}
	}
}

func TestCodec(t *testing.T) {
	ExampleCodecs()
	t.Run("mp3float", func(t *testing.T) {
		c := FindDecoder(CodecIDMp3)
		require.NotNil(t, c)
		require.Equal(t, c.ID(), CodecIDMp3)
		require.Nil(t, c.ChannelLayouts())
		require.True(t, c.IsDecoder())
		require.False(t, c.IsEncoder())
		require.Nil(t, c.PixelFormats())
		require.Equal(t, []SampleFormat{SampleFormatFltp, SampleFormatFlt}, c.SampleFormats())
		require.Equal(t, "mp3float", c.Name())
		require.Equal(t, "mp3float", c.String())
	})
	t.Run("aac", func(t *testing.T) {
		c := FindDecoderByName("aac")
		require.NotNil(t, c)
		els := []ChannelLayout{
			ChannelLayoutMono,
			ChannelLayoutStereo,
			ChannelLayoutSurround,
			ChannelLayout4Point0,
			ChannelLayout5Point0Back,
			ChannelLayout5Point1Back,
			ChannelLayout7Point1WideBack,
			ChannelLayout6Point1Back,
			ChannelLayout7Point1,
			ChannelLayout22Point2,
			ChannelLayout5Point1Point2Back,
		}
		gls := c.ChannelLayouts()
		require.Len(t, gls, len(els))
		for idx := range els {
			require.True(t, els[idx].Equal(gls[idx]))
		}
		require.True(t, c.IsDecoder())
		require.False(t, c.IsEncoder())
		require.Equal(t, []SampleFormat{SampleFormatFltp}, c.SampleFormats())
		require.Equal(t, "aac", c.Name())
		require.Equal(t, "aac", c.String())
	})
	t.Run("mjpeg/ByID", func(t *testing.T) {
		c := FindEncoder(CodecIDMjpeg)
		require.NotNil(t, c)
		require.False(t, c.IsDecoder())
		require.True(t, c.IsEncoder())
		require.Contains(t, c.PixelFormats(), PixelFormatYuvj420P)
		require.Nil(t, c.SampleFormats())
		require.Contains(t, c.Name(), "mjpeg")
		require.Contains(t, c.String(), "mjpeg")
	})

	t.Run("mjpeg/ByName", func(t *testing.T) {
		c := FindEncoderByName("mjpeg")
		require.NotNil(t, c)
		require.False(t, c.IsDecoder())
		require.True(t, c.IsEncoder())
		require.Equal(t, []PixelFormat{
			PixelFormatYuvj420P,
			PixelFormatYuvj422P,
			PixelFormatYuvj444P,
			PixelFormatYuv420P,
			PixelFormatYuv422P,
			PixelFormatYuv444P,
		}, c.PixelFormats())
		require.Equal(t, "mjpeg", c.Name())
		require.Equal(t, "mjpeg", c.String())
	})
	t.Run("invalid", func(t *testing.T) {
		c := FindDecoderByName("invalid")
		require.Nil(t, c)
	})
	t.Run("Codecs", func(t *testing.T) {
		var found bool
		for _, c := range Codecs() {
			if c.ID() == CodecIDMjpeg {
				found = true
			}
		}
		require.True(t, found)
	})
}
