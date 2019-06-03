package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/robbiemcmichael/auth-mux/internal/config"
	"github.com/robbiemcmichael/auth-mux/internal/input"
	"github.com/robbiemcmichael/auth-mux/internal/output"
)

func handler(inputHandler input.HandlerFunc, outputHandler output.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := inputHandler(r)
		if err != nil {
			log.Printf("input handler: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		log.Printf("Authentication result: %+v", result)

		if err := outputHandler(w, result); err != nil {
			log.Printf("output handler: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	for _, i := range config.Inputs {
		for _, o := range config.Outputs {
			httpPath := path.Clean("/" + i.Path + "/" + o.Path)
			handler := handler(i.Config.Handler, o.Config.Handler)
			http.HandleFunc(httpPath, handler)
			log.Printf("Added handler: %s", httpPath)
		}
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
