version = 0.0.1
binary_dir = dist
binary_name = master
coverage_file_name = coverage.out
main_package_path = ./cmd/master/main.go
ldflags = -X 'main.version=$(version)'
platforms := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64
tmp = $(subst /, , $@)
os = $(word 1, $(tmp))
arch = $(word 2, $(tmp))

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
	rm -f ${binary_name} ${coverage_file_name}
	rm -rf ${binary_dir}

.PHONY: test
test:
	go test -race -buildvcs ./...

.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=${coverage_file_name} ./...
	go tool cover -html=${coverage_file_name}

.PHONY: release
release: $(platforms)

.PHONY: $(platforms)
$(platforms):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags "$(ldflags)" -o '$(binary_dir)/$(binary_name)-$(os)-$(arch)' $(main_package_path)
