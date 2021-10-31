package processor

import (
	"bytes"
	"embed"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/util"
	"github.com/iancoleman/strcase"
)

var (
	templateRootDirname = "__template__"
	//go:embed __template__
	templateResources embed.FS
)

type TemplateProcessor struct {
	config              *config.ProcessorConfig
	fileContentTemplate *template.Template
	filenameTemplate    *template.Template
}

func NewTemplateProcessor(processorConfig *config.ProcessorConfig) *TemplateProcessor {
	processor := &TemplateProcessor{
		config: processorConfig,
	}

	processor.initTemplates()
	return processor
}

func (processor *TemplateProcessor) initTemplates() {
	fileContentTemplate := template.New("file_content")
	filenameTemplate := template.New("filename")

	// Register custom pipelines
	fileContentTemplate.Funcs(template.FuncMap{
		"to_snake": strcase.ToSnake,
	})

	// Define custom delimiters
	filenameTemplate.Delims("[", "]")

	processor.filenameTemplate = filenameTemplate
	processor.fileContentTemplate = fileContentTemplate
}

func collectTemplateFilenames(dir string) (filenames []string) {
	files, err := templateResources.ReadDir(dir)
	util.PanicIfError(err)

	for _, file := range files {
		if file.IsDir() {
			filenames = append(filenames, collectTemplateFilenames(path.Join(dir, file.Name()))...)
		} else {
			filenames = append(filenames, path.Join(dir, file.Name()))
		}
	}

	return
}

func (processor *TemplateProcessor) processFile(filename string) {
	templateFileContent, err := templateResources.ReadFile(filename)
	util.PanicIfError(err)

	contentTpl, err := processor.fileContentTemplate.New("").Parse(string(templateFileContent))
	util.PanicIfError(err)

	nameTpl, err := processor.filenameTemplate.New("").Parse(filename)
	util.PanicIfError(err)

	// Process filename
	filenameBuffer := bytes.NewBufferString("")
	nameTpl.Execute(filenameBuffer, config.ProcessorConfig{
		TypeName:       strcase.ToSnake(processor.config.TypeName),
		EntityName:     strcase.ToSnake(processor.config.EntityName),
		RepositoryName: strcase.ToSnake(processor.config.RepositoryName),
		UseCaseName:    strcase.ToSnake(processor.config.UseCaseName),
	})

	targetFilename := strings.Replace(
		filenameBuffer.String(),
		templateRootDirname,
		path.Join(processor.config.TargetPath, strcase.ToSnake(processor.config.TypeName)),
		1,
	)
	targetFilename = strings.Replace(targetFilename, "gotpl", "go", 1)

	// Create directory
	dirname := path.Dir(targetFilename)

	util.PanicIfError(os.MkdirAll(dirname, os.ModePerm))

	// Create file
	file, err := os.Create(targetFilename)
	util.PanicIfError(err)

	// Generate file
	util.PanicIfError(contentTpl.Execute(file, processor.config))
}

func (processor *TemplateProcessor) Process() {
	filenames := collectTemplateFilenames(templateRootDirname)

	for _, filename := range filenames {
		processor.processFile(filename)
	}
}
