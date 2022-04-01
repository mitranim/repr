package repr

import (
	"encoding/json"
	"go/format"
	"testing"

	"github.com/mitranim/repr/test"
)

func ExampleString() {
	type Data struct {
		Number int
		String string
		List   []int
	}

	_, _ = Println(Data{
		Number: 123,
		String: "hello world!",
		List:   []int{10, 20, 30},
	})

	// Output:
	// repr.Data{
	// 	Number: 123,
	// 	String: "hello world!",
	// 	List: []int{10, 20, 30},
	// }
}

func TestValidSyntax(t *testing.T) {
	code := Bytes(testStructure)
	_, err := format.Source(code)
	if err != nil {
		t.Fatalf("failed to format via gofmt: %v", err)
	}
}

func TestDefault(t *testing.T) {
	actual := String(testStructure)
	expected := testOutputDefault
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestSingleLine(t *testing.T) {
	conf := Config{}
	actual := StringC(testStructure, conf)
	expected := testOutputSingleLine
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestWithoutPackageName(t *testing.T) {
	conf := Config{
		Indent: Default.Indent,
		PackageMap: map[string]string{
			"github.com/mitranim/repr/test": "",
		},
	}
	actual := StringC(testStructure, conf)
	expected := testOutputWithoutPackageName
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestSingleLineWithoutPackageName(t *testing.T) {
	conf := Config{
		PackageMap: map[string]string{
			"github.com/mitranim/repr/test": "",
		},
	}
	actual := StringC(testStructure, conf)
	expected := testOutputSingleLineWithoutPackageName
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestRenamed(t *testing.T) {
	conf := Config{
		Indent: Default.Indent,
		PackageMap: map[string]string{
			"github.com/mitranim/repr/test": "renamed",
		},
	}
	actual := StringC(testStructure, conf)
	expected := testOutputRenamed
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestSingleLineRenamed(t *testing.T) {
	conf := Config{
		PackageMap: map[string]string{
			"github.com/mitranim/repr/test": "renamed",
		},
	}
	actual := StringC(testStructure, conf)
	expected := testOutputSingleLineRenamed
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestZeroFields(t *testing.T) {
	conf := Config{
		Indent:     Default.Indent,
		ZeroFields: true,
	}
	actual := StringC(testStructure, conf)
	expected := testOutputWithZeroFields
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestForceConstructorNames(t *testing.T) {
	conf := Default
	conf.ForceConstructorName = true
	actual := StringC(testStructure, conf)
	expected := testOutputWithForcedConstructorNames
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func TestBytesHex(t *testing.T) {
	actual := String(testBytes)
	expected := testOutputBytesHex
	if actual != expected {
		t.Fatalf("expected output:\n%v\nactual output:\n%v", expected, actual)
	}
}

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Bytes(testStructure)
	}
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String(testStructure)
	}
}

func BenchmarkBytesWithFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := format.Source(Bytes(testStructure))
		if err != nil {
			b.Fatalf("failed to format via gofmt: %v", err)
		}
	}
}

func BenchmarkJsonForComparison(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testStructure)
		if err != nil {
			b.Fatalf("failed to json-marshal: %v", err)
		}
	}
}

var testStructure = test.Abi{
	test.AbiFunction{
		Type:     "function",
		Name:     "two",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "uint256",
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
			test.AbiParam{
				Type: "uint32[]",
				AbiType: test.AbiType{
					Type: "uint32[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint32",
						Kind: 2,
					},
				},
			},
			test.AbiParam{
				Type: "bytes10",
				AbiType: test.AbiType{
					Type:     "bytes10",
					Kind:     6,
					ArrayLen: 10,
					FixedLen: true,
				},
			},
			test.AbiParam{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
		},
		Outputs:         []test.AbiParam{},
		StateMutability: "pure",
		Selector:        [4]uint8{38, 21, 241, 68},
	},
	test.AbiFunction{
		Type:     "function",
		Name:     "one",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
			test.AbiParam{
				Type: "bool",
				AbiType: test.AbiType{
					Type: "bool",
					Kind: 1,
				},
			},
			test.AbiParam{
				Type: "uint256[]",
				AbiType: test.AbiType{
					Type: "uint256[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint256",
						Kind: 2,
					},
				},
			},
		},
		Outputs:         []test.AbiParam{},
		StateMutability: "pure",
		Selector:        [4]uint8{85, 203, 146, 205},
	},
	test.AbiFunction{
		Type:     "function",
		Name:     "three",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
		},
		Outputs:         []test.AbiParam{},
		StateMutability: "pure",
		Selector:        [4]uint8{156, 130, 203, 37},
	},
	test.AbiFunction{
		Type:     "function",
		Name:     "four",
		Constant: true,
		Inputs:   []test.AbiParam{},
		Outputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "string",
				AbiType: test.AbiType{
					Type: "string",
					Kind: 6,
				},
			},
		},
		StateMutability: "pure",
		Selector:        [4]uint8{161, 252, 162, 182},
	},
	test.AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type:    "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type:    "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type:    "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
		Selector: test.Word{221, 242, 82, 173, 27, 226, 200, 155, 105, 194, 176, 104, 252, 55, 141, 170, 149, 43, 167, 241, 99, 196, 161, 22, 40, 245, 90, 77, 245, 35, 179, 239},
		IndexedInputs: []test.AbiParam{
			test.AbiParam{
				Type:    "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type:    "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type:    "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
	},
}

