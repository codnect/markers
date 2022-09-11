package processor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var placeHolderRegex = regexp.MustCompile(`{(.*?)}`)

type Context struct {
	dirs           []string
	loadResult     *packages.LoadResult
	registry       *marker.Registry
	collector      *marker.Collector
	config         Config
	configFilePath string
	goModuleDir    string
	packageId      string
	version        string
	errors         []error
	values         map[string]any
	args           []string
}

func (ctx *Context) Directories() []string {
	return ctx.dirs
}

func (ctx *Context) Registry() *marker.Registry {
	return ctx.registry
}

func (ctx *Context) ModuleRoot() string {
	return ctx.goModuleDir
}

func (ctx *Context) OutputPath() string {
	outputPath, _ := ctx.ParameterValue("OUTPUT_PATH")
	return outputPath.Value
}

func (ctx *Context) Value(name string) (any, bool) {
	if value, exits := ctx.values[name]; exits {
		return value, true
	}

	return nil, false
}

func (ctx *Context) Set(name string, value any) {
	ctx.values[name] = value
}

func (ctx *Context) Args() []string {
	return ctx.args
}

func (ctx *Context) Parameters() []Parameter {
	parametersMap := make(map[string]Parameter, 0)
	parameters := make([]Parameter, 0)

	for _, parameter := range ctx.config.Parameters {
		parametersMap[parameter.Name] = Parameter{
			Name:  parameter.Name,
			Value: ctx.resolveParameters(parameter.Name, parameter.Value),
		}
	}

	overrides := ctx.findOverrides()

	for _, parameter := range overrides {
		parametersMap[parameter.Name] = Parameter{
			Name:  parameter.Name,
			Value: ctx.resolveParameters(parameter.Name, parameter.Value),
		}
	}

	for _, value := range parametersMap {
		parameters = append(parameters, value)
	}

	return parameters
}

func (ctx *Context) ParameterValue(name string) (Parameter, bool) {
	if name == "MODULE_ROOT" {
		return Parameter{
			"MODULE_ROOT",
			ctx.goModuleDir,
		}, true
	}

	param, exists := ctx.findParameterInGlobal(name)
	overrideParam, overrideExists := ctx.findParameterInOverrides(name)

	if !exists || !overrideExists {
		if name == "OUTPUT_PATH" {
			return Parameter{
				"OUTPUT_PATH",
				fmt.Sprintf("%s/generated", ctx.goModuleDir),
			}, true
		}

		return Parameter{}, false
	}

	if overrideExists {
		return *overrideParam, true
	}

	if exists {
		return *param, true
	}

	return Parameter{}, false
}

func (ctx *Context) findParameterInGlobal(name string) (*Parameter, bool) {
	if len(ctx.config.Parameters) == 0 {
		return nil, false
	}

	for _, parameter := range ctx.config.Parameters {
		if parameter.Name == name {
			return &Parameter{
				parameter.Name,
				ctx.resolveParameters(parameter.Name, parameter.Value),
			}, true
		}
	}

	return nil, false
}

func (ctx *Context) findOverrides() []Parameter {
	for _, override := range ctx.config.Overrides {
		if override.Package == ctx.packageId && override.Version == ctx.version {
			return override.Parameters
		}
	}

	return nil
}

func (ctx *Context) findParameterInOverrides(name string) (*Parameter, bool) {
	if len(ctx.config.Overrides) == 0 {
		return nil, false
	}

	overrideParameters := ctx.findOverrides()

	for _, parameter := range overrideParameters {
		if parameter.Name == name {
			return &Parameter{
				parameter.Name,
				ctx.resolveParameters(parameter.Name, parameter.Value),
			}, true
		}
	}

	return nil, false
}

func (ctx *Context) resolveParameters(name, value string) string {
	matches := placeHolderRegex.FindAllStringSubmatch(value, -1)
	placeHolders := make([]string, 0)
	for _, match := range matches {
		if len(match) == 2 {
			placeHolders = append(placeHolders, match[1])
		}
	}

	for _, placeHolder := range placeHolders {
		parameter, exits := ctx.ParameterValue(placeHolder)
		if !exits {
			continue
		}

		value = strings.ReplaceAll(value, fmt.Sprintf("${%s}", parameter.Name), parameter.Value)
	}

	if name == "OUTPUT_PATH" {
		return filepath.FromSlash(value)
	}

	return value
}

func (ctx *Context) Config() Config {
	return ctx.config
}

func (ctx *Context) ConfigFilePath() string {
	return ctx.configFilePath
}

func (ctx *Context) Collector() *marker.Collector {
	return ctx.collector
}

func (ctx *Context) LoadResult() *packages.LoadResult {
	return ctx.loadResult
}

func (ctx *Context) Error(err error) {
	ctx.errors = append(ctx.errors, err)
}

func (ctx *Context) printError(err error) {
	if err != nil {

		switch typedErr := err.(type) {
		case marker.ErrorList:
			ctx.printErrors(typedErr)
			return
		}

		log.Println(err)
		return
	}
}

func (ctx *Context) printErrors(errorList marker.ErrorList) {
	if errorList == nil || len(errorList) == 0 {
		return
	}

	for _, err := range errorList {
		switch typedErr := err.(type) {
		case marker.Error:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ParserError:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ErrorList:
			ctx.printErrors(typedErr)
		default:
			ctx.printError(err)
		}
	}
}
