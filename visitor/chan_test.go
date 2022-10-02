package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChan_Direction(t *testing.T) {
	c := &Chan{
		direction: SendDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, SendDir, c.Direction())

	c = &Chan{
		direction: ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, ReceiveDir, c.Direction())

	c = &Chan{
		direction: SendDir | ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, BothDir, c.Direction())
}

func TestChan_Name(t *testing.T) {
	c := &Chan{
		direction: SendDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "chan<- bool", c.Name())

	c = &Chan{
		direction: ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "<-chan bool", c.Name())

	c = &Chan{
		direction: SendDir | ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "chan bool", c.Name())
}

func TestChan_String(t *testing.T) {
	c := &Chan{
		direction: SendDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "chan<- bool", c.String())

	c = &Chan{
		direction: ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "<-chan bool", c.String())

	c = &Chan{
		direction: SendDir | ReceiveDir,
		elem:      basicTypesMap["bool"],
	}

	assert.Equal(t, "chan bool", c.String())
}

func TestChan_Elem(t *testing.T) {
	elem := basicTypesMap["bool"]
	c := &Chan{
		direction: SendDir | ReceiveDir,
		elem:      elem,
	}

	assert.Equal(t, elem, c.Elem())
	assert.Equal(t, "bool", c.Elem().Name())
	assert.Equal(t, "bool", c.Elem().String())
}

func TestChan_Underlying(t *testing.T) {
	elem := basicTypesMap["bool"]
	c := &Chan{
		direction: SendDir | ReceiveDir,
		elem:      elem,
	}

	assert.Equal(t, c, c.Underlying())
}
