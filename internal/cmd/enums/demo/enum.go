package demo

//#include "enum.h"
import "C"

//go:generate go run ../enumflags.go -enum_types EnumType

type EnumType int

const (
	EnumTypeOne   = EnumType(C.DEMO_ENUM_ONE)
	EnumTypeTwo   = EnumType(C.DEMO_ENUM_TWO)
	EnumTypeThree = EnumType(C.DEMO_ENUM_THREE)
	EnumTypeFour  = EnumType(C.DEMO_ENUM_FOUR)
)
