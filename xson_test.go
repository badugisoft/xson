package xson_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/badugisoft/xson"
	"github.com/badugisoft/xson/data"
)

func TestMarshal(t *testing.T) {
	types := map[xson.Type]string{
		xson.JSON: "data/data.json",
		xson.YAML: "data/data.yaml",
		xson.XML:  "data/data.xml",
		xson.TOML: "data/data.toml",
	}

	for tp, file := range types {
		d, err := xson.Marshal(tp, data.SampleData)
		_panic(err)

		read, err := ioutil.ReadFile(file)
		_panic(err)

		_true(t, reflect.DeepEqual(_trimBytes(d), _trimBytes(read)))
	}
}

func TestMarshalIndent(t *testing.T) {
	types := map[xson.Type]string{
		xson.JSON: "data/data.indent.json",
		xson.XML:  "data/data.indent.xml",
	}

	for tp, file := range types {
		d, err := xson.MarshalIndent(tp, data.SampleData, "", "  ")
		_panic(err)

		read, err := ioutil.ReadFile(file)
		_panic(err)

		_true(t, reflect.DeepEqual(_trimBytes(d), _trimBytes(read)))
	}
}

func TestUnmarshal(t *testing.T) {
	types := map[xson.Type]string{
		xson.JSON:      "data/data.json",
		xson.YAML:      "data/data.yaml",
		xson.FLAT_YAML: "data/data.flat.yaml",
		xson.XML:       "data/data.xml",
		xson.TOML:      "data/data.toml",
	}

	for tp, file := range types {
		read, err := ioutil.ReadFile(file)
		_panic(err)

		var d data.Data
		err = xson.Unmarshal(tp, read, &d)
		_panic(err)

		_true(t, reflect.DeepEqual(data.SampleData, d))
	}
}
