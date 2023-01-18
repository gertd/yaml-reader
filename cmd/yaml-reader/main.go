package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "", "filepath to yaml input file")
	flag.Parse()

	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})

	if err := yaml.Unmarshal(buf, &data); err != nil {
		log.Fatal(err)
	}

	for k, v := range data {
		// echo "SELECTED_COLOR=green" >> $GITHUB_OUTPUT
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
}