var testBytes = []byte{
	0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15,
	0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd,
	0x5b, 0x50, 0x61, 0x03, 0x2d, 0x80, 0x61, 0x00,
	0x20, 0x60, 0x00, 0x39, 0x60, 0x00, 0xf3, 0x00,
	0x60, 0x80, 0x60, 0x40, 0x52, 0x60, 0x04, 0x36,
	0x10, 0x61, 0x00, 0x61, 0x57, 0x63, 0xff, 0xff,
	0xff, 0xff, 0x7c, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x60, 0x00, 0x35, 0x04, 0x16, 0x63, 0x26, 0x15,
	0xf1, 0x44, 0x81, 0x14, 0x61, 0x00, 0x66, 0x57,
	0x80, 0x63, 0x55, 0xcb, 0x92, 0xcd, 0x14, 0x61,
	0x01, 0x1f, 0x57, 0x80, 0x63, 0x9c, 0x82, 0xcb,
	0x25, 0x14, 0x61, 0x01, 0xba, 0x57, 0x80, 0x63,
	0xa1, 0xfc, 0xa2, 0xb6, 0x14, 0x61, 0x01, 0xe8,
	0x57, 0x5b, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x34,
	0x80, 0x15, 0x61, 0x00, 0x72, 0x57, 0x60, 0x00,
	0x80, 0xfd, 0x5b, 0x50, 0x60, 0x40, 0x80, 0x51,
	0x60, 0x20, 0x60, 0x04, 0x60, 0x24, 0x80, 0x35,
	0x82, 0x81, 0x01, 0x35, 0x84, 0x81, 0x02, 0x80,
	0x87, 0x01, 0x86, 0x01, 0x90, 0x97, 0x52, 0x80,
	0x86, 0x52, 0x61, 0x01, 0x1d, 0x96, 0x84, 0x35,
	0x96, 0x36, 0x96, 0x60, 0x44, 0x95, 0x91, 0x94,
	0x90, 0x91, 0x01, 0x92, 0x91, 0x82, 0x91, 0x85,
	0x01, 0x90, 0x84, 0x90, 0x80, 0x82, 0x84, 0x37,
	0x50, 0x50, 0x60, 0x40, 0x80, 0x51, 0x60, 0x20,
	0x60, 0x1f, 0x81, 0x8a, 0x01, 0x35, 0x8b, 0x01,
	0x80, 0x35, 0x91, 0x82, 0x01, 0x83, 0x90, 0x04,
	0x83, 0x02, 0x84, 0x01, 0x83, 0x01, 0x85, 0x52,
	0x81, 0x84, 0x52, 0x98, 0x9b, 0x75, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x19, 0x8b, 0x35, 0x16,
	0x9b, 0x90, 0x9a, 0x90, 0x99, 0x94, 0x01, 0x97,
	0x50, 0x91, 0x95, 0x50, 0x91, 0x82, 0x01, 0x93,
	0x50, 0x91, 0x50, 0x81, 0x90, 0x84, 0x01, 0x83,
	0x82, 0x80, 0x82, 0x84, 0x37, 0x50, 0x94, 0x97,
	0x50, 0x61, 0x02, 0xa8, 0x96, 0x50, 0x50, 0x50,
	0x50, 0x50, 0x50, 0x50, 0x56, 0x5b, 0x00, 0x5b,
	0x34, 0x80, 0x15, 0x61, 0x01, 0x2b, 0x57, 0x60,
	0x00, 0x80, 0xfd, 0x5b, 0x50, 0x60, 0x40, 0x80,
	0x51, 0x60, 0x20, 0x60, 0x04, 0x80, 0x35, 0x80,
	0x82, 0x01, 0x35, 0x60, 0x1f, 0x81, 0x01, 0x84,
	0x90, 0x04, 0x84, 0x02, 0x85, 0x01, 0x84, 0x01,
	0x90, 0x95, 0x52, 0x84, 0x84, 0x52, 0x61, 0x01,
	0x1d, 0x94, 0x36, 0x94, 0x92, 0x93, 0x60, 0x24,
	0x93, 0x92, 0x84, 0x01, 0x91, 0x90, 0x81, 0x90,
	0x84, 0x01, 0x83, 0x82, 0x80, 0x82, 0x84, 0x37,
	0x50, 0x50, 0x60, 0x40, 0x80, 0x51, 0x60, 0x20,
	0x80, 0x89, 0x01, 0x35, 0x8a, 0x01, 0x80, 0x35,
	0x80, 0x83, 0x02, 0x84, 0x81, 0x01, 0x84, 0x01,
	0x86, 0x52, 0x81, 0x85, 0x52, 0x99, 0x9c, 0x8b,
	0x35, 0x15, 0x15, 0x9c, 0x90, 0x9b, 0x90, 0x9a,
	0x95, 0x01, 0x98, 0x50, 0x92, 0x96, 0x50, 0x81,
	0x01, 0x94, 0x50, 0x90, 0x92, 0x50, 0x82, 0x91,
	0x90, 0x85, 0x01, 0x90, 0x84, 0x90, 0x80, 0x82,
	0x84, 0x37, 0x50, 0x94, 0x97, 0x50, 0x61, 0x02,
	0xae, 0x96, 0x50, 0x50, 0x50, 0x50, 0x50, 0x50,
	0x50, 0x56, 0x5b, 0x34, 0x80, 0x15, 0x61, 0x01,
	0xc6, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50,
	0x61, 0x01, 0x1d, 0x73, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x60, 0x04, 0x35, 0x16, 0x61, 0x02, 0xb3, 0x56,
	0x5b, 0x34, 0x80, 0x15, 0x61, 0x01, 0xf4, 0x57,
	0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x61, 0x01,
	0xfd, 0x61, 0x02, 0xb6, 0x56, 0x5b, 0x60, 0x40,
	0x51, 0x80, 0x83, 0x73, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x16, 0x73, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x16, 0x81,
	0x52, 0x60, 0x20, 0x01, 0x80, 0x60, 0x20, 0x01,
	0x82, 0x81, 0x03, 0x82, 0x52, 0x83, 0x81, 0x81,
	0x51, 0x81, 0x52, 0x60, 0x20, 0x01, 0x91, 0x50,
	0x80, 0x51, 0x90, 0x60, 0x20, 0x01, 0x90, 0x80,
	0x83, 0x83, 0x60, 0x00, 0x5b, 0x83, 0x81, 0x10,
	0x15, 0x61, 0x02, 0x6c, 0x57, 0x81, 0x81, 0x01,
	0x51, 0x83, 0x82, 0x01, 0x52, 0x60, 0x20, 0x01,
	0x61, 0x02, 0x54, 0x56, 0x5b, 0x50, 0x50, 0x50,
	0x50, 0x90, 0x50, 0x90, 0x81, 0x01, 0x90, 0x60,
	0x1f, 0x16, 0x80, 0x15, 0x61, 0x02, 0x99, 0x57,
	0x80, 0x82, 0x03, 0x80, 0x51, 0x60, 0x01, 0x83,
	0x60, 0x20, 0x03, 0x61, 0x01, 0x00, 0x0a, 0x03,
	0x19, 0x16, 0x81, 0x52, 0x60, 0x20, 0x01, 0x91,
	0x50, 0x5b, 0x50, 0x93, 0x50, 0x50, 0x50, 0x50,
	0x60, 0x40, 0x51, 0x80, 0x91, 0x03, 0x90, 0xf3,
	0x5b, 0x50, 0x50, 0x50, 0x50, 0x56, 0x5b, 0x50,
	0x50, 0x50, 0x56, 0x5b, 0x50, 0x56, 0x5b, 0x60,
	0x40, 0x80, 0x51, 0x80, 0x82, 0x01, 0x90, 0x91,
	0x52, 0x60, 0x04, 0x81, 0x52, 0x7f, 0x74, 0x65,
	0x73, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x20,
	0x82, 0x01, 0x52, 0x72, 0xa3, 0x29, 0xc0, 0x64,
	0x87, 0x69, 0xa7, 0x3a, 0xfa, 0xc7, 0xf9, 0x38,
	0x1e, 0x08, 0xfb, 0x43, 0xdb, 0xea, 0x72, 0x91,
	0x56, 0x00, 0xa1, 0x65, 0x62, 0x7a, 0x7a, 0x72,
	0x30, 0x58, 0x20, 0xde, 0xe1, 0x6d, 0xfb, 0xcb,
	0x1d, 0x75, 0x1e, 0x36, 0xfd, 0x09, 0xe7, 0x7d,
	0xe2, 0x3a, 0xfc, 0xf9, 0x38, 0xd6, 0xd3, 0xa6,
	0x74, 0x02, 0x48, 0xa9, 0x11, 0x9e, 0x4d, 0x28,
	0x53, 0x22, 0x87, 0x00, 0x29,
}

