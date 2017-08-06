package template

import (
	"bitbucket.org/pkg/inflect"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type state struct {
	n int
}

func (s *state) Set(n int) int {
	s.n = n
	return n
}

func (s *state) Inc() int {
	s.n++
	return s.n
}

var s state
var (
	funcMap = template.FuncMap{
		"Pluralize":  inflect.Pluralize,
		"Underscore": inflect.Underscore,
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
	}
)

type Builder struct {
	TemplateName string
	TemplatePath string
}

func NewBuilder(templatePath string) *Builder {
	if templatePath[0:1] != "/" {
		templatePath = TemplatePath(templatePath)
	}

	templateName := filepath.Base(templatePath)
	builder := &Builder{
		TemplateName: templateName,
		TemplatePath: templatePath,
	}

	return builder
}

func (builder *Builder) Template() *template.Template {
	contents := LoadTemplateFromFile(builder.TemplatePath)
	tmpl := template.Must(template.New(builder.TemplateName).Funcs(funcMap).Parse(contents))

	return tmpl
}

func (builder *Builder) Write(writer io.Writer, data interface{}) {
	tmpl := builder.Template()
	err := tmpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

func (builder *Builder) WriteToPath(outputPath string, data interface{}) {
	printAction("green+h:black", "create", outputPath)
	if _, err := os.Stat(outputPath); err == nil {
		printAction("red+h:black", "skip", outputPath)
		return
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	builder.Write(file, data)
}

func (builder *Builder) InsertAfterToPath(outputPath string, after string, data interface{}) {
	printAction("cyan+h:black", "insert", outputPath)

	newFilePath := outputPath + ".new"

	file, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	outputFile, err := os.Create(newFilePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()

		writer.WriteString(line + "\n")
		if strings.HasPrefix(line, after) {
			builder.Write(writer, data)
		}
	}

	writer.Flush()
	outputFile.Close()

	os.Rename(newFilePath, outputPath)
}
