SRC_DIR = .
PROTOS := $(shell find $(SRC_DIR) -name '*.proto')

all: proto

proto:
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--gofast_out=plugins=grpc:${SRC_DIR} \
		claw/motor/*.proto
	protoc -I/usr/local/include -I${SRC_DIR} \
		-I${GOPATH}/src \
		--gofast_out=plugins=grpc:${SRC_DIR} \
		clawserver/*.proto
