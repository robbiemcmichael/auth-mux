.PHONY: build
build: auth-mux

auth-mux:
	go build -o $@ ./cmd/$@

.PHONY: docker
docker:
	docker build .

.PHONY: clean
clean:
	rm -f auth-mux
