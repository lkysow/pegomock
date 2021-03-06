package filehandling

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/petergtz/pegomock/mockgen"
	"github.com/petergtz/pegomock/model"
	"github.com/petergtz/pegomock/modelgen/gomock"
	"github.com/petergtz/pegomock/modelgen/loader"
	"github.com/petergtz/pegomock/pegomock/util"
)

func GenerateMockFileInOutputDir(
	args []string,
	outputDirPath string,
	outputFilePathOverride string,
	packageOut string,
	selfPackage string,
	debugParser bool,
	out io.Writer,
	useExperimentalModelGen bool) {
	GenerateMockFile(
		args,
		OutputFilePath(args, outputDirPath, outputFilePathOverride),
		packageOut,
		selfPackage,
		debugParser,
		out,
		useExperimentalModelGen)
}

func OutputFilePath(args []string, outputDirPath string, outputFilePathOverride string) string {
	if outputFilePathOverride != "" {
		return outputFilePathOverride
	} else if util.SourceMode(args) {
		return filepath.Join(outputDirPath, "mock_"+strings.TrimSuffix(args[0], ".go")+"_test.go")
	} else {
		return filepath.Join(outputDirPath, "mock_"+strings.ToLower(args[len(args)-1])+"_test.go")
	}
}

func GenerateMockFile(args []string, outputFilePath string, packageOut string, selfPackage string, debugParser bool, out io.Writer, useExperimentalModelGen bool) {
	output := GenerateMockSourceCode(args, packageOut, selfPackage, debugParser, out, useExperimentalModelGen)

	err := ioutil.WriteFile(outputFilePath, output, 0664)
	if err != nil {
		panic(fmt.Errorf("Failed writing to destination: %v", err))
	}
}

func GenerateMockSourceCode(args []string, packageOut string, selfPackage string, debugParser bool, out io.Writer, useExperimentalModelGen bool) []byte {
	var err error

	var ast *model.Package
	var src string
	if util.SourceMode(args) {
		ast, err = gomock.ParseFile(args[0])
		src = args[0]
	} else {
		if len(args) != 2 {
			log.Fatal("Expected exactly two arguments, but got " + fmt.Sprint(args))
		}
		if useExperimentalModelGen {
			ast, err = loader.GenerateModel(args[0], args[1])

		} else {
			ast, err = gomock.Reflect(args[0], strings.Split(args[1], ","))
		}
		src = fmt.Sprintf("%v (interfaces: %v)", args[0], args[1])
	}
	if err != nil {
		panic(fmt.Errorf("Loading input failed: %v", err))
	}

	if debugParser {
		ast.Print(out)
	}

	output, err := mockgen.GenerateOutput(ast, src, packageOut, selfPackage)
	if err != nil {
		panic(fmt.Errorf("Failed generating mock: %v", err))
	}
	return output
}