const testOutputDefault = `test.Abi{
	test.AbiFunction{
		Type: "function",
		Name: "two",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Type: "uint256",
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
			{
				Type: "uint32[]",
				AbiType: test.AbiType{
					Type: "uint32[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint32",
						Kind: 2,
					},
				},
			},
			{
				Type: "bytes10",
				AbiType: test.AbiType{
					Type: "bytes10",
					Kind: 6,
					ArrayLen: 10,
					FixedLen: true,
				},
			},
			{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44},
	},
	test.AbiFunction{
		Type: "function",
		Name: "one",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
			{
				Type: "bool",
				AbiType: test.AbiType{
					Type: "bool",
					Kind: 1,
				},
			},
			{
				Type: "uint256[]",
				AbiType: test.AbiType{
					Type: "uint256[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint256",
						Kind: 2,
					},
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd},
	},
	test.AbiFunction{
		Type: "function",
		Name: "three",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25},
	},
	test.AbiFunction{
		Type: "function",
		Name: "four",
		Constant: true,
		Inputs: []test.AbiParam{},
		Outputs: []test.AbiParam{
			{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "string",
				AbiType: test.AbiType{
					Type: "string",
					Kind: 6,
				},
			},
		},
		StateMutability: "pure",
		Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6},
	},
	test.AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []test.AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
		Selector: test.Word{
			0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
			0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
			0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
			0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
		},
		IndexedInputs: []test.AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
	},
}`

