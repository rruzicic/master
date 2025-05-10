binary_name = master
coverage_file_name = coverage.out
main_package_path = ./cmd/master/main.go

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: build
build:
	go build -o ${binary_name} ${main_package_path}

.PHONY: run
run: build
	./${binary_name}

.PHONY: clean
clean:
	rm ${binary_name} ${coverage_file_name}

.PHONY: test
test:
	go test -race -buildvcs ./...

.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=${coverage_file_name} ./...
	go tool cover -html=${coverage_file_name}