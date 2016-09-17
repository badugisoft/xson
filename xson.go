package xson

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Type int

const (
	UNKNOWN Type = 1 + iota
	JSON
	YAML
	XML
)

type Marshaller struct {
	Marshal       func(v interface{}) ([]byte, error)
	MarshalIndent func(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal     func(data []byte, v interface{}) error
	Extesions     []string
}

var (
	UnknownTypeError  = errors.New("Unknown type")
	NotSupportedError = errors.New("Not supported function")

	marshallers = map[Type]Marshaller{
		JSON: Marshaller{
			Marshal:       json.Marshal,
			MarshalIndent: json.MarshalIndent,
			Unmarshal:     json.Unmarshal,
			Extesions:     []string{"json"},
		},
		YAML: Marshaller{
			Marshal:       yaml.Marshal,
			MarshalIndent: nil,
			Unmarshal:     yaml.Unmarshal,
			Extesions:     []string{"yaml", "yml"},
		},
		XML: Marshaller{
			Marshal:       xml.Marshal,
			MarshalIndent: xml.MarshalIndent,
			Unmarshal:     xml.Unmarshal,
			Extesions:     []string{"xml"},
		},
	}
)

func GetType(filenameOrExtension string) Type {
	filename := strings.ToLower(filenameOrExtension)
	extension := filepath.Ext(filename)
	for t, m := range marshallers {
		for _, e := range m.Extesions {
			if e == filename || e == extension {
				return t
			}
		}
	}
	return UNKNOWN
}

func GetTypes() []Type {
	return []Type{JSON, YAML, XML}
}

func Marshal(t Type, v interface{}) ([]byte, error) {
	marshaller, found := marshallers[t]
	if !found {
		return nil, UnknownTypeError
	}
	if marshaller.Marshal == nil {
		return nil, NotSupportedError
	}
	return marshaller.Marshal(v)
}

func MarshalIndent(t Type, v interface{}, prefix, indent string) ([]byte, error) {
	marshaller, found := marshallers[t]
	if !found {
		return nil, UnknownTypeError
	}
	if marshaller.MarshalIndent == nil {
		return nil, NotSupportedError
	}
	return marshaller.MarshalIndent(v, prefix, indent)
}

func Unmarshal(t Type, data []byte, v interface{}) error {
	marshaller, found := marshallers[t]
	if !found {
		return UnknownTypeError
	}
	if marshaller.Unmarshal == nil {
		return NotSupportedError
	}
	return marshaller.Unmarshal(data, v)
}
