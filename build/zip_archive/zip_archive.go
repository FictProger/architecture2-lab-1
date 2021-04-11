package zip_archive

import (
	"fmt"
	"path"
	"strings"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/FictProger/architecture2-lab-1/build/zip_archive")

	// Ninja rule to execute zip.
	zipRule = pctx.StaticRule("zipArchive", blueprint.RuleParams{
		Command:     "zip $outputFile $files",
		Description: "zipping into $outputFile",
	}, "workDir", "outputFile", "files")
)

// zipArchiveType implements the simplest zipping.
type zipArchiveType struct {
	blueprint.SimpleName

	properties struct {
		// Archive name.
		Name string
		// List of input files.
		Srcs []string
	}
}

func (zipper *zipArchiveType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	outputPath := path.Join(config.BaseOutputDir, "archives")
	outputFile := path.Join(outputPath, zipper.properties.Name) + ".zip"

	var inputs []string

	for _, src := range zipper.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			inputs = append(inputs, matches...)
		}
	}
	for i, _ := range inputs {
		inputs[i] = "/" + inputs[i]
	}
	filesStr := strings.Join(inputs, " ")

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Archiving into '%s'", name),
		Rule:        zipRule,
		Outputs:     []string{outputPath},
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"outputFile": outputFile,
			"files":      filesStr,
		},
	})
}

// SimpleBinFactory is a factory for zip-archive module type which supports Go command packages without running tests.
func SimpleZipArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &zipArchiveType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
