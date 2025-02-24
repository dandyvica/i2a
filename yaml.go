package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// This will match the YAML configuration file where all settings are defined
type YAMLConfig struct {
	Perf struct {
		Channels struct {
			Jobs int64
			Sql  int64
		}
		Workers int64
	}
	Hash struct {
		Maxsize string
	}
}

// Read the YAML configuration file
func ReadYAMLConfig(configFile string) (*YAMLConfig, error) {
	// read whole file in memory
	yamlData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error <%v> opening YAML configuration file: <%s>", err, configFile)
		return nil, err
	}

	// read YAML into struct
	yamlConf := &YAMLConfig{}
	fmt.Printf("%+v\n", yamlConf)
	err = yaml.Unmarshal(yamlData, yamlConf)
	if err != nil {
		log.Fatalf("error <%v> reading YAML configuration file: <%s>", err, configFile)
		return nil, err
	}

	return yamlConf, nil
}
