/*
Overview

Prints Go data structures as syntactically valid Go code. Useful for code
generation. The name "repr" stands for "representation" and alludes to the
Python function with the same name.

Solves a problem unaddressed by https://github.com/davecgh/go-spew/spew and
directly competes with https://github.com/shurcooL/go-goon

Has no dependencies outside the standard library.

Why

Motives:

• Dumping data as code is useful for code generation and debugging.

• fmt.Sprintf("%#v") doesn't always output valid code, has no multiline mode.

• https://github.com/davecgh/go-spew/spew doesn't output valid code, has no single-line mode.

• https://github.com/shurcooL/go-goon outputs too much noise, has no single-line mode.

Features

Supports single-line and multiline modes. Defaults to multiline.

The output looks like something you'd write by hand, and is almost exactly
compliant with gofmt. Unlike gofmt, it doesn't align field values in struct
literals. Use "go/format" to fix that, at a 50x performance cost:

	import (
		"go/format"
		"github.com/mitranim/repr"
	)

	code := repr.Bytes(someDataStructure)
	code, err := format.Source(code)

Zero-initialized fields in structs are omitted by default (configurable).

Bytes are printed in hex notation. In multiline mode, byte arrays have 8 bytes
per row:

	var output = []uint8{
		0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15,
		0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd,
		0x5b, 0x50, 0x61, 0x03, 0x2d, 0x80, 0x61, 0x00,
		0x20, 0x60, 0x00, 0x39, 0x60, 0x00, 0xf3, 0x00,
		// ...
	}

Supports package renaming, which is useful for code generation. See Config for
details.

Limitations

Some of these limitations may be lifted in future versions.

• Fancy types such as "big.Int" or "time.Time" are printed as empty structs;
ideally they would be printed as constructor calls.

• Funcs are treated as nil.

• Chans are treated as nil.

• Pointers to primitive types are not supported and cause a panic.

• "byte" is printed as "uint8".

• "rune" is printed as "int32".

• Runes are printed as integers, not character literals.

• Enum-style constants are not mapped back to identifers.

• On structs, only exported fields are included.

• Cyclic structures cause infinite recursion.

Note: pointers to composite types such as structs, arrays, slices and maps are
supported by prefixing literals with "&", but Go currently doesn't support this
for primitive literals.

Installation

Shell:

	go get -u github.com/mitranim/repr

Usage

Example:

	import "github.com/mitranim/repr"

	type Data struct {
		Number int
		String string
		List   []int
	}

	repr.Println(Data{
		Number: 123,
		String: "hello world!",
		List:   []int{10, 20, 30},
	})

	// Output
	Data{
		Number: 123,
		String: "hello world!",
		List: []int{10, 20, 30},
	}

Misc

I'm receptive to suggestions. If this package almost satisfies you but needs
changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
*/
package repr

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

/*
Format settings. Can be passed to "StringC", "BytesC" and "AppendC". A
zero-initialized config is almost the same as the "Default" config.
*/
type Config struct {
	/**
	If true, put everything on the same line. If false (default), produce
	multiline output, using tabs for indentation.
	*/
	SingleLine bool

	/**
	If true, include zero fields in struct literals. If false (default), omit
	zero fields from struct literals.

	"Zero value" for any type is defined as:

		bool    = false
		number  = 0
		string  = ""
		nilable = nil
		array   = every byte is 0
		struct  = every byte is 0
	*/
	ZeroFields bool

	/**
	If true, always print constructor names for elements in arrays and slices. If
	false (default), elide them wherever possible.
	*/
	ForceConstructorName bool

	/**
	Maps fully-qualified packages to short aliases. Useful for code generation.
	An empty string causes the package name to be stripped. The default config
	strips "main." from the output:

		map[string]string{"main": ""}

	For external packages, keys should be fully-qualified:

		map[string]string{
			"golang.org/x/sys": "sys",
		}
	*/
	PackageMap map[string]string
}

/*
Global/default settings. Used by functions like "String". Custom configs can be
passed to functions like "StringC".
*/
var Default = Config{
	PackageMap: map[string]string{
		"main": "",
	},
}

/*
Formats the value using the "Default" config. See "Config" for details.
*/
func String(val interface{}) string {
	return bytesToMutableString(appendAny(nil, val, state{conf: Default}))
}

