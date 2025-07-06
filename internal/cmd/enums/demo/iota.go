package demo

//go:generate go run ../enumflags.go -enum_types IotaType

type IotaType int

const (
	IotaTypeOne IotaType = iota
	IotaTypeTwo
	IotaTypeThree
	IotaTypeFour

	ValueOne   = 10
	ValueTwo   = 20
	ValueThree = 30
)
