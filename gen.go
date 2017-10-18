//go:generate protoc -I protocol/ protocol/protocol.proto --go_out=plugins=grpc:protocol
package accord
