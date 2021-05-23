package codec

import (
	"bytes"
	"errors"
	"github.com/appmanch/go-commons/textutils"
	"io"
	"reflect"
	"strings"
)

//var knownTypes map[string][]FieldMeta=make(map[string][]FieldMeta)
// Bool
//	Int
//	Int8
//	Int16
//	Int32
//	Int64
//	Uint
//	Uint8
//	Uint16
//	Uint32
//	Uint64
//	Float32
//	Float64
//	Array
//	Func
//	Interface
//	Map
//	Slice
//	String
//	Struct

type FieldMeta struct{
	Name       string
	FieldName  string
	Type       reflect.Type
	Dimension  int
	Required bool
	TargetNames map[string]string

}

type StringFieldMeta struct {
	FieldMeta
	DefaultVal string
	Pattern string
	Format string
	OmitEmpty bool
}

type Int8FieldMeta struct {
	FieldMeta
	DefaultVal int8
	Min int8
	Max int8
}

type Int16FieldMeta struct {
	FieldMeta
	DefaultVal int16
	Min int16
	Max int16
}





type IntFieldMeta struct {
	FieldMeta
	DefaultVal int
	Min int
	Max int
}

type Int64FieldMeta struct {
	FieldMeta
	DefaultVal int64
	Min int64
	Max int64
}

type Float32FieldMeta struct {
	FieldMeta
	DefaultVal float32
	Min float32
	Max float32
}


type Float64FieldMeta struct {
	FieldMeta
	DefaultVal float64
	Min float64
	Max float64
}

type BooleanFieldMeta struct {
	FieldMeta
	DefaultVal bool
}

type BooleanFieldMeta struct {
	FieldMeta
	DefaultVal bool
}

type baseCodec struct {

}

func Get() baseCodec{
	return baseCodec{}
}



func handleType(v interface{}){

}

func (d baseCodec) DecodeString(s string, v interface{}) error {

	r := strings.NewReader(s)
	return d.Read(r, v)
}

func (d baseCodec) DecodeBytes(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	return d.Read(r, v)
}

func (d baseCodec) EncodeToBytes(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	e := d.Write(v, buf)
	if e == nil {
		return buf.Bytes(), e
	} else {
		return nil, e
	}
}

func (d baseCodec) EncodeToString(v interface{}) (string, error) {
	buf := &bytes.Buffer{}
	e := d.Write(v, buf)
	if e == nil {
		return buf.String(), e
	} else {
		return textutils.EmptyStr, e
	}
}

func (d baseCodec) Read(r io.Reader, v interface{}) error {

	return errors.New("Reader is not implemented in base codec")
}

func (d baseCodec) Write(v interface{}, w io.Writer) error {
	return errors.New("Writer is not implemented in base codec")
}

type StringEncoder interface {
	//EncodeToString will encode  a type to string
	EncodeToString(v interface{}) string
}

type BytesEncoder interface {
	// EncodeToBytes will encode the provided type to []byte
	EncodeToBytes(v interface{}) []byte
}

type StringDecoder interface {
	//DecodeString will decode  a type from string
	DecodeString(s string, v interface{}) error
}

type BytesDecoder interface {
	//DecodeBytes will decode a type from an array of bytes
	DecodeBytes(b []byte, v interface{}) error
}

type Encoder interface {
	StringEncoder
	BytesEncoder
}


type EncoderWriter interface {
	Encoder
	//Write a type to writer
	Write(v interface{}, w io.Writer) error
}

type Decoder interface {
	StringDecoder
	BytesDecoder
}



type DecoderReader interface {
	Decoder
	//Read a type from a reader
	Read(r io.Reader, v interface{}) error
}

//Generic Codec Interface
type Codec interface {
	EncoderWriter
	DecoderReader
}