const testOutputSingleLine = `test.Abi{test.AbiFunction{Type: "function", Name: "two", Constant: true, Inputs: []test.AbiParam{{Type: "uint256", AbiType: test.AbiType{Type: "uint256", Kind: 2}}, {Type: "uint32[]", AbiType: test.AbiType{Type: "uint32[]", Kind: 7, Elem: &test.AbiType{Type: "uint32", Kind: 2}}}, {Type: "bytes10", AbiType: test.AbiType{Type: "bytes10", Kind: 6, ArrayLen: 10, FixedLen: true}}, {Type: "bytes", AbiType: test.AbiType{Type: "bytes", Kind: 6}}}, Outputs: []test.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44}}, test.AbiFunction{Type: "function", Name: "one", Constant: true, Inputs: []test.AbiParam{{Type: "bytes", AbiType: test.AbiType{Type: "bytes", Kind: 6}}, {Type: "bool", AbiType: test.AbiType{Type: "bool", Kind: 1}}, {Type: "uint256[]", AbiType: test.AbiType{Type: "uint256[]", Kind: 7, Elem: &test.AbiType{Type: "uint256", Kind: 2}}}}, Outputs: []test.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd}}, test.AbiFunction{Type: "function", Name: "three", Constant: true, Inputs: []test.AbiParam{{Type: "address", AbiType: test.AbiType{Type: "address", Kind: 4}}}, Outputs: []test.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25}}, test.AbiFunction{Type: "function", Name: "four", Constant: true, Inputs: []test.AbiParam{}, Outputs: []test.AbiParam{{Type: "address", AbiType: test.AbiType{Type: "address", Kind: 4}}, {Type: "string", AbiType: test.AbiType{Type: "string", Kind: 6}}}, StateMutability: "pure", Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6}}, test.AbiEvent{Type: "event", Name: "Transfer", Inputs: []test.AbiParam{{Type: "address", Indexed: true, AbiType: test.AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: test.AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: test.AbiType{Type: "uint256", Kind: 2}}}, Selector: test.Word{0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b, 0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa, 0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16, 0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef}, IndexedInputs: []test.AbiParam{{Type: "address", Indexed: true, AbiType: test.AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: test.AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: test.AbiType{Type: "uint256", Kind: 2}}}}}`

const testOutputWithoutPackageName = `Abi{
	AbiFunction{
		Type: "function",
		Name: "two",
		Constant: true,
		Inputs: []AbiParam{
			{
				Type: "uint256",
				AbiType: AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
			{
				Type: "uint32[]",
				AbiType: AbiType{
					Type: "uint32[]",
					Kind: 7,
					Elem: &AbiType{
						Type: "uint32",
						Kind: 2,
					},
				},
			},
			{
				Type: "bytes10",
				AbiType: AbiType{
					Type: "bytes10",
					Kind: 6,
					ArrayLen: 10,
					FixedLen: true,
				},
			},
			{
				Type: "bytes",
				AbiType: AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
		},
		Outputs: []AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44},
	},
	AbiFunction{
		Type: "function",
		Name: "one",
		Constant: true,
		Inputs: []AbiParam{
			{
				Type: "bytes",
				AbiType: AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
			{
				Type: "bool",
				AbiType: AbiType{
					Type: "bool",
					Kind: 1,
				},
			},
			{
				Type: "uint256[]",
				AbiType: AbiType{
					Type: "uint256[]",
					Kind: 7,
					Elem: &AbiType{
						Type: "uint256",
						Kind: 2,
					},
				},
			},
		},
		Outputs: []AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd},
	},
	AbiFunction{
		Type: "function",
		Name: "three",
		Constant: true,
		Inputs: []AbiParam{
			{
				Type: "address",
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
		},
		Outputs: []AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25},
	},
	AbiFunction{
		Type: "function",
		Name: "four",
		Constant: true,
		Inputs: []AbiParam{},
		Outputs: []AbiParam{
			{
				Type: "address",
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "string",
				AbiType: AbiType{
					Type: "string",
					Kind: 6,
				},
			},
		},
		StateMutability: "pure",
		Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6},
	},
	AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
		Selector: Word{
			0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
			0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
			0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
			0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
		},
		IndexedInputs: []AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
	},
}`

