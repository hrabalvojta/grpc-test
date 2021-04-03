generate:
	buf generate

run: generate
	go run main.go

build: generate
	go build -o "bin/grpc-server" main.go

test: 
	go test

lint:
	buf lint
	buf breaking --against 'https://github.com/johanbrandhorst/grpc-gateway-boilerplate.git#branch=master'

dbuild: generate
	docker build -t vhrabal/grpc-test:latest .
	docker image prune -f

drun: dbuild
	docker run --rm --name grpc_test -p10000:10000 -p11000:11000 -d vhrabal/grpc-test:latest
	
BUF_VERSION:=0.41.0

install:
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
  	chmod +x "$(shell go env GOPATH)/bin/buf"
