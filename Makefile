
lint:
	golangci-lint run --config config/.golangci.yml

format:
	gofumpt -extra -l -w . && \
    gci write -s standard -s default -s blank .

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go install mvdan.cc/gofumpt@v0.5.0
	go install github.com/daixiang0/gci@v0.11.0