const testOutputSingleLineWithoutPackageName = `Abi{AbiFunction{Type: "function", Name: "two", Constant: true, Inputs: []AbiParam{{Type: "uint256", AbiType: AbiType{Type: "uint256", Kind: 2}}, {Type: "uint32[]", AbiType: AbiType{Type: "uint32[]", Kind: 7, Elem: &AbiType{Type: "uint32", Kind: 2}}}, {Type: "bytes10", AbiType: AbiType{Type: "bytes10", Kind: 6, ArrayLen: 10, FixedLen: true}}, {Type: "bytes", AbiType: AbiType{Type: "bytes", Kind: 6}}}, Outputs: []AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44}}, AbiFunction{Type: "function", Name: "one", Constant: true, Inputs: []AbiParam{{Type: "bytes", AbiType: AbiType{Type: "bytes", Kind: 6}}, {Type: "bool", AbiType: AbiType{Type: "bool", Kind: 1}}, {Type: "uint256[]", AbiType: AbiType{Type: "uint256[]", Kind: 7, Elem: &AbiType{Type: "uint256", Kind: 2}}}}, Outputs: []AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd}}, AbiFunction{Type: "function", Name: "three", Constant: true, Inputs: []AbiParam{{Type: "address", AbiType: AbiType{Type: "address", Kind: 4}}}, Outputs: []AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25}}, AbiFunction{Type: "function", Name: "four", Constant: true, Inputs: []AbiParam{}, Outputs: []AbiParam{{Type: "address", AbiType: AbiType{Type: "address", Kind: 4}}, {Type: "string", AbiType: AbiType{Type: "string", Kind: 6}}}, StateMutability: "pure", Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6}}, AbiEvent{Type: "event", Name: "Transfer", Inputs: []AbiParam{{Type: "address", Indexed: true, AbiType: AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: AbiType{Type: "uint256", Kind: 2}}}, Selector: Word{0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b, 0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa, 0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16, 0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef}, IndexedInputs: []AbiParam{{Type: "address", Indexed: true, AbiType: AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: AbiType{Type: "uint256", Kind: 2}}}}}`

const testOutputRenamed = `renamed.Abi{
	renamed.AbiFunction{
		Type: "function",
		Name: "two",
		Constant: true,
		Inputs: []renamed.AbiParam{
			{
				Type: "uint256",
				AbiType: renamed.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
			{
				Type: "uint32[]",
				AbiType: renamed.AbiType{
					Type: "uint32[]",
					Kind: 7,
					Elem: &renamed.AbiType{
						Type: "uint32",
						Kind: 2,
					},
				},
			},
			{
				Type: "bytes10",
				AbiType: renamed.AbiType{
					Type: "bytes10",
					Kind: 6,
					ArrayLen: 10,
					FixedLen: true,
				},
			},
			{
				Type: "bytes",
				AbiType: renamed.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
		},
		Outputs: []renamed.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44},
	},
	renamed.AbiFunction{
		Type: "function",
		Name: "one",
		Constant: true,
		Inputs: []renamed.AbiParam{
			{
				Type: "bytes",
				AbiType: renamed.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
			{
				Type: "bool",
				AbiType: renamed.AbiType{
					Type: "bool",
					Kind: 1,
				},
			},
			{
				Type: "uint256[]",
				AbiType: renamed.AbiType{
					Type: "uint256[]",
					Kind: 7,
					Elem: &renamed.AbiType{
						Type: "uint256",
						Kind: 2,
					},
				},
			},
		},
		Outputs: []renamed.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd},
	},
	renamed.AbiFunction{
		Type: "function",
		Name: "three",
		Constant: true,
		Inputs: []renamed.AbiParam{
			{
				Type: "address",
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
		},
		Outputs: []renamed.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25},
	},
	renamed.AbiFunction{
		Type: "function",
		Name: "four",
		Constant: true,
		Inputs: []renamed.AbiParam{},
		Outputs: []renamed.AbiParam{
			{
				Type: "address",
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "string",
				AbiType: renamed.AbiType{
					Type: "string",
					Kind: 6,
				},
			},
		},
		StateMutability: "pure",
		Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6},
	},
	renamed.AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []renamed.AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
		Selector: renamed.Word{
			0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
			0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
			0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
			0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
		},
		IndexedInputs: []renamed.AbiParam{
			{
				Type: "address",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "address",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			{
				Type: "uint256",
				Indexed: true,
				AbiType: renamed.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
	},
}`

