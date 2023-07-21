package run

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/common/execute"
	"github.com/vision-cli/common/file"
	"github.com/vision-cli/common/tmpl"
)

const (
	goTemplateDir = "_templates"
	workflowDir   = ".github/workflows"
)

//go:embed all:_templates/tf
var templateFiles embed.FS

func Create(p *api_v1.PluginPlaceholders, executor execute.Executor, t tmpl.TmplWriter) error {
	var err error

	if err = tmpl.GenerateFS(templateFiles, goTemplateDir, "infra", p, false, t); err != nil {
		return fmt.Errorf("generating structure from the template: %w", err)
	}

	if err = genWorkflow(p); err != nil {
		return fmt.Errorf("generating service workflow with target dir: [%s]: %w", workflowDir, err)
	}

	return nil
}

//go:embed _templates/workflows/go.yml.tmpl
var goWorkflow string

func genWorkflow(p *api_v1.PluginPlaceholders) error {
	workflowName := "infra.yml"

	if err := Generate(goWorkflow, workflowDir, workflowName, p); err != nil {
		return fmt.Errorf("generating service workflow: %w", err)
	}

	return nil
}

// Generate writes template to filename in the targetDir, substituting placeholder values.
// Any existing files will be overwritten.
func Generate(template string, targetDir string, filename string, p any) error {
	t, err := tmpl.New(targetDir, template)
	if err != nil {
		return fmt.Errorf("parsing template file: %w", err)
	}

	if err = file.CreateDir(targetDir); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}
	newF, err := os.Create(filepath.Join(targetDir, filename))
	if err != nil {
		return fmt.Errorf("creating target file: %w", err)
	}
	defer newF.Close()

	if err = t.Execute(newF, p); err != nil {
		return fmt.Errorf("writing contents to %s: %w", filename, err)
	}

	return nil
}
