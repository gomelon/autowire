package autowire

import (
	"fmt"
	"github.com/gomelon/meta"
	"os"
	"testing"
	"text/template"
)

func TestTemplateGen(t *testing.T) {

	workdir, _ := os.Getwd()
	path := workdir + "/testdata/bar"
	generator, err := meta.NewTmplPkgGen(path, TmplWire, meta.WithOutputFilename("wire_set"),
		meta.WithFuncMapFactory(func(generator *meta.TmplPkgGen) template.FuncMap {
			return NewFunctions(generator).FuncMap()
		}),
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//err = generator.Print()
	err = generator.Generate()
	//err = generator.Generate()
	if err != nil {
		fmt.Println(err.Error())
	}
}
