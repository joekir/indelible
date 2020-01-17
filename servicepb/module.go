//go:generate protoc -I/usr/local/include -I. -I$GOPATH/src --go_out=plugins=grpc:. indelible.proto

package servicepb
