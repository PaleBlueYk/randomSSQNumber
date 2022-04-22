.PHONY: all build

appName = "randSSQ"
buildDir = "bin"
configSourceFile = "config/config.toml"

all: build

build:
	@echo "building..."
	@rm -rf ${buildDir}
	@mkdir ${buildDir}
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${buildDir}/${appName} main.go
	@cp ${configSourceFile} ${buildDir}/config.toml
	@cp -r ./source/ ${buildDir}/source
	@echo "build complete and the bin file is located at <this project path>/bin"