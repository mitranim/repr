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

	* Dumping data as code is useful for code generation and debugging

	* fmt.Sprintf("%#v") doesn't always output valid code, has no multiline mode

	* https://github.com/davecgh/go-spew/spew doesn't output valid code, has no single-line mode

	* https://github.com/shurcooL/go-goon outputs too much noise, has no single-line mode

Features

Supports single-line and multiline modes. Defaults to multiline.

The output looks like something you'd write by hand, and is almost exactly
compliant with gofmt. Unlike gofmt, it doesn't align field values in struct
literals. Use "go/format" to fix that, at a 50x performance cost:

	import (
		"go/format"
		"github.com/Mitranim/repr"
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

	* Fancy types such as "big.Int" or "time.Time" are printed as empty structs;
	  ideally they would be printed as constructor calls

	* Funcs are treated as nil

	* Chans are treated as nil

	* Pointers to primitive types are not supported and cause a panic

	* "byte" is printed as "uint8"

	* "rune" is printed as "int32"

	* Runes are printed as integers, not character literals

	* Enum-style constants are not mapped back to identifers

	* Omits private fields from structs

	* Cyclic structures cause infinite recursion

Pointers to composite types such as structs, arrays, slices and maps are
supported by prefixing literals with "&". Go currently doesn't support this for
primitive literals. On the bright side, you probably shouldn't use pointers to
primitive types anyway.

Installation

Shell:

	go get -u github.com/Mitranim/repr

Usage

Example:

	import "github.com/Mitranim/repr"

	type Data struct {
		Number int
		String string
		List   []int
	}

	fmt.Println(repr.String(Data{
		Number: 123,
		String: "hello world!",
		List:   []int{10, 20, 30},
	}))

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
Global settings. Used by "String", "Bytes" and "Append". Feel free to modify.
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
	return bytesToMutableString(appendAny(nil, val, Default, false, 0))
}

/*
Short for "String with config". Formats the value using the provided config. See
"Config" for details.
*/
func StringC(val interface{}, conf Config) string {
	return bytesToMutableString(appendAny(nil, val, conf, false, 0))
}

/*
Formats the value using the "Default" config. See "Config" for details.
*/
func Bytes(val interface{}) []byte {
	return appendAny(nil, val, Default, false, 0)
}

/*
"Short for "Bytes with config". Formats the value using the provided config. See
"Config" for details.
*/
func BytesC(val interface{}, conf Config) []byte {
	return appendAny(nil, val, conf, false, 0)
}

/*
Formats the value using the "Default" config, appending the output to the
provided buffer. See "Config" for details.
*/
func Append(out []byte, val interface{}) []byte {
	return appendAny(nil, val, Default, false, 0)
}

/*
Short for "Append with config". Formats the value using the provided config,
appending the output to the provided buffer. See "Config" for details.
*/
func AppendC(out []byte, val interface{}, conf Config) []byte {
	return appendAny(out, val, conf, false, 0)
}

var (
	byteType      = reflect.TypeOf(byte(0))
	byteSliceType = reflect.TypeOf([]byte(nil))
)

