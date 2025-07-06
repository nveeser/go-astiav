package demo

//#include "enum.h"
import "C"

//go:generate go run ../enumflags.go -flag_types FlagType

type FlagType int

const (
	FlagTypeOne   = FlagType(C.DEMO_FLAG_ONE)
	FlagTypeTwo   = FlagType(C.DEMO_FLAG_TWO)
	FlagTypeThree = FlagType(C.DEMO_FLAG_THREE)
	FlagTypeFour  = FlagType(C.DEMO_FLAG_FOUR)
)