const testOutputSingleLineRenamed = `renamed.Abi{renamed.AbiFunction{Type: "function", Name: "two", Constant: true, Inputs: []renamed.AbiParam{{Type: "uint256", AbiType: renamed.AbiType{Type: "uint256", Kind: 2}}, {Type: "uint32[]", AbiType: renamed.AbiType{Type: "uint32[]", Kind: 7, Elem: &renamed.AbiType{Type: "uint32", Kind: 2}}}, {Type: "bytes10", AbiType: renamed.AbiType{Type: "bytes10", Kind: 6, ArrayLen: 10, FixedLen: true}}, {Type: "bytes", AbiType: renamed.AbiType{Type: "bytes", Kind: 6}}}, Outputs: []renamed.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44}}, renamed.AbiFunction{Type: "function", Name: "one", Constant: true, Inputs: []renamed.AbiParam{{Type: "bytes", AbiType: renamed.AbiType{Type: "bytes", Kind: 6}}, {Type: "bool", AbiType: renamed.AbiType{Type: "bool", Kind: 1}}, {Type: "uint256[]", AbiType: renamed.AbiType{Type: "uint256[]", Kind: 7, Elem: &renamed.AbiType{Type: "uint256", Kind: 2}}}}, Outputs: []renamed.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd}}, renamed.AbiFunction{Type: "function", Name: "three", Constant: true, Inputs: []renamed.AbiParam{{Type: "address", AbiType: renamed.AbiType{Type: "address", Kind: 4}}}, Outputs: []renamed.AbiParam{}, StateMutability: "pure", Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25}}, renamed.AbiFunction{Type: "function", Name: "four", Constant: true, Inputs: []renamed.AbiParam{}, Outputs: []renamed.AbiParam{{Type: "address", AbiType: renamed.AbiType{Type: "address", Kind: 4}}, {Type: "string", AbiType: renamed.AbiType{Type: "string", Kind: 6}}}, StateMutability: "pure", Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6}}, renamed.AbiEvent{Type: "event", Name: "Transfer", Inputs: []renamed.AbiParam{{Type: "address", Indexed: true, AbiType: renamed.AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: renamed.AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: renamed.AbiType{Type: "uint256", Kind: 2}}}, Selector: renamed.Word{0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b, 0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa, 0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16, 0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef}, IndexedInputs: []renamed.AbiParam{{Type: "address", Indexed: true, AbiType: renamed.AbiType{Type: "address", Kind: 4}}, {Type: "address", Indexed: true, AbiType: renamed.AbiType{Type: "address", Kind: 4}}, {Type: "uint256", Indexed: true, AbiType: renamed.AbiType{Type: "uint256", Kind: 2}}}}}`

const testOutputWithZeroFields = `test.Abi{
	test.AbiFunction{
		Type: "function",
		Name: "two",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Name: "",
				Type: "uint256",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "uint32[]",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "uint32[]",
					Kind: 7,
					ArrayLen: 0,
					FixedLen: false,
					Elem: &test.AbiType{
						Type: "uint32",
						Kind: 2,
						ArrayLen: 0,
						FixedLen: false,
						Elem: nil,
					},
				},
			},
			{
				Name: "",
				Type: "bytes10",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "bytes10",
					Kind: 6,
					ArrayLen: 10,
					FixedLen: true,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "bytes",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
		},
		Outputs: []test.AbiParam{},
		Payable: false,
		StateMutability: "pure",
		Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44},
	},
	test.AbiFunction{
		Type: "function",
		Name: "one",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Name: "",
				Type: "bytes",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "bool",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "bool",
					Kind: 1,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "uint256[]",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "uint256[]",
					Kind: 7,
					ArrayLen: 0,
					FixedLen: false,
					Elem: &test.AbiType{
						Type: "uint256",
						Kind: 2,
						ArrayLen: 0,
						FixedLen: false,
						Elem: nil,
					},
				},
			},
		},
		Outputs: []test.AbiParam{},
		Payable: false,
		StateMutability: "pure",
		Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd},
	},
	test.AbiFunction{
		Type: "function",
		Name: "three",
		Constant: true,
		Inputs: []test.AbiParam{
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
		},
		Outputs: []test.AbiParam{},
		Payable: false,
		StateMutability: "pure",
		Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25},
	},
	test.AbiFunction{
		Type: "function",
		Name: "four",
		Constant: true,
		Inputs: []test.AbiParam{},
		Outputs: []test.AbiParam{
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "string",
				Components: nil,
				Indexed: false,
				AbiType: test.AbiType{
					Type: "string",
					Kind: 6,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
		},
		Payable: false,
		StateMutability: "pure",
		Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6},
	},
	test.AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []test.AbiParam{
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "uint256",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
		},
		Anonymous: false,
		Selector: test.Word{
			0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
			0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
			0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
			0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
		},
		IndexedInputs: []test.AbiParam{
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "address",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
			{
				Name: "",
				Type: "uint256",
				Components: nil,
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
					ArrayLen: 0,
					FixedLen: false,
					Elem: nil,
				},
			},
		},
		NonIndexedInputs: nil,
	},
}`