func appendAny(out []byte, val interface{}, conf Config, elideType bool, indent int) []byte {
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
		if !elideType {
			out = append(out, "[]uint8"...)
		}
		out = appendBytes(out, val, conf, indent)
		return out
	}

	rval := reflect.ValueOf(val)
	typ := rval.Type()

	switch typ.Kind() {
	case reflect.Bool:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		if rval.Bool() {
			out = append(out, "true"...)
		} else {
			out = append(out, "false"...)
		}
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendInt(out, rval.Int(), 10)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendUint(out, rval.Uint(), 10)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Uintptr:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendUint(append(out, '0', 'x'), rval.Uint(), 16)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Float32:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendFloat(out, rval.Float(), 'f', -1, 32)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Float64:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendFloat(out, rval.Float(), 'f', -1, 64)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Complex64, reflect.Complex128:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = appendComplex128(out, rval.Convert(reflect.TypeOf(complex128(0))).Complex())
		out = appendCastSuffix(out, rval, elideType)

	case reflect.String:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = strconv.AppendQuote(out, rval.String())
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Chan:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = append(out, "nil"...)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Func:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		out = append(out, "nil"...)
		out = appendCastSuffix(out, rval, elideType)

	// Pretty sure this should never match
	case reflect.Interface:
		panic("repr currently doesn't support printing an interface")

	case reflect.UnsafePointer:
		out = appendCastPrefix(out, rval, elideType, conf.PackageMap)
		ptr := rval.Convert(reflect.TypeOf(unsafe.Pointer(nil))).Interface().(unsafe.Pointer)
		out = strconv.AppendUint(append(out, '0', 'x'), uint64(uintptr(ptr)), 16)
		out = appendCastSuffix(out, rval, elideType)

	case reflect.Ptr:
		switch typ.Elem().Kind() {
		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			if isZeroOrShouldOmit(rval) {
				out = append(out, "nil"...)
			} else {
				out = append(out, '&')
				out = appendAny(out, rval.Elem().Interface(), conf, elideType, indent)
			}
		default:
			panic("repr currently doesn't support pointers to non-composite types")
		}

	case reflect.Array:
		if !elideType {
			out = append(out, typeName(rval.Type(), conf.PackageMap)...)
		}
		if typ.Elem() == byteType {
			out = appendBytes(out, byteArrayToSlice(rval), conf, indent)
		} else {
			out = appendList(out, rval, conf, indent)
		}

	case reflect.Slice:
		if !elideType {
			out = append(out, typeName(rval.Type(), conf.PackageMap)...)
		}
		if rval.IsNil() {
			if elideType {
				out = append(out, "nil"...)
			} else {
				out = append(out, "(nil)"...)
			}
		} else if typ.Elem() == byteType {
			out = appendBytes(out, rval.Bytes(), conf, indent)
		} else {
			out = appendList(out, rval, conf, indent)
		}

	case reflect.Struct:
		if !elideType {
			out = append(out, typeName(rval.Type(), conf.PackageMap)...)
		}
		out = appendStruct(out, rval, conf, indent)

	case reflect.Map:
		if !elideType {
			out = append(out, typeName(rval.Type(), conf.PackageMap)...)
		}
		if rval.IsNil() {
			if elideType {
				out = append(out, "nil"...)
			} else {
				out = append(out, "(nil)"...)
			}
		} else {
			out = appendMap(out, rval, conf, indent)
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

func appendList(out []byte, rval reflect.Value, conf Config, indent int) []byte {
	elemType := rval.Type().Elem()
	elideElemType := isPrimitive(elemType)
	count := rval.Len()

	if conf.SingleLine || (!mayRequireMultiline(elemType) && count < 48) {
		out = append(out, '{')
		for i := 0; i < count; i++ {
			out = appendAny(out, rval.Index(i).Interface(), conf, elideElemType, 0)
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
		indent++
	}

	for i := 0; i < count; i++ {
		out = appendIndent(out, indent)
		out = appendAny(out, rval.Index(i).Interface(), conf, elideElemType, indent)
		out = append(out, ',', '\n')
	}

	if count > 0 {
		indent--
		out = appendIndent(out, indent)
	}

	out = append(out, '}')
	return out
}

func appendStruct(out []byte, rval reflect.Value, conf Config, indent int) []byte {
	typ := rval.Type()

	if conf.SingleLine {
		out = append(out, '{')
		count := rval.NumField()
		var hasFields bool

		for i := 0; i < count; i++ {
			field := rval.Field(i)
			if !conf.ZeroFields && isZeroOrShouldOmit(field) {
				continue
			}

			fieldType := typ.Field(i)
			// Skip unexported field
			if fieldType.PkgPath != "" {
				continue
			}

			if hasFields {
				out = append(out, ',', ' ')
			}
			hasFields = true

			out = append(out, fieldType.Name...)
			out = append(out, ':', ' ')
			elideElemType := isPrimitive(field.Type()) || isNil(field)
			out = appendAny(out, field.Interface(), conf, elideElemType, 0)
		}
		out = append(out, '}')
		return out
	}

	out = append(out, '{')
	count := 0

	for i := 0; i < rval.NumField(); i++ {
		field := rval.Field(i)
		if !conf.ZeroFields && isZeroOrShouldOmit(field) {
			continue
		}

		fieldType := typ.Field(i)
		// Skip unexported field
		if fieldType.PkgPath != "" {
			continue
		}

		count++
		if count == 1 {
			out = append(out, '\n')
			indent++
		}

		out = appendIndent(out, indent)
		out = append(out, fieldType.Name...)
		out = append(out, ':', ' ')
		elideElemType := isPrimitive(field.Type()) || isNil(field)
		out = appendAny(out, field.Interface(), conf, elideElemType, indent)
		out = append(out, ',', '\n')
	}

	if count > 0 {
		indent--
		out = appendIndent(out, indent)
	}

	out = append(out, '}')
	return out
}

func appendMap(out []byte, rval reflect.Value, conf Config, indent int) []byte {
	typ := rval.Type()
	keyType := typ.Key()
	elemType := typ.Elem()
	elideKeyType := isPrimitive(keyType)
	elideElemType := isPrimitive(elemType)

	if conf.SingleLine {
		out = append(out, '{')
		keys := rval.MapKeys()

		for i, key := range keys {
			out = appendAny(out, key.Interface(), conf, elideKeyType, 0)
			out = append(out, ':', ' ')
			out = appendAny(out, rval.MapIndex(key).Interface(), conf, elideElemType, 0)
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
			indent++
		}

		out = appendIndent(out, indent)
		out = appendAny(out, key.Interface(), conf, elideKeyType, indent)
		out = append(out, ':', ' ')
		out = appendAny(out, rval.MapIndex(key).Interface(), conf, elideElemType, indent)

		out = append(out, ',', '\n')
	}

	if len(keys) > 0 {
		indent--
		out = appendIndent(out, indent)
	}

	out = append(out, '}')
	return out
}

// Similar to fmt.Sprintf("%#02v", val), but multiline: large inputs are printed
// as a column with 8 bytes per row.
func appendBytes(out []byte, val []byte, conf Config, indent int) []byte {
	if conf.SingleLine || len(val) <= 8 {
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

	indent++
	out = append(out, '{', '\n')

	for i, char := range val {
		if i == 0 {
			out = appendIndent(out, indent)
		} else if i%8 == 0 {
			out = append(out, ',', '\n')
			out = appendIndent(out, indent)
		} else {
			out = append(out, ',', ' ')
		}
		out = appendByteHex(out, char)
	}

	indent--
	out = append(out, ',', '\n')
	out = appendIndent(out, indent)
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

func appendCastPrefix(out []byte, rval reflect.Value, elideType bool, packageMap map[string]string) []byte {
	if elideType {
		return out
	}
	out = append(out, typeName(rval.Type(), packageMap)...)
	out = append(out, '(')
	return out
}

func appendCastSuffix(out []byte, rval reflect.Value, elideType bool) []byte {
	if elideType {
		return out
	}
	return append(out, ')')
}

func appendIndent(out []byte, indent int) []byte {
	for i := 0; i < indent; i++ {
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

func mayRequireMultiline(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.String, reflect.Struct:
		return true
	default:
		return false
	}
}

func isPrimitive(typ reflect.Type) bool {
	switch typ.Kind() {
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

func isNil(rval reflect.Value) bool {
	switch rval.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rval.IsNil()
	default:
		return false
	}
}

// Inefficient, but not a bottleneck. TODO improve.
func typeName(typ reflect.Type, packageMap map[string]string) string {
	name := typ.Name()
	if name == "" {
		switch typ.Kind() {
		case reflect.Array:
			return "[" + strconv.Itoa(typ.Len()) + "]" + typeName(typ.Elem(), packageMap)

		case reflect.Slice:
			return "[]" + typeName(typ.Elem(), packageMap)

		case reflect.Map:
			return "map[" + typeName(typ.Key(), packageMap) + "]" + typeName(typ.Elem(), packageMap)
		}
		return typ.String()
	}

	pkg := typ.PkgPath()
	if pkg == "" {
		return typ.String()
	}

	pkg, ok := packageMap[pkg]
	if !ok {
		return typ.String()
	}
	if pkg == "" {
		return name
	}
	return pkg + "." + name
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
		typ uintptr
		dat unsafe.Pointer
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
