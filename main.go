package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	envBasePath = "INPUT_BASEPATH"
)

type report struct {
	Issues []result.Issue `json:"Issues"`
}

type annotation struct {
	file string
	line int
	col  int
	text string
}

func (a annotation) Output() string {
	return fmt.Sprintf("::error file=%s,line=%d,col=%d::%s\n", a.file, a.line, a.col, a.text)
}

type config struct {
	basePath string
}

func loadConfig() config {
	return config{
		basePath: os.Getenv(envBasePath),
	}
}

func createAnotations(cfg config, issues []result.Issue) []annotation {
	ann := make([]annotation, len(issues))
	for i := range issues {
		pos := issues[i].Pos
		file := filepath.Join(cfg.basePath, pos.Filename)
		ann[i] = annotation{
			file: file,
			line: pos.Line,
			col:  pos.Column,
			text: fmt.Sprintf("[%s] %s", issues[i].FromLinter, issues[i].Text),
		}
	}
	return ann
}

func reportFailures(cfg config, failures []result.Issue) {
	anns := createAnotations(cfg, failures)
	for _, ann := range anns {
		fmt.Println(ann.Output())
	}
}

func main() {
	cfg := loadConfig()

	var r report
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&r); err != nil {
		panic(err)
	}

	if len(r.Issues) > 0 {
		reportFailures(cfg, r.Issues)
	}

	os.Exit(0)
}
