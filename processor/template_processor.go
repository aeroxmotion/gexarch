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

func collectTemplateFilenames(dir string) (filenames []string) {
	files, err := templateResources.ReadDir(dir)
	util.PanicIfError(err)

	for _, file := range files {
		curDir := path.Join(dir, file.Name())

		if file.IsDir() {
			filenames = append(filenames, collectTemplateFilenames(curDir)...)
		} else {
			filenames = append(filenames, curDir)
		}
	}

	return
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

func (processor *TemplateProcessor) processFile(rootDir, filename string) {
	templateFileContent, err := templateResources.ReadFile(filename)
	util.PanicIfError(err)

	contentTpl, err := processor.fileContentTemplate.New("").Parse(string(templateFileContent))
	util.PanicIfError(err)

	nameTpl, err := processor.filenameTemplate.New("").Parse(filename)
	util.PanicIfError(err)

	// Process filename
	filenameBuffer := bytes.NewBufferString("")
	nameTpl.Execute(filenameBuffer, processor.config.ToSnakeValues())

	targetFilename := strings.Replace(
		filenameBuffer.String(),
		rootDir,
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

func (processor *TemplateProcessor) processFiles(suffixDir string) {
	rootDir := path.Join(templateRootDirname, suffixDir)
	filenames := collectTemplateFilenames(rootDir)

	for _, filename := range filenames {
		processor.processFile(rootDir, filename)
	}
}

func (processor *TemplateProcessor) ProcessByType() {
	processor.processFiles("type")
}
