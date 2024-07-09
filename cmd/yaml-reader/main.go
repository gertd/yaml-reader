package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const (
	resultRoot = "data"
)

type parameters struct {
	filename   string
	jsonOutput bool
}

func main() {
	p := parameters{}

	flag.StringVar(&p.filename, "file", "", "filepath to yaml input file")
	flag.BoolVar(&p.jsonOutput, "json", false, "output as JSON variable")
	flag.Parse()

	if ok, err := p.validate(); !ok {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err := exec(&p); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}

func (p *parameters) validate() (bool, error) {
	if p.filename == "" {
		return false, errors.Errorf("--file not specified")
	}

	if !fileExist(p.filename) {
		return false, errors.Errorf("--file [%s] does not exist", p.filename)
	}

	return true, nil
}

func exec(p *parameters) error {
	data, err := readFile(p.filename)
	if err != nil {
		return err
	}

	writer, err := openGitHubOutput()
	if err != nil {
		return err
	}

	defer func() {
		_ = writer.Sync()
		_ = writer.Close()
	}()

	if p.jsonOutput {
		return outputJSON(writer, data)
	}

	return outputVariables(writer, data)
}

func readFile(filename string) (map[string]interface{}, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})

	if err := yaml.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func outputJSON(w *os.File, data map[string]interface{}) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "")

	if err := enc.Encode(data); err != nil {
		return err
	}

	return write(w, resultRoot, buf.String())
}

func outputVariables(w *os.File, data map[string]interface{}) error {
	for k, v := range data {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if err := write(w, k, s); err != nil {
			return err
		}
	}
	return nil
}

func write(w *os.File, key, value string) error {
	_, err := fmt.Fprintf(w, "%s=%s\n", key, strings.TrimSpace(value))
	return err
}

func openGitHubOutput() (*os.File, error) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return os.Stdout, nil
	}

	filename := os.Getenv("GITHUB_OUTPUT")
	if filename == "" {
		return nil, os.ErrNotExist
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

func fileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