/*
Short for "String with config". Formats the value using the provided config. See
"Config" for details.
*/
func StringC(val interface{}, conf Config) string {
	return bytesToMutableString(appendAny(nil, val, state{conf: conf}))
}

/*
Formats the value using the "Default" config. See "Config" for details.
*/
func Bytes(val interface{}) []byte {
	return appendAny(nil, val, state{conf: Default})
}

/*
"Short for "Bytes with config". Formats the value using the provided config. See
"Config" for details.
*/
func BytesC(val interface{}, conf Config) []byte {
	return appendAny(nil, val, state{conf: conf})
}

/*
Formats the value using the "Default" config, appending the output to the
provided buffer. See "Config" for details.
*/
func Append(out []byte, val interface{}) []byte {
	return appendAny(nil, val, state{conf: Default})
}

/*
Short for "Append with config". Formats the value using the provided config,
appending the output to the provided buffer. See "Config" for details.
*/
func AppendC(out []byte, val interface{}, conf Config) []byte {
	return appendAny(out, val, state{conf: conf})
}

/*
Shortcut for `fmt.Println(repr.String(val))`.
*/
func Println(val interface{}) (int, error) {
	return fmt.Println(String(val))
}

/*
Shortcut for `fmt.Println(repr.StringC(val, conf))`.
*/
func PrintlnC(val interface{}, conf Config) (int, error) {
	return fmt.Println(StringC(val, conf))
}

var (
	byteType      = reflect.TypeOf(byte(0))
	byteSliceType = reflect.TypeOf([]byte(nil))
)

type state struct {
	conf      Config
	indent    int
	elideType bool
}

