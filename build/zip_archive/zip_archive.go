package zip_archive

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/FictProger/architecture2-lab-1/build/zip_archive")

	// Ninja rule to execute zip.
	zipRule = pctx.StaticRule("zipArchive", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outputPath",
		Description: "zipping into $outputPath",
	}, "workDir", "outputPath")
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
	outputPath := path.Join(config.BaseOutputDir, "archives", zipper.properties.Name)
    
	var inputs []string
	
	for _, src := range zipper.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			inputs = append(inputs, matches...)
		}
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Archiving into '%s'", name),
		Rule:        zipRule,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"outputPath": outputPath,
		},
	})
}

// SimpleBinFactory is a factory for zip-archive module type which supports Go command packages without running tests.
func SimpleZipArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &zipArchiveType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}