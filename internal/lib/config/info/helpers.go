package info

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	whiteSpace = 32
)

type StringerWriter interface {
	io.Writer
	fmt.Stringer
}

type indentWrapper struct {
	writer StringerWriter
	indent int
}

func newIndent(writer StringerWriter, indent int) *indentWrapper {
	return &indentWrapper{writer: writer, indent: indent}
}

func (i *indentWrapper) Write(p []byte) (n int, err error) {
	indenter := []byte{whiteSpace}
	return i.writer.Write(append(bytes.Repeat(indenter, i.indent), p...))
}

func (i *indentWrapper) Bytes() []byte {
	return []byte(i.writer.String())
}

func ConfigInfoTags[Type any](cfg Type) []byte {

	_type := reflect.TypeOf(cfg)

	if _type.Kind() == reflect.Pointer {
		_type = _type.Elem()
	}

	if _type.Kind() != reflect.Struct {
		return nil
	}

	builder := new(strings.Builder)
	builder.Write([]byte(_type.Name() + "\n"))

	writer := newIndent(builder, 1)
	configInfoTags(writer, _type)

	return writer.Bytes()
}

func configInfoTags(
	writer io.Writer,
	_type reflect.Type,
) {

	for i := range _type.NumField() {

		field := _type.Field(i)

		if field.Type.Kind() == reflect.Struct {
			configInfoTags(writer, field.Type)
			return
		}

		writer.Write([]byte(fmt.Sprintf("%s %s", field.Name, field.Type)))

		if envTag, ok := field.Tag.Lookup("env"); ok {
			writer.Write([]byte("env:" + envTag))
		}

		if envTag, ok := field.Tag.Lookup("yaml"); ok {
			writer.Write([]byte("yaml:" + envTag))
		}

		if envTag, ok := field.Tag.Lookup("json"); ok {
			writer.Write([]byte("json:" + envTag))
		}

		if defaultTag, ok := field.Tag.Lookup("default"); ok {
			writer.Write([]byte("default:" + defaultTag))
		}

		if descTag, ok := field.Tag.Lookup("desc"); ok {
			writer.Write([]byte("- " + descTag))
		}

		writer.Write([]byte("\n"))
	}

}