func appendAny(out []byte, val interface{}, state state) []byte {
	// Well-known types
	switch val := val.(type) {
	case bool:
		if val {
			return append(out, "true"...)
		}
		return append(out, "false"...)
	case uint8: // = byte
		return appendByteHex(out, val)
	case uint16:
		return strconv.AppendUint(out, uint64(val), 10)
	case uint32:
		return strconv.AppendUint(out, uint64(val), 10)
	case uint64:
		return strconv.AppendUint(out, uint64(val), 10)
	case uint:
		return strconv.AppendUint(out, uint64(val), 10)
	case uintptr:
		return strconv.AppendUint(append(out, '0', 'x'), uint64(val), 16)
	case unsafe.Pointer:
		return strconv.AppendUint(append(out, '0', 'x'), uint64(uintptr(val)), 16)
	case int8:
		return strconv.AppendInt(out, int64(val), 10)
	case int16:
		return strconv.AppendInt(out, int64(val), 10)
	case int32: // = rune
		return strconv.AppendInt(out, int64(val), 10)
	case int64:
		return strconv.AppendInt(out, int64(val), 10)
	case int:
		return strconv.AppendInt(out, int64(val), 10)
	case float32:
		return strconv.AppendFloat(out, float64(val), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(out, float64(val), 'f', -1, 64)
	case complex64:
		return appendComplex128(out, complex128(val))
	case complex128:
		return appendComplex128(out, val)
	case string:
		return strconv.AppendQuote(out, val)
	case []byte:
		if !state.elideType {
			out = append(out, "[]uint8"...)
		}
		out = appendBytes(out, val, state)
		return out
	}

	rval := reflect.ValueOf(val)
	if !rval.IsValid() {
		out = append(out, "nil"...)
		return out
	}

	rtype := rval.Type()

	switch rtype.Kind() {
	case reflect.Bool:
		out = appendCastPrefix(out, rval, state)
		if rval.Bool() {
			out = append(out, "true"...)
		} else {
			out = append(out, "false"...)
		}
		out = appendCastSuffix(out, rval, state)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendInt(out, rval.Int(), 10)
		out = appendCastSuffix(out, rval, state)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendUint(out, rval.Uint(), 10)
		out = appendCastSuffix(out, rval, state)

	case reflect.Uintptr:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendUint(append(out, '0', 'x'), rval.Uint(), 16)
		out = appendCastSuffix(out, rval, state)

	case reflect.Float32:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendFloat(out, rval.Float(), 'f', -1, 32)
		out = appendCastSuffix(out, rval, state)

	case reflect.Float64:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendFloat(out, rval.Float(), 'f', -1, 64)
		out = appendCastSuffix(out, rval, state)

	case reflect.Complex64, reflect.Complex128:
		out = appendCastPrefix(out, rval, state)
		out = appendComplex128(out, rval.Convert(reflect.TypeOf(complex128(0))).Complex())
		out = appendCastSuffix(out, rval, state)

	case reflect.String:
		out = appendCastPrefix(out, rval, state)
		out = strconv.AppendQuote(out, rval.String())
		out = appendCastSuffix(out, rval, state)

	case reflect.Chan:
		out = appendCastPrefix(out, rval, state)
		out = append(out, "nil"...)
		out = appendCastSuffix(out, rval, state)

	case reflect.Func:
		out = appendCastPrefix(out, rval, state)
		out = append(out, "nil"...)
		out = appendCastSuffix(out, rval, state)

	// Pretty sure this should never match
	case reflect.Interface:
		panic("repr currently doesn't support printing an interface")

	case reflect.UnsafePointer:
		out = appendCastPrefix(out, rval, state)
		ptr := rval.Convert(reflect.TypeOf(unsafe.Pointer(nil))).Interface().(unsafe.Pointer)
		out = strconv.AppendUint(append(out, '0', 'x'), uint64(uintptr(ptr)), 16)
		out = appendCastSuffix(out, rval, state)

	case reflect.Ptr:
		switch rtype.Elem().Kind() {
		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			if isZeroOrShouldOmit(rval) {
				out = append(out, "nil"...)
			} else {
				out = append(out, '&')
				out = appendAny(out, rval.Elem().Interface(), state)
			}
		default:
			panic("repr currently doesn't support pointers to non-composite types")
		}

	case reflect.Array:
		if !state.elideType {
			out = appendTypeName(out, rval.Type(), state)
		}
		if rtype.Elem() == byteType {
			out = appendBytes(out, byteArrayToSlice(rval), state)
		} else {
			out = appendList(out, rval, state)
		}

	case reflect.Slice:
		if rval.IsNil() {
			if state.elideType {
				out = append(out, "nil"...)
			} else {
				out = appendTypeName(out, rval.Type(), state)
				out = append(out, "(nil)"...)
			}
		} else {
			out = appendTypeName(out, rval.Type(), state)
			if rtype.Elem() == byteType {
				out = appendBytes(out, rval.Bytes(), state)
			} else {
				out = appendList(out, rval, state)
			}
		}

	case reflect.Struct:
		if !state.elideType {
			out = appendTypeName(out, rval.Type(), state)
		}
		out = appendStruct(out, rval, state)

	case reflect.Map:
		if rval.IsNil() {
			if state.elideType {
				out = append(out, "nil"...)
			} else {
				out = appendTypeName(out, rval.Type(), state)
				out = append(out, "(nil)"...)
			}
		} else {
			out = appendTypeName(out, rval.Type(), state)
			out = appendMap(out, rval, state)
		}
	}

	return out
}

func appendComplex128(out []byte, val complex128) []byte {
	out = append(out, '(')
	out = strconv.AppendFloat(out, real(val), 'f', -1, 64)
	i := imag(val)
	if !(i < 0) {
		out = append(out, '+')
	}
	out = strconv.AppendFloat(out, i, 'f', -1, 64)
	out = append(out, 'i', ')')
	return out
}

func appendList(out []byte, rval reflect.Value, state state) []byte {
	elemType := rval.Type().Elem()
	state.elideType = canElideType(elemType, state)
	count := rval.Len()

	if state.conf.SingleLine || (!mayRequireMultiline(elemType) && count < 48) {
		state.indent = 0
		out = append(out, '{')
		for i := 0; i < count; i++ {
			out = appendAny(out, rval.Index(i).Interface(), state)
			if i < count-1 {
				out = append(out, ',', ' ')
			}
		}
		out = append(out, '}')
		return out
	}

	out = append(out, '{')
	if count > 0 {
		out = append(out, '\n')
		state.indent++
	}

	for i := 0; i < count; i++ {
		out = appendIndent(out, state)
		out = appendAny(out, rval.Index(i).Interface(), state)
		out = append(out, ',', '\n')
	}

	if count > 0 {
		state.indent--
		out = appendIndent(out, state)
	}

	out = append(out, '}')
	return out
}

func appendStruct(out []byte, rval reflect.Value, state state) []byte {
	rtype := rval.Type()

	if state.conf.SingleLine {
		state.indent = 0
		var hasFields bool

		out = append(out, '{')
		for i := 0; i < rtype.NumField(); i++ {
			sfield := rtype.Field(i)
			if !isSfieldExported(sfield) {
				continue
			}

			rfield := rval.Field(i)
			if !state.conf.ZeroFields && isZeroOrShouldOmit(rfield) {
				continue
			}

			if hasFields {
				out = append(out, ',', ' ')
			}
			hasFields = true

			out = append(out, sfield.Name...)
			out = append(out, ':', ' ')

			state := state
			state.elideType = isPrimitive(rfield.Type()) || isNil(rfield)
			out = appendAny(out, rfield.Interface(), state)
		}
		out = append(out, '}')
		return out
	}

	count := 0
	out = append(out, '{')

	for i := 0; i < rtype.NumField(); i++ {
		sfield := rtype.Field(i)
		if !isSfieldExported(sfield) {
			continue
		}

		rfield := rval.Field(i)
		if !state.conf.ZeroFields && isZeroOrShouldOmit(rfield) {
			continue
		}

		count++
		if count == 1 {
			out = append(out, '\n')
			state.indent++
		}

		out = appendIndent(out, state)
		out = append(out, sfield.Name...)
		out = append(out, ':', ' ')

		state := state
		state.elideType = isPrimitive(rfield.Type()) || isNil(rfield)
		out = appendAny(out, rfield.Interface(), state)
		out = append(out, ',', '\n')
	}

	if count > 0 {
		state.indent--
		out = appendIndent(out, state)
	}

	out = append(out, '}')
	return out
}

// TODO: the test doesn't cover constructor elision in maps.
func appendMap(out []byte, rval reflect.Value, state state) []byte {
	rtype := rval.Type()
	keyType := rtype.Key()
	elemType := rtype.Elem()
	elideKeyType := canElideType(keyType, state)
	elideElemType := canElideType(elemType, state)

	if state.conf.SingleLine {
		state.indent = 0

		keyState := state
		keyState.elideType = elideKeyType

		elemState := state
		elemState.elideType = elideElemType

		keys := rval.MapKeys()

		out = append(out, '{')
		for i, key := range keys {
			out = appendAny(out, key.Interface(), keyState)
			out = append(out, ':', ' ')
			out = appendAny(out, rval.MapIndex(key).Interface(), elemState)
			if i < len(keys)-1 {
				out = append(out, ',', ' ')
			}
		}
		out = append(out, '}')
		return out
	}

	out = append(out, '{')
	keys := rval.MapKeys()

	for i, key := range keys {
		if i == 0 {
			out = append(out, '\n')
			state.indent++
		}

		keyState := state
		keyState.elideType = elideKeyType

		elemState := state
		elemState.elideType = elideElemType

		out = appendIndent(out, state)
		out = appendAny(out, key.Interface(), keyState)
		out = append(out, ':', ' ')
		out = appendAny(out, rval.MapIndex(key).Interface(), elemState)

		out = append(out, ',', '\n')
	}

	if len(keys) > 0 {
		state.indent--
		out = appendIndent(out, state)
	}

	out = append(out, '}')
	return out
}

// Similar to fmt.Sprintf("%#02v", val), but multiline: large inputs are printed
// as a column with 8 bytes per row.
func appendBytes(out []byte, val []byte, state state) []byte {
	if state.conf.SingleLine || len(val) <= 8 {
		out = append(out, '{')

		for i, char := range val {
			out = appendByteHex(out, char)
			if i < len(val)-1 {
				out = append(out, ',', ' ')
			}
		}

		out = append(out, '}')
		return out
	}

	state.indent++
	out = append(out, '{', '\n')

	for i, char := range val {
		if i == 0 {
			out = appendIndent(out, state)
		} else if i%8 == 0 {
			out = append(out, ',', '\n')
			out = appendIndent(out, state)
		} else {
			out = append(out, ',', ' ')
		}
		out = appendByteHex(out, char)
	}

	state.indent--
	out = append(out, ',', '\n')
	out = appendIndent(out, state)
	out = append(out, '}')
	return out
}

func appendByteHex(out []byte, char byte) []byte {
	const hexDigits = "0123456789abcdef"
	return append(out, '0', 'x', hexDigits[int(char>>4)], hexDigits[int(char&^0xf0)])
}

func byteArrayToSlice(rval reflect.Value) []byte {
	type sliceHeader struct {
		dat unsafe.Pointer
		len int
		cap int
	}
	ptr, size := raw(rval)
	slice := sliceHeader{ptr, int(size), int(size)}
	return *(*[]byte)(unsafe.Pointer(&slice))
}

func appendCastPrefix(out []byte, rval reflect.Value, state state) []byte {
	if state.elideType {
		return out
	}
	out = appendTypeName(out, rval.Type(), state)
	out = append(out, '(')
	return out
}

func appendCastSuffix(out []byte, rval reflect.Value, state state) []byte {
	if state.elideType {
		return out
	}
	return append(out, ')')
}

func appendIndent(out []byte, state state) []byte {
	for i := 0; i < state.indent; i++ {
		out = append(out, '\t')
	}
	return out
}

func isZeroOrShouldOmit(rval reflect.Value) bool {
	switch rval.Type().Kind() {
	case reflect.Bool:
		return !rval.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rval.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rval.Uint() == 0

	case reflect.Uintptr:
		return rval.Uint() == 0

	case reflect.UnsafePointer:
		return rval.Convert(reflect.TypeOf(unsafe.Pointer(nil))).Interface().(unsafe.Pointer) == nil

	case reflect.Float32:
		return rval.Float() == 0

	case reflect.Float64:
		return rval.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return rval.Complex() == 0

	case reflect.Array:
		return isZero(rval)

	case reflect.Chan:
		return true

	case reflect.Func:
		return true

	case reflect.Interface:
		return rval.Interface() == nil

	case reflect.Map:
		return rval.IsNil()

	case reflect.Ptr:
		return rval.IsNil()

	case reflect.Slice:
		return rval.IsNil()

	case reflect.String:
		return rval.String() == ""

	case reflect.Struct:
		return isZero(rval)

	default:
		return false
	}
}

func mayRequireMultiline(rtype reflect.Type) bool {
	switch rtype.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.String, reflect.Struct:
		return true
	default:
		return false
	}
}

