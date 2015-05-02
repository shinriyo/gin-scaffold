package command

import (
	//"fmt"
	"bitbucket.org/pkg/inflect"
	"github.com/dcu/gin-scaffold/template"
	"path/filepath"
	//"strings"
)

type ControllerCommand struct {
	PackageName        string
	ControllerName     string
	ModelName          string
	InstanceName       string
	InstanceNamePlural string
	RoutePath          string
	TemplateName       string
	Fields             map[string]string
}

func (command *ControllerCommand) Execute(args []string) {
	command.ControllerName = args[0]
	command.RoutePath = inflect.Underscore(command.ControllerName)
	command.ModelName = inflect.Singularize(command.ControllerName)

	command.InstanceName = inflect.CamelizeDownFirst(command.ControllerName)
	command.InstanceNamePlural = inflect.Pluralize(command.InstanceName)
	command.PackageName = template.PackageName()

	outputPath := filepath.Join("controllers", inflect.Underscore(command.ControllerName)+".go")

	builder := template.NewBuilder("controller.go.tmpl")
	builder.Write(outputPath, command)
}