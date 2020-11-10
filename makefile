GOPATH := $(shell go env GOPATH)
$(info GOPATH = $(GOPATH))
PATH := ${GOPATH}/bin:${CURDIR}/bin:$(PATH)
$(info PATH = $(PATH))

goexe = $(shell go env GOEXE)

.PHONY: build
build: test
	ls cmd | xargs -I {} go build -o bin/cmd/{} cmd/{}/main.go
	for file in ${STATIC_FILES}; do cp -R $$file bin/cmd/; done

.PHONY: test
test:
	go test -v ./...

.PHONY: codegen
codegen: bin/protoc-gen-gogoslick$(go_exe) bin/mockgen$(go_exe)
	protoc \
		-I=. \
		-I=$(GOPATH)/src \
		-I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf \
		--gogoslick_out=\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,\
plugins=grpc,\
paths=source_relative:. \
		./pkg/*.proto
	go generate ./...

bin/mockgen$(go_exe): go.mod
	go build -o $@ github.com/golang/mock/mockgen
	go get github.com/golang/mock/mockgen

bin/protoc-gen-gogoslick$(go_exe): go.mod
	go build -o $@ github.com/gogo/protobuf/protoc-gen-gogoslick
	go get github.com/gogo/protobuf/protoc-gen-gogoslick

