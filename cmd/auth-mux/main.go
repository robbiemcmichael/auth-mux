package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/robbiemcmichael/auth-mux/internal"
)

func handler(i internal.Input, o internal.Output) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validation, err := i.Config.Handler(r)
		if err != nil {
			log.Printf("input handler for %q: %v", i.Name, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		log.Printf("Authentication successful: %+v", validation)

		if err := o.Config.Handler(w, validation); err != nil {
			log.Printf("output handler for %q: %v", o.Name, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	configPath := flag.String("c", "config.yaml", "path to the config file")
	flag.Parse()

	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config internal.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	for _, i := range config.Inputs {
		for _, o := range config.Outputs {
			httpPath := path.Clean("/" + i.Path + "/" + o.Path)
			handler := handler(i, o)
			http.HandleFunc(httpPath, handler)
			log.Printf("Added handler: %s", httpPath)
		}
	}

	bind := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Fatal(http.ListenAndServeTLS(bind, config.Cert, config.Key, nil))
}
