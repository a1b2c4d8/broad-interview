.PHONY: build
build:
	@mkdir -p _output
	@go build -o _output/main cmd/main.go

.PHONY: clean
clean:
	@rm -rf _output