func isPrimitive(rtype reflect.Type) bool {
	switch rtype.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

func isInterface(rtype reflect.Type) bool {
	return rtype.Kind() == reflect.Interface
}

func isNil(rval reflect.Value) bool {
	switch rval.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rval.IsNil()
	default:
		return false
	}
}

func appendTypeName(out []byte, rtype reflect.Type, state state) []byte {
	name := rtype.Name()

	if name == "" {
		switch rtype.Kind() {
		case reflect.Array:
			out = append(out, '[')
			out = strconv.AppendInt(out, int64(rtype.Len()), 10)
			out = append(out, ']')
			out = appendTypeName(out, rtype.Elem(), state)
			return out

		case reflect.Slice:
			out = append(out, "[]"...)
			out = appendTypeName(out, rtype.Elem(), state)
			return out

		case reflect.Map:
			out = append(out, "map["...)
			out = appendTypeName(out, rtype.Key(), state)
			out = append(out, ']')
			out = appendTypeName(out, rtype.Elem(), state)
			return out
		}
		return append(out, rtype.String()...)
	}

	pkg := rtype.PkgPath()
	if pkg == "" {
		return append(out, rtype.String()...)
	}

	pkg, ok := state.conf.PackageMap[pkg]
	if !ok {
		return append(out, rtype.String()...)
	}

	if pkg == "" {
		return append(out, name...)
	}

	out = append(out, pkg...)
	out = append(out, '.')
	out = append(out, name...)
	return out
}

// Questionable
func isZero(rval reflect.Value) bool {
	ptr, size := raw(rval)
	for i := uintptr(0); i < size; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(ptr) + i)) != 0 {
			return false
		}
	}
	return true
}

func raw(rval reflect.Value) (unsafe.Pointer, uintptr) {
	if rval.CanAddr() {
		return unsafe.Pointer(rval.UnsafeAddr()), rval.Type().Size()
	}

	type emptyInterface struct {
		rtype uintptr
		dat   unsafe.Pointer
	}
	iface := rval.Interface()
	return (*emptyInterface)(unsafe.Pointer(&iface)).dat, rval.Type().Size()
}

/*
Reinterprets a byte slice as a string, saving an allocation.
Borrowed from the standard library. Reasonably safe.
*/
func bytesToMutableString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func isSfieldExported(sfield reflect.StructField) bool {
	return sfield.PkgPath == ""
}

func canElideType(rtype reflect.Type, state state) bool {
	return !state.conf.ForceConstructorName && !isInterface(rtype)
}
