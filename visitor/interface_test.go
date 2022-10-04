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
		name:     "BakeryShop",
		fileName: "dessert.go",
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
		name:     "Dessert",
		fileName: "dessert.go",
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
					Name: "NewYearsEveCookie",
				},
			},
		},
		name:     "NewYearsEveCookie",
		fileName: "dessert.go",
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
		name:     "SweetShop",
		fileName: "dessert.go",
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
		embeddedTypes: []string{"NewYearsEveCookie", "Dessert"},
	}
)

func assertInterfaces(t *testing.T, file *File, interfaces map[string]interfaceInfo) bool {

	if len(interfaces) != file.Interfaces().Len() {
		t.Errorf("the number of the interface should be %d, but got %d", len(interfaces), file.Interfaces().Len())
		return false
	}

	for expectedInterfaceName, expectedInterface := range interfaces {
		actualInterface, ok := file.Interfaces().FindByName(expectedInterfaceName)

		if !ok {
			t.Errorf("interface with name %s is not found", expectedInterfaceName)
			continue
		}

		if expectedInterface.fileName != actualInterface.File().Name() {
			t.Errorf("the file name for interface %s should be %s, but got %s", expectedInterfaceName, expectedInterface.fileName, actualInterface.File().Name())
		}

		if actualInterface.NumMethods() == 0 && !actualInterface.IsEmpty() {
			t.Errorf("the interface %s should be empty", actualInterface.Name())
		} else if actualInterface.NumMethods() != 0 && actualInterface.IsEmpty() {
			t.Errorf("the interface %s should not be empty", actualInterface.Name())
		}

		if actualInterface.NumMethods() != len(expectedInterface.methods) {
			t.Errorf("the number of the methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.methods), actualInterface.NumMethods())
			continue
		}

		if actualInterface.NumExplicitMethods() != len(expectedInterface.explicitMethods) {
			t.Errorf("the number of the explicit methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.explicitMethods), actualInterface.NumExplicitMethods())
			continue
		}

		if actualInterface.NumEmbeddedTypes() != len(expectedInterface.embeddedTypes) {
			t.Errorf("the number of the embedded types of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.embeddedTypes), actualInterface.NumEmbeddedTypes())
			continue
		}

		for index, expectedEmbeddedType := range expectedInterface.embeddedTypes {
			actualEmbeddedType := actualInterface.EmbeddedTypes()[index]
			if expectedEmbeddedType != actualEmbeddedType.Name() {
				t.Errorf("the interface %s should have the embedded type %s at index %d, but got %s", expectedInterfaceName, expectedEmbeddedType, index, actualEmbeddedType.Name())
				continue
			}
		}

		assert.Equal(t, actualInterface, actualInterface.Underlying())

		assert.Equal(t, expectedInterface.position, actualInterface.Position(), "the position of the interface %s should be %w, but got %w",
			expectedInterfaceName, expectedInterface.position, actualInterface.Position())

		assertFunctions(t, fmt.Sprintf("interface %s", actualInterface.Name()), actualInterface.Methods(), expectedInterface.methods)
		assertMarkers(t, expectedInterface.markers, actualInterface.Markers(), fmt.Sprintf("interface %s", expectedInterfaceName))

	}

	return true
}
