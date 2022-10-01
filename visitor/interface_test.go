package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"testing"
)

type interfaceInfo struct {
	markers         marker.MarkerValues
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

		assertFunctions(t, fmt.Sprintf("interface %s", actualInterface.Name()), actualInterface.Methods(), expectedInterface.methods)
		assertMarkers(t, expectedInterface.markers, actualInterface.Markers(), fmt.Sprintf("interface %s", expectedInterfaceName))

	}

	return true
}
