package common

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type OtherStruct struct {
	Data       []byte
	Other      string  `json:"other"`
	Another    int     `json:",omitempty"`
	FloatPropA float64 `json:"float_property_a"`
	FloatPropB float64 `json:"float_property_b,omitempty"`
	Nested     *NestedStruct
}

type NestedStruct struct {
	InnerData []byte `json:"-,omitempty"`
}

type Example struct {
	Data *OtherStruct
}

func (e *Example) MarshalJSON() ([]byte, error) {
	serialized := SerializeBytes(e)
	return json.Marshal(serialized)
}

func (e *Example) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	deserialized := DeserializeBytes(raw, reflect.TypeOf(*e))
	*e = deserialized.(Example)
	return nil
}

func TestKeyExists(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": {"b": {"c": 1}}}`)
	asst.True(KeyExists(data, "a", "b", "c"), "test KeyPathExists() failed")
	asst.False(KeyExists(data, "a", "b", "c", "d"), "test KeyPathExists() failed")
}

func TestKeyPathExists(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": {"b": {"c": 1}}}`)
	asst.True(KeyPathExists(data, "a.b.c"), "test KeyPathExists() failed")
	asst.False(KeyPathExists(data, "a.b.c"), "test KeyPathExists() failed")
}

func TestGetLength(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": [1, 2, 3]}`)
	length, err := GetLength(data, "a")
	asst.Nil(err, "test GetLength() failed")
	asst.Equal(3, length, "test GetLength() failed")
}

func TestSerializeBytes(t *testing.T) {
	asst := assert.New(t)

	nested := &NestedStruct{
		InnerData: []byte("Nested data!"),
	}

	os := &OtherStruct{
		Data:       []byte("Hello, World!"),
		Other:      "Some other data",
		Another:    123,
		FloatPropA: 1.23,
		FloatPropB: 4.0,
		Nested:     nested,
	}

	tBytes, err := json.Marshal(t)
	asst.Nil(err, "test SerializeBytes() failed")
	t.Logf("T Serialized: %s", string(tBytes))

	ex1 := &Example{
		Data: os,
	}

	// serialize Example struct
	serialized, err := json.Marshal(ex1)
	asst.Nil(err, "test SerializeBytes() failed")
	t.Logf("Serialized: %s", string(serialized))

	var ex2 Example
	err = json.Unmarshal(serialized, &ex2)
	asst.Nil(err, "test SerializeBytes() failed")
	t.Logf("Deserialized: %+v", ex2)

	asst.Equal(ex1.Data.Data, ex2.Data.Data, "test SerializeBytes() failed")
	asst.Equal(ex1.Data.Other, ex2.Data.Other, "test SerializeBytes() failed")
	asst.Equal(ex1.Data.Another, ex2.Data.Another, "test SerializeBytes() failed")
	asst.Equal(ex1.Data.FloatPropA, ex2.Data.FloatPropA, "test SerializeBytes() failed")
	asst.Equal(ex1.Data.FloatPropB, ex2.Data.FloatPropB, "test SerializeBytes() failed")
	asst.Equal(ex1.Data.Nested.InnerData, ex2.Data.Nested.InnerData, "test SerializeBytes() failed")
}

func TestJSON_MaskJSON(t *testing.T) {
	asst := assert.New(t)

	data := []byte(`{"a": {"b": {"password": "sssss"}}}`)
	masked, err := MaskJSON(data, "pass")
	asst.Nil(err, "test MaskJSON() failed")
	asst.Equal(`{"a":{"b":{"password":"******"}}}`, string(masked), "test MaskJSON() failed")
}
