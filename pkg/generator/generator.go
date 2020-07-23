package generator

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/maildealru/iniziamo/pkg/parser"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) WriteFile(info *parser.Info) {
	baseName := info.Name[:len(info.Name)-len(".go")]
	iniziamoName := filepath.Join(info.Dir, baseName+"_iniziamo.go")

	f, err := os.Create(iniziamoName)
	if err != nil {
		log.Fatalf("failed to create file: %s", iniziamoName)
	}

	g.executeTemplate(f, info)
	if err := exec.Command("go", "fmt", iniziamoName).Run(); err != nil {
		_ = os.Remove(iniziamoName)
		log.Fatalf("failed to format file: %s", err.Error())
	}
	if err := exec.Command("goimports", "-w", iniziamoName).Run(); err != nil {
		_ = os.Remove(iniziamoName)
		log.Fatalf("failed to format file imports: %s", err.Error())
	}
}

func (g *Generator) executeTemplate(f *os.File, info *parser.Info) {
	t, err := template.New("iniziamo").Parse(IniziamoTemplate)
	if err != nil {
		panic(err)
	}
	if err := t.Execute(f, *info); err != nil {
		panic(err)
	}
}
