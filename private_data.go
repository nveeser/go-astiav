package astiav

import "unsafe"

// Deprecated - use FindChildClasses() for access AVClass instances
// contained within an existing AVOptions enabled struct.
type PrivateData struct {
	c unsafe.Pointer
}

func newPrivateDataFromC(c unsafe.Pointer) *PrivateData {
	if c == nil {
		return nil
	}
	return &PrivateData{c: c}
}

// Deprecated - use Class().GetOptions()
func (pd *PrivateData) Options() *Options {
	return &Options{pd.Class()}
}

func (pd *PrivateData) Class() *Class {
	return newClassFromC(pd.c)
}
