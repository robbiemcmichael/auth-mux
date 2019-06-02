package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"github.com/robbiemcmichael/auth-mux/internal/config"
)

func main() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", config)
	fmt.Printf("%+v\n", config.Inputs[0].Config)
	fmt.Printf("%+v\n", config.Outputs[0].Config)
}
