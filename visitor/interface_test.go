package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/stretchr/testify/assert"
	"testing"
)

type interfaceInfo struct {
	markers         marker.MarkerValues
	name            string
	fileName        string
	position        Position
	explicitMethods map[string]functionInfo
	methods         map[string]functionInfo
	embeddedTypes   []string
	isExported      bool
}

// interfaces
var (
	bakeryShopInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "BakeryShop",
				},
			},
		},
		name:       "BakeryShop",
		fileName:   "dessert.go",
		isExported: true,
		position: Position{
			Line:   13,
			Column: 6,
		},
		explicitMethods: map[string]functionInfo{
			"Bread": breadFunction,
		},
		methods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
			"Bread":    breadFunction,
		},
		embeddedTypes: []string{"Dessert"},
	}

	dessertInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "Dessert",
				},
			},
		},
		name:       "Dessert",
		fileName:   "dessert.go",
		isExported: true,
		position: Position{
			Line:   79,
			Column: 6,
		},
		explicitMethods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
		methods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
	}

	newYearsEveCookieInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "newYearsEveCookie",
				},
			},
		},
		name:       "newYearsEveCookie",
		fileName:   "dessert.go",
		isExported: false,
		position: Position{
			Line:   48,
			Column: 6,
		},
		methods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
		},
		explicitMethods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
		},
	}

	sweetShopInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "SweetShop",
				},
			},
		},
		name:       "SweetShop",
		fileName:   "dessert.go",
		isExported: true,
		position: Position{
			Line:   125,
			Column: 6,
		},
		explicitMethods: map[string]functionInfo{
			"Macaron": macaronFunction,
		},
		methods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
			"Macaron":  macaronFunction,
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
		embeddedTypes: []string{"newYearsEveCookie", "Dessert"},
	}
)

func assertInterfaces(t *testing.T, file *File, interfaces map[string]interfaceInfo) bool {

	if len(interfaces) != file.Interfaces().Len() {
		t.Errorf("the number of the interface should be %d, but got %d", len(interfaces), file.Interfaces().Len())
		return false
	}

	index := 0
	for expectedInterfaceName, expectedInterface := range interfaces {
		actualInterface, ok := file.Interfaces().FindByName(expectedInterfaceName)

		if !ok {
			t.Errorf("interface with name %s is not found", expectedInterfaceName)
			continue
		}

		if actualInterface.InterfaceType() == nil {
			t.Errorf("InterfaceType() for interface %s should not return nil", actualInterface.Name())
		}

		if expectedInterface.fileName != actualInterface.File().Name() {
			t.Errorf("the file name for interface %s should be %s, but got %s", expectedInterfaceName, expectedInterface.fileName, actualInterface.File().Name())
		}

		if file.Interfaces().elements[index] != file.Interfaces().At(index) {
			t.Errorf("interface with name %s does not match with interface at index %d", actualInterface.Name(), index)
			continue
		}

		if actualInterface.IsExported() && !expectedInterface.isExported {
			t.Errorf("interface with name %s is exported, but should be unexported", actualInterface.Name())
		} else if !actualInterface.IsExported() && expectedInterface.isExported {
			t.Errorf("interface with name %s is not exported, but should be exported", actualInterface.Name())
		}

		if actualInterface.NumMethods() == 0 && !actualInterface.IsEmpty() {
			t.Errorf("the interface %s should be empty", actualInterface.Name())
		} else if actualInterface.NumMethods() != 0 && actualInterface.IsEmpty() {
			t.Errorf("the interface %s should not be empty", actualInterface.Name())
		}

		if actualInterface.NumMethods() != len(expectedInterface.methods) {
			t.Errorf("the number of the methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.methods), actualInterface.NumMethods())
		}

		if actualInterface.NumExplicitMethods() != len(expectedInterface.explicitMethods) {
			t.Errorf("the number of the explicit methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.explicitMethods), actualInterface.NumExplicitMethods())
		}

		if actualInterface.NumEmbeddedTypes() != len(expectedInterface.embeddedTypes) {
			t.Errorf("the number of the embedded types of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.embeddedTypes), actualInterface.NumEmbeddedTypes())
		}

		assert.Equal(t, actualInterface, actualInterface.Underlying())

		assert.Equal(t, expectedInterface.position, actualInterface.Position(), "the position of the interface %s should be %w, but got %w",
			expectedInterfaceName, expectedInterface.position, actualInterface.Position())

		assertInterfaceEmbeddedTypes(t, fmt.Sprintf("interface %s", actualInterface.Name()), actualInterface.EmbeddedTypes(), expectedInterface.embeddedTypes)
		assertFunctions(t, fmt.Sprintf("interface %s", actualInterface.Name()), actualInterface.Methods(), expectedInterface.methods)
		assertFunctions(t, fmt.Sprintf("interface %s", actualInterface.Name()), actualInterface.ExplicitMethods(), expectedInterface.explicitMethods)
		assertMarkers(t, expectedInterface.markers, actualInterface.Markers(), fmt.Sprintf("interface %s", expectedInterfaceName))

		index++
	}

	return true
}

func assertInterfaceEmbeddedTypes(t *testing.T, interfaceName string, actualEmbeddedTypes *Types, expectedEmbeddedTypes []string) bool {

	if len(expectedEmbeddedTypes) != actualEmbeddedTypes.Len() {
		t.Errorf("the number of the embedded types should be %d, but got %d", len(expectedEmbeddedTypes), actualEmbeddedTypes.Len())
		return false
	}

	for index, expectedTypeName := range expectedEmbeddedTypes {
		actualEmbeddedType, ok := actualEmbeddedTypes.FindByName(expectedTypeName)

		if !ok {
			t.Errorf("embedded type with name %s for %s is not found", expectedTypeName, interfaceName)
			continue
		}

		if actualEmbeddedTypes.elements[index] != actualEmbeddedTypes.At(index) {
			t.Errorf("embedded type with name %s does not match with interface at index %d", actualEmbeddedType.Name(), index)
			continue
		}

		if actualEmbeddedType.Name() != expectedTypeName {
			t.Errorf("expected type name shoud be %s, but got %s", expectedTypeName, actualEmbeddedType.Name())
		}
	}

	return true
}
