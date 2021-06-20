package marker

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type ArgumentType int

const (
	InvalidType ArgumentType = iota
	RawType
	AnyType
	BoolType
	IntType
	StringType
	SliceType
	MapType
)

var (
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	rawArgsType   = reflect.TypeOf((*[]byte)(nil)).Elem()
)

type Argument struct {
	Name     string
	Type     ArgumentType
	Pointer  bool
	Required bool

	ItemType *ArgumentType
}

type Definition struct {
	Name       string
	Level      TargetLevel
	Output     reflect.Type
	Fields     map[string]Argument
	FieldNames map[string]string
}

func (definition *Definition) extract() error {

	for index := 0; index < definition.Output.NumField(); index++ {
		field := definition.Output.Field(index)

		if field.PkgPath != "" {
			continue
		}

		fieldInfo, err := getArgumentInfo(field)

		if err != nil {
			return err
		}

		definition.Fields[fieldInfo.Name] = fieldInfo
		definition.FieldNames[fieldInfo.Name] = field.Name
	}

	return nil
}

func getArgumentInfo(structField reflect.StructField) (Argument, error) {
	fieldName := lowerCamelCase(structField.Name)

	markerTag, tagExists := structField.Tag.Lookup("marker")
	markerTagValues := strings.Split(markerTag, ",")

	if tagExists && markerTagValues[0] != "" {
		fieldName = markerTagValues[0]
	}

	optionalOption := false

	for _, tagOption := range markerTagValues[1:] {

		if tagOption == "optional" {
			optionalOption = true
		}

	}

	fieldType := structField.Type
	argumentType, err := getArgumentType(fieldType)

	if err != nil {
		return Argument{}, err
	}

	var argumentItemType ArgumentType

	if argumentType == SliceType || argumentType == MapType {
		itemType, err := getArgumentType(fieldType.Elem())

		if err != nil && argumentType == SliceType {
			return Argument{}, fmt.Errorf("bad slice item type: %w", err)
		} else if err != nil && argumentType == MapType {
			return Argument{}, fmt.Errorf("bad map item type: %w", err)
		}

		argumentItemType = itemType
	}

	isPointer := false
	isOptional := false

	if fieldType.Kind() == reflect.Ptr {
		isPointer = true
		isOptional = true
	}

	optionalOption = optionalOption || isOptional

	return Argument{
		Name:     fieldName,
		Type:     argumentType,
		Pointer:  isPointer,
		Required: !optionalOption,
		ItemType: &argumentItemType,
	}, nil
}

func getArgumentType(rawType reflect.Type) (ArgumentType, error) {

	if rawType == rawArgsType {
		return RawType, nil
	}

	if rawType == interfaceType {
		return AnyType, nil
	}

	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	switch rawType.Kind() {
	case reflect.String:
		return StringType, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return IntType, nil
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		return IntType, nil
	case reflect.Bool:
		return BoolType, nil
	case reflect.Slice:
		return SliceType, nil
	case reflect.Map:

		if rawType.Key().Kind() != reflect.String {
			return InvalidType, fmt.Errorf("bad map key type: map key must be string")
		}

		return MapType, nil
	default:
		return InvalidType, fmt.Errorf("type has unsupported kind %s", rawType.Kind())
	}
}

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}

func MakeDefinition(name string, level TargetLevel, output interface{}) (*Definition, error) {
	outputType := reflect.TypeOf(output)

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	definition := &Definition{
		Name:       name,
		Level:      level,
		Output:     outputType,
		FieldNames: make(map[string]string),
	}

	err := definition.extract()

	if err != nil {
		return nil, err
	}

	return definition, nil
}

type Registry struct {
	definitionMap map[string]*Definition

	initOnce sync.Once
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (registry *Registry) initialize() {
	registry.initOnce.Do(func() {

		if registry.definitionMap == nil {
			registry.definitionMap = make(map[string]*Definition)
		}

	})
}

func (registry *Registry) Register(name string, level TargetLevel, output interface{}) error {
	registry.initialize()

	def, err := MakeDefinition(name, level, output)

	if err != nil {
		return err
	}

	return registry.RegisterWithDefinition(def)
}

func (registry *Registry) RegisterWithDefinition(definition *Definition) error {
	registry.initialize()

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if definition.Level == 0 {
		return fmt.Errorf("specify target levels for the definition : %v", definition.Name)
	}

	if _, ok := registry.definitionMap[definition.Name]; ok {
		return fmt.Errorf("there is already registered definition : %v", definition.Name)
	}

	registry.definitionMap[definition.Name] = definition

	return nil
}

func (registry *Registry) Lookup(marker string) *Definition {
	registry.initialize()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	name, anonymousName, _ := splitMarker(marker)

	if def, exists := registry.definitionMap[anonymousName]; exists {
		return def
	}

	return registry.definitionMap[name]
}
