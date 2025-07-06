package demo

import (
	"fmt"
)

func ExampleIotaType() {
	fmt.Printf("String: %s\n", IotaTypeOne)
	fmt.Printf("Number: %d\n", IotaTypeTwo)
	fmt.Printf("Binary: %b\n", IotaTypeThree)
	// Output:
	// String: One
	// Number: 1
	// Binary: 10
}

func ExampleEnumType() {
	fmt.Printf("String: %s\n", EnumTypeOne)
	fmt.Printf("CString: %c\n", EnumTypeOne)
	fmt.Printf("Number: %d\n", EnumTypeThree)
	fmt.Printf("Binary: %b\n", EnumTypeFour)
	// Output:
	// String: One
	// CString: DEMO_ENUM_ONE
	// Number: 3
	// Binary: 100
}

func ExampleFlagType() {
	var flags = FlagTypeOne | FlagTypeTwo
	fmt.Printf("Has One: %t\n", flags.Has(FlagTypeOne))
	fmt.Printf("Has Three: %t\n", flags.Has(FlagTypeThree))

	fmt.Printf("String: %s\n", FlagTypeOne)
	fmt.Printf("CString: %c\n", EnumTypeTwo)
	fmt.Printf("Number: %d\n", EnumTypeThree)
	fmt.Printf("Binary: %b\n", EnumTypeFour)
	fmt.Printf("String (Set): %s\n", FlagTypeOne|FlagTypeTwo)
	fmt.Printf("CString (Set): %c\n", FlagTypeOne|FlagTypeTwo)
	fmt.Printf("Number (Set): %d\n", FlagTypeThree|FlagTypeFour)
	fmt.Printf("Binary (Set): %b\n", FlagTypeThree|FlagTypeFour)
	// Output:
	// Has One: true
	// Has Three: false
	// String: One
	// CString: DEMO_ENUM_TWO
	// Number: 3
	// Binary: 100
	// String (Set): One|Two
	// CString (Set): DEMO_FLAG_ONE|DEMO_FLAG_TWO
	// Number (Set): 48
	// Binary (Set): 110000
}
