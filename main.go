package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	filePath string
	module   string
)

type ResourceInstanceAttributes struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ResourceInstance struct {
	SchemaVersion int                        `json:"schema_version"`
	Attributes    ResourceInstanceAttributes `json:"attributes,omitempty"`
}

type Resource struct {
	Module    string             `json:"module,omitempty"`
	Mode      string             `json:"mode"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
}

type TerraformState struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

func main() {
	flag.StringVar(&filePath, "file", "state.json", "path to state file")
	flag.StringVar(&module, "module", "all", "name of module to filter resources by")
	flag.Parse()

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	var state TerraformState
	json.Unmarshal(b, &state)

	for r := 0; r < len(state.Resources); r++ {
		// If we're only looking at a particular resource, skip
		// all other resources that aren't in the selected module
		if module != "all" && state.Resources[r].Module != module {
			continue
		}

		// Skip when state item is not a resource (mode=managed)
		if state.Resources[r].Mode != "managed" {
			continue
		}

		// Deal with resources that have a `count` set
		if len(state.Resources[r].Instances) >= 2 {
			for i := 0; i < len(state.Resources[r].Instances); i++ {
				fmt.Printf(
					"terraform import \"%s.%s.%s[%s]\" \"%s\"\n",
					state.Resources[r].Module,
					state.Resources[r].Type,
					state.Resources[r].Name,
					strconv.Itoa(i),
					state.Resources[r].Instances[i].Attributes.ID,
				)
			}
    } else if len(state.Resources[r].Instances) == 1 {
      fmt.Printf(
        "terraform import \"%s.%s.%s\" \"%s\"\n",
        state.Resources[r].Module,
        state.Resources[r].Type,
        state.Resources[r].Name,
        state.Resources[r].Instances[int(0)].Attributes.ID,
      )
		} else if len(state.Resources[r].Instances) == 0 {
      continue
    }
	}
}
