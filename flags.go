package astiav

func Flags[F ~int64](fs ...F) F {
	var flag F
	for _, f := range fs {
		flag |= f
	}
	return flag
}

type (
	// BuffersinkFlags is an alias for BuffersinkFlag
	// Deprecated - use BuffersinkFlag
	BuffersinkFlags = BuffersinkFlag
	// BuffersrcFlags is an alias for BuffersrcFlag
	// Deprecated - use BuffersrcFlag
	BuffersrcFlags = BuffersrcFlag
	// CodecContextFlags is an alias for CodecContextFlag
	// Deprecated - use CodecContextFlag
	CodecContextFlags = CodecContextFlag
	// CodecContextFlags2 is an alias for CodecContextFlag2
	// Deprecated - use CodecContextFlag2
	CodecContextFlags2 = CodecContextFlag2
	// CodecHardwareConfigMethodFlags is an alias for CodecHardwareConfigMethodFlag
	// Deprecated - use CodecHardwareConfigMethodFlag directly
	CodecHardwareConfigMethodFlags = CodecHardwareConfigMethodFlag
	// DictionaryFlags is an alias for DictionaryFlag
	// Deprecated - use DictionaryFlag directly
	DictionaryFlags = DictionaryFlag
	// FilterFlags is an alias for FilterFlag
	// Deprecated - use FilterFlag directly
	FilterFlags = FilterFlag
	// FilterCommandFlags is an alias for FilterCommandFlag
	// Deprecated - use FilterCommandFlag directly
	FilterCommandFlags = FilterCommandFlag
	// FormatContextFlags is an alias for FormatContextFlag
	// Deprecated - use FormatContextFlag directly
	FormatContextFlags = FormatContextFlag
	// FormatContextCtxFlags is an alias for FormatContextCtxFlag
	// Deprecated - use FormatContextCtxFlag directly
	FormatContextCtxFlags = FormatContextCtxFlag
	// FormatEventFlags is an alias for FormatEventFlag
	// Deprecated - use FormatEventFlag directly
	FormatEventFlags = FormatEventFlag
	// IOContextFlags is an alias for IOContextFlag
	// Deprecated - use IOContextFlag directly
	IOContextFlags = IOContextFlag
	// IOFormatFlags is an alias for IOFormatFlag
	// Deprecated - use IOFormatFlag directly
	IOFormatFlags = IOFormatFlag
	// OptionSearchFlags is an alias for OptionSearchFlag
	// Deprecated - use OptionSearchFlag directly
	OptionSearchFlags = OptionSearchFlag
	// PacketFlags is an alias for PacketFlag
	// Deprecated - use PacketFlag directly
	PacketFlags = PacketFlag
	// PixelFormatDescriptorFlags is an alias for PixelFormatDescriptorFlag
	// Deprecated - use PixelFormatDescriptorFlag directly
	PixelFormatDescriptorFlags = PixelFormatDescriptorFlag
	// SeekFlags is an alias for SeekFlag
	// Deprecated - use SeekFlag directly
	SeekFlags = SeekFlag
	// SoftwareScaleContextFlags is an alias for SoftwareScaleContextFlag
	// Deprecated - use SoftwareScaleContextFlag directly
	SoftwareScaleContextFlags = SoftwareScaleContextFlag
	// StreamEventFlags is an alias for StreamEventFlag
	// Deprecated - use StreamEventFlag directly
	StreamEventFlags = StreamEventFlag
)

var (
	NewBuffersinkFlags                = Flags[BuffersinkFlag]
	NewBuffersrcFlags                 = Flags[BuffersrcFlag]
	NewCodecContextFlags              = Flags[CodecContextFlag]
	NewCodecContextFlags2             = Flags[CodecContextFlag2]
	NewCodecHardwareConfigMethodFlags = Flags[CodecHardwareConfigMethodFlag]
	NewDictionaryFlags                = Flags[DictionaryFlag]
	NewFilterFlags                    = Flags[FilterFlag]
	NewFilterCommandFlags             = Flags[FilterCommandFlag]
	NewFormatContextFlags             = Flags[FormatContextFlag]
	NewFormatContextCtxFlags          = Flags[FormatContextCtxFlag]
	NewFormatEventFlags               = Flags[FormatEventFlag]
	NewIOContextFlags                 = Flags[IOContextFlag]
	NewIOFormatFlags                  = Flags[IOFormatFlag]
	NewOptionSearchFlags              = Flags[OptionSearchFlag]
	NewPacketFlags                    = Flags[PacketFlag]
	NewPixelFormatDescriptorFlags     = Flags[PixelFormatDescriptorFlag]
	NewSeekFlags                      = Flags[SeekFlag]
	NewSoftwareScaleContextFlags      = Flags[SoftwareScaleContextFlag]
	NewStreamEventFlags               = Flags[StreamEventFlag]
)