const testOutputWithForcedConstructorNames = `test.Abi{
	test.AbiFunction{
		Type: "function",
		Name: "two",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "uint256",
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
			test.AbiParam{
				Type: "uint32[]",
				AbiType: test.AbiType{
					Type: "uint32[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint32",
						Kind: 2,
					},
				},
			},
			test.AbiParam{
				Type: "bytes10",
				AbiType: test.AbiType{
					Type: "bytes10",
					Kind: 6,
					ArrayLen: 10,
					FixedLen: true,
				},
			},
			test.AbiParam{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x26, 0x15, 0xf1, 0x44},
	},
	test.AbiFunction{
		Type: "function",
		Name: "one",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "bytes",
				AbiType: test.AbiType{
					Type: "bytes",
					Kind: 6,
				},
			},
			test.AbiParam{
				Type: "bool",
				AbiType: test.AbiType{
					Type: "bool",
					Kind: 1,
				},
			},
			test.AbiParam{
				Type: "uint256[]",
				AbiType: test.AbiType{
					Type: "uint256[]",
					Kind: 7,
					Elem: &test.AbiType{
						Type: "uint256",
						Kind: 2,
					},
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x55, 0xcb, 0x92, 0xcd},
	},
	test.AbiFunction{
		Type: "function",
		Name: "three",
		Constant: true,
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
		},
		Outputs: []test.AbiParam{},
		StateMutability: "pure",
		Selector: [4]uint8{0x9c, 0x82, 0xcb, 0x25},
	},
	test.AbiFunction{
		Type: "function",
		Name: "four",
		Constant: true,
		Inputs: []test.AbiParam{},
		Outputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "string",
				AbiType: test.AbiType{
					Type: "string",
					Kind: 6,
				},
			},
		},
		StateMutability: "pure",
		Selector: [4]uint8{0xa1, 0xfc, 0xa2, 0xb6},
	},
	test.AbiEvent{
		Type: "event",
		Name: "Transfer",
		Inputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
		Selector: test.Word{
			0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
			0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
			0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
			0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
		},
		IndexedInputs: []test.AbiParam{
			test.AbiParam{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "address",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "address",
					Kind: 4,
				},
			},
			test.AbiParam{
				Type: "uint256",
				Indexed: true,
				AbiType: test.AbiType{
					Type: "uint256",
					Kind: 2,
				},
			},
		},
	},
}`

