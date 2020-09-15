package test

/*
Excerpt from https://github.com/purelabio/eth/blob/master/abi.go

This is a separate package, instead of simply being included in the test file,
in order to test how the printer handles package names.
*/

type Abi []AbiMethod

type AbiMethod interface{}

type AbiConstructor struct {
	Type            string
	Name            string
	Inputs          []AbiParam
	Payable         bool
	StateMutability string
}

type AbiFunction struct {
	Type            string
	Name            string
	Constant        bool
	Inputs          []AbiParam
	Outputs         []AbiParam
	Payable         bool
	StateMutability string
	Selector        [4]byte
}

type AbiEvent struct {
	Type             string
	Name             string
	Inputs           []AbiParam
	Anonymous        bool
	Selector         Word
	IndexedInputs    []AbiParam
	NonIndexedInputs []AbiParam
}

type AbiParam struct {
	Name       string
	Type       string
	Components []AbiParam
	Indexed    bool
	AbiType    AbiType
}

type AbiKind byte

const (
	AbiKindBool AbiKind = iota + 1
	AbiKindUint
	AbiKindInt
	AbiKindAddress
	AbiKindFunction
	AbiKindDenseArray
	AbiKindSparseArray
)

func (self AbiKind) String() string {
	switch self {
	case AbiKindBool:
		return "AbiKindBool"
	case AbiKindUint:
		return "AbiKindUint"
	case AbiKindInt:
		return "AbiKindInt"
	case AbiKindAddress:
		return "AbiKindAddress"
	case AbiKindFunction:
		return "AbiKindFunction"
	case AbiKindDenseArray:
		return "AbiKindDenseArray"
	case AbiKindSparseArray:
		return "AbiKindSparseArray"
	default:
		return ""
	}
}

type AbiType struct {
	Type     string
	Kind     AbiKind
	ArrayLen int
	FixedLen bool
	Elem     *AbiType
}

type Word [32]byte
