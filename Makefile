.PHONY: build
build: auth-mux

auth-mux:
	go build -o $@ ./cmd/$@

.PHONY: clean
clean:
	rm -f auth-mux
