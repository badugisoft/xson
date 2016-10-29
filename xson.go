package xson

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"

	yaml "gopkg.in/yaml.v2"
)

type Type int

const (
	UNKNOWN Type = 1 + iota
	JSON
	YAML
	XML
	TOML
	FLAT_JSON
	FLAT_YAML
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
		TOML: Marshaller{
			Marshal: func(v interface{}) ([]byte, error) {
				var b bytes.Buffer
				err := toml.NewEncoder(&b).Encode(v)
				if err != nil {
					return nil, err
				}
				return b.Bytes(), nil
			},
			MarshalIndent: nil,
			Unmarshal:     toml.Unmarshal,
			Extesions:     []string{"toml"},
		},
		FLAT_JSON: Marshaller{
			Marshal:       nil,
			MarshalIndent: nil,
			Unmarshal: func(data []byte, v interface{}) error {
				tmp := map[interface{}]interface{}{}
				err := json.Unmarshal(data, &tmp)
				if err != nil {
					return err
				}

				err = unflatten(&tmp)
				if err != nil {
					return err
				}

				d, err := json.Marshal(tmp)
				if err != nil {
					return err
				}

				return json.Unmarshal(d, v)
			},
			Extesions: []string{"flat.json"},
		},
		FLAT_YAML: Marshaller{
			Marshal:       nil,
			MarshalIndent: nil,
			Unmarshal: func(data []byte, v interface{}) error {
				tmp := map[interface{}]interface{}{}
				err := yaml.Unmarshal(data, &tmp)
				if err != nil {
					return err
				}

				err = unflatten(&tmp)
				if err != nil {
					return err
				}

				d, err := yaml.Marshal(tmp)
				if err != nil {
					return err
				}

				return yaml.Unmarshal(d, v)
			},
			Extesions: []string{"flat.yaml", "flat.yml"},
		},
	}
)

func unflatten(m *map[interface{}]interface{}) error {
	for k, v := range *m {
		kk := k.(string)
		if strings.Contains(kk, ".") {
			kind := reflect.TypeOf(v).Kind()
			if kind == reflect.Map {
				vv := v.(map[interface{}]interface{})
				err := unflatten(&vv)
				if err != nil {
					return err
				}
				v = vv
			}

			p := *m
			tokens := strings.Split(kk, ".")
			for i, token := range tokens {
				if i == len(tokens)-1 {
					p[token] = v
				} else {
					if _, found := p[token]; !found {
						p[token] = map[interface{}]interface{}{}
					}
					p, _ = p[token].(map[interface{}]interface{})
				}
			}

			delete(*m, k)
		}
	}
	return nil
}

func GetType(filenameOrExtension string) Type {
	lowerName := strings.ToLower(filenameOrExtension)
	maxExtLen, tp := 0, UNKNOWN

	for t, m := range marshallers {
		for _, e := range m.Extesions {
			if e == lowerName || strings.HasSuffix(lowerName, "."+e) {
				extLen := len(e)
				if extLen > maxExtLen {
					tp = t
					maxExtLen = extLen
				}
				break
			}
		}
	}
	return tp
}

func GetTypes() []Type {
	return []Type{JSON, YAML, XML, TOML, FLAT_JSON, FLAT_YAML}
}

func GetExtensions(t Type) []string {
	marshaller, found := marshallers[t]
	if !found {
		return []string{}
	}
	return marshaller.Extesions[:]
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
