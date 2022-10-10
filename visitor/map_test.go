package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap_Name(t *testing.T) {
	key := basicTypesMap["string"]
	val := &Slice{
		elem: basicTypesMap["int32"],
	}
	m := &Map{
		key:  key,
		elem: val,
	}

	assert.Equal(t, "map[string][]int32", m.Name())
}

func TestMap_String(t *testing.T) {
	key := basicTypesMap["string"]
	val := &Interface{}
	m := &Map{
		key:  key,
		elem: val,
	}

	assert.Equal(t, "map[string]interface{}", m.String())
}

func TestMap_Key(t *testing.T) {
	key := basicTypesMap["bool"]
	val := &Array{
		len:  5,
		elem: basicTypesMap["rune"],
	}
	m := &Map{
		key:  key,
		elem: val,
	}

	assert.Equal(t, key, m.Key())
	assert.Equal(t, "bool", m.Key().Name())
	assert.Equal(t, "bool", m.Key().String())
}

func TestMap_Elem(t *testing.T) {
	key := basicTypesMap["string"]
	val := &Slice{
		elem: basicTypesMap["int32"],
	}
	m := &Map{
		key:  key,
		elem: val,
	}

	assert.Equal(t, val, m.Elem())
	assert.Equal(t, "[]int32", m.Elem().Name())
	assert.Equal(t, "[]int32", m.Elem().String())

	assert.Equal(t, m, m.Underlying())
}

func TestMap_Underlying(t *testing.T) {
	key := basicTypesMap["string"]
	val := &Slice{
		elem: basicTypesMap["int32"],
	}
	m := &Map{
		key:  key,
		elem: val,
	}

	assert.Equal(t, m, m.Underlying())
}
