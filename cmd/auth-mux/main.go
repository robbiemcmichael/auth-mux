package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/robbiemcmichael/auth-mux/internal/config"
)

func handler(i config.Input, o config.Output) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := i.Config.Handler(r)
		if err != nil {
			log.Printf("input handler for %q: %v", i.Name, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		log.Printf("Authentication result: %+v", result)

		if err := o.Config.Handler(w, result); err != nil {
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

	var config config.Config

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