const testOutputBytesHex = `[]uint8{
	0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15,
	0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd,
	0x5b, 0x50, 0x61, 0x03, 0x2d, 0x80, 0x61, 0x00,
	0x20, 0x60, 0x00, 0x39, 0x60, 0x00, 0xf3, 0x00,
	0x60, 0x80, 0x60, 0x40, 0x52, 0x60, 0x04, 0x36,
	0x10, 0x61, 0x00, 0x61, 0x57, 0x63, 0xff, 0xff,
	0xff, 0xff, 0x7c, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x60, 0x00, 0x35, 0x04, 0x16, 0x63, 0x26, 0x15,
	0xf1, 0x44, 0x81, 0x14, 0x61, 0x00, 0x66, 0x57,
	0x80, 0x63, 0x55, 0xcb, 0x92, 0xcd, 0x14, 0x61,
	0x01, 0x1f, 0x57, 0x80, 0x63, 0x9c, 0x82, 0xcb,
	0x25, 0x14, 0x61, 0x01, 0xba, 0x57, 0x80, 0x63,
	0xa1, 0xfc, 0xa2, 0xb6, 0x14, 0x61, 0x01, 0xe8,
	0x57, 0x5b, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x34,
	0x80, 0x15, 0x61, 0x00, 0x72, 0x57, 0x60, 0x00,
	0x80, 0xfd, 0x5b, 0x50, 0x60, 0x40, 0x80, 0x51,
	0x60, 0x20, 0x60, 0x04, 0x60, 0x24, 0x80, 0x35,
	0x82, 0x81, 0x01, 0x35, 0x84, 0x81, 0x02, 0x80,
	0x87, 0x01, 0x86, 0x01, 0x90, 0x97, 0x52, 0x80,
	0x86, 0x52, 0x61, 0x01, 0x1d, 0x96, 0x84, 0x35,
	0x96, 0x36, 0x96, 0x60, 0x44, 0x95, 0x91, 0x94,
	0x90, 0x91, 0x01, 0x92, 0x91, 0x82, 0x91, 0x85,
	0x01, 0x90, 0x84, 0x90, 0x80, 0x82, 0x84, 0x37,
	0x50, 0x50, 0x60, 0x40, 0x80, 0x51, 0x60, 0x20,
	0x60, 0x1f, 0x81, 0x8a, 0x01, 0x35, 0x8b, 0x01,
	0x80, 0x35, 0x91, 0x82, 0x01, 0x83, 0x90, 0x04,
	0x83, 0x02, 0x84, 0x01, 0x83, 0x01, 0x85, 0x52,
	0x81, 0x84, 0x52, 0x98, 0x9b, 0x75, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x19, 0x8b, 0x35, 0x16,
	0x9b, 0x90, 0x9a, 0x90, 0x99, 0x94, 0x01, 0x97,
	0x50, 0x91, 0x95, 0x50, 0x91, 0x82, 0x01, 0x93,
	0x50, 0x91, 0x50, 0x81, 0x90, 0x84, 0x01, 0x83,
	0x82, 0x80, 0x82, 0x84, 0x37, 0x50, 0x94, 0x97,
	0x50, 0x61, 0x02, 0xa8, 0x96, 0x50, 0x50, 0x50,
	0x50, 0x50, 0x50, 0x50, 0x56, 0x5b, 0x00, 0x5b,
	0x34, 0x80, 0x15, 0x61, 0x01, 0x2b, 0x57, 0x60,
	0x00, 0x80, 0xfd, 0x5b, 0x50, 0x60, 0x40, 0x80,
	0x51, 0x60, 0x20, 0x60, 0x04, 0x80, 0x35, 0x80,
	0x82, 0x01, 0x35, 0x60, 0x1f, 0x81, 0x01, 0x84,
	0x90, 0x04, 0x84, 0x02, 0x85, 0x01, 0x84, 0x01,
	0x90, 0x95, 0x52, 0x84, 0x84, 0x52, 0x61, 0x01,
	0x1d, 0x94, 0x36, 0x94, 0x92, 0x93, 0x60, 0x24,
	0x93, 0x92, 0x84, 0x01, 0x91, 0x90, 0x81, 0x90,
	0x84, 0x01, 0x83, 0x82, 0x80, 0x82, 0x84, 0x37,
	0x50, 0x50, 0x60, 0x40, 0x80, 0x51, 0x60, 0x20,
	0x80, 0x89, 0x01, 0x35, 0x8a, 0x01, 0x80, 0x35,
	0x80, 0x83, 0x02, 0x84, 0x81, 0x01, 0x84, 0x01,
	0x86, 0x52, 0x81, 0x85, 0x52, 0x99, 0x9c, 0x8b,
	0x35, 0x15, 0x15, 0x9c, 0x90, 0x9b, 0x90, 0x9a,
	0x95, 0x01, 0x98, 0x50, 0x92, 0x96, 0x50, 0x81,
	0x01, 0x94, 0x50, 0x90, 0x92, 0x50, 0x82, 0x91,
	0x90, 0x85, 0x01, 0x90, 0x84, 0x90, 0x80, 0x82,
	0x84, 0x37, 0x50, 0x94, 0x97, 0x50, 0x61, 0x02,
	0xae, 0x96, 0x50, 0x50, 0x50, 0x50, 0x50, 0x50,
	0x50, 0x56, 0x5b, 0x34, 0x80, 0x15, 0x61, 0x01,
	0xc6, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50,
	0x61, 0x01, 0x1d, 0x73, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x60, 0x04, 0x35, 0x16, 0x61, 0x02, 0xb3, 0x56,
	0x5b, 0x34, 0x80, 0x15, 0x61, 0x01, 0xf4, 0x57,
	0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x61, 0x01,
	0xfd, 0x61, 0x02, 0xb6, 0x56, 0x5b, 0x60, 0x40,
	0x51, 0x80, 0x83, 0x73, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x16, 0x73, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x16, 0x81,
	0x52, 0x60, 0x20, 0x01, 0x80, 0x60, 0x20, 0x01,
	0x82, 0x81, 0x03, 0x82, 0x52, 0x83, 0x81, 0x81,
	0x51, 0x81, 0x52, 0x60, 0x20, 0x01, 0x91, 0x50,
	0x80, 0x51, 0x90, 0x60, 0x20, 0x01, 0x90, 0x80,
	0x83, 0x83, 0x60, 0x00, 0x5b, 0x83, 0x81, 0x10,
	0x15, 0x61, 0x02, 0x6c, 0x57, 0x81, 0x81, 0x01,
	0x51, 0x83, 0x82, 0x01, 0x52, 0x60, 0x20, 0x01,
	0x61, 0x02, 0x54, 0x56, 0x5b, 0x50, 0x50, 0x50,
	0x50, 0x90, 0x50, 0x90, 0x81, 0x01, 0x90, 0x60,
	0x1f, 0x16, 0x80, 0x15, 0x61, 0x02, 0x99, 0x57,
	0x80, 0x82, 0x03, 0x80, 0x51, 0x60, 0x01, 0x83,
	0x60, 0x20, 0x03, 0x61, 0x01, 0x00, 0x0a, 0x03,
	0x19, 0x16, 0x81, 0x52, 0x60, 0x20, 0x01, 0x91,
	0x50, 0x5b, 0x50, 0x93, 0x50, 0x50, 0x50, 0x50,
	0x60, 0x40, 0x51, 0x80, 0x91, 0x03, 0x90, 0xf3,
	0x5b, 0x50, 0x50, 0x50, 0x50, 0x56, 0x5b, 0x50,
	0x50, 0x50, 0x56, 0x5b, 0x50, 0x56, 0x5b, 0x60,
	0x40, 0x80, 0x51, 0x80, 0x82, 0x01, 0x90, 0x91,
	0x52, 0x60, 0x04, 0x81, 0x52, 0x7f, 0x74, 0x65,
	0x73, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x20,
	0x82, 0x01, 0x52, 0x72, 0xa3, 0x29, 0xc0, 0x64,
	0x87, 0x69, 0xa7, 0x3a, 0xfa, 0xc7, 0xf9, 0x38,
	0x1e, 0x08, 0xfb, 0x43, 0xdb, 0xea, 0x72, 0x91,
	0x56, 0x00, 0xa1, 0x65, 0x62, 0x7a, 0x7a, 0x72,
	0x30, 0x58, 0x20, 0xde, 0xe1, 0x6d, 0xfb, 0xcb,
	0x1d, 0x75, 0x1e, 0x36, 0xfd, 0x09, 0xe7, 0x7d,
	0xe2, 0x3a, 0xfc, 0xf9, 0x38, 0xd6, 0xd3, 0xa6,
	0x74, 0x02, 0x48, 0xa9, 0x11, 0x9e, 0x4d, 0x28,
	0x53, 0x22, 0x87, 0x00, 0x29,
}`
