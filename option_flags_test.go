package astiav

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionFlag(t *testing.T) {
	t.Run("BitSet", func(t *testing.T) {
		var flag = OptionFlagAudioParam
		assert.Equal(t, true, flag.IsSingleFlag())
		flag = OptionFlagAudioParam | OptionFlagReadonly
		assert.Equal(t, false, flag.IsSingleFlag())
		assert.Equal(t, true, flag.Has(OptionFlagAudioParam))
		assert.Equal(t, false, flag.Has(OptionFlagBsfParam))
		flag = flag.Add(OptionFlagBsfParam)
		assert.Equal(t, true, flag.Has(OptionFlagBsfParam))
	})
	t.Run("FormatSingle", func(t *testing.T) {
		var flag = OptionFlagAudioParam
		assert.Equal(t, "AudioParam", fmt.Sprintf("%s", flag))
		assert.Equal(t, "AudioParam(8)", fmt.Sprintf("%+s", flag))
		assert.Equal(t, "AV_OPT_FLAG_AUDIO_PARAM", fmt.Sprintf("%c", flag))
		assert.Equal(t, "AV_OPT_FLAG_AUDIO_PARAM(8)", fmt.Sprintf("%+c", flag))
		assert.Equal(t, "8", fmt.Sprintf("%d", flag))
		assert.Equal(t, "1000", fmt.Sprintf("%b", flag))
	})
	t.Run("FormatMulti", func(t *testing.T) {
		var flags = OptionFlagAudioParam | OptionFlagReadonly | OptionFlagExport
		assert.Equal(t, "AudioParam|Export|Readonly", fmt.Sprintf("%s", flags))
		assert.Equal(t, "AudioParam(8)|Export(64)|Readonly(128)", fmt.Sprintf("%+s", flags))
		assert.Equal(t, "AV_OPT_FLAG_AUDIO_PARAM|AV_OPT_FLAG_EXPORT|AV_OPT_FLAG_READONLY", fmt.Sprintf("%c", flags))
		assert.Equal(t, "AV_OPT_FLAG_AUDIO_PARAM(8)|AV_OPT_FLAG_EXPORT(64)|AV_OPT_FLAG_READONLY(128)", fmt.Sprintf("%+c", flags))
		assert.Equal(t, "200", fmt.Sprintf("%d", flags))
		assert.Equal(t, "11001000", fmt.Sprintf("%b", flags))
	})
}

func TestOptionType(t *testing.T) {
	t.Run("Format", func(t *testing.T) {
		var flag = OptionTypeChlayout
		assert.Equal(t, "Chlayout", fmt.Sprintf("%s", flag))
		assert.Equal(t, "Chlayout(19)", fmt.Sprintf("%+s", flag))
		assert.Equal(t, "AV_OPT_TYPE_CHLAYOUT", fmt.Sprintf("%c", flag))
		assert.Equal(t, "AV_OPT_TYPE_CHLAYOUT(19)", fmt.Sprintf("%+c", flag))
		assert.Equal(t, "19", fmt.Sprintf("%d", flag))
		assert.Equal(t, "10011", fmt.Sprintf("%b", flag))
	})
}
