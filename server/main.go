package main

import (
	//"net"
	//"context"
	//"log"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/dynamicpb"
	//"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

	//"google.golang.org/protobuf/reflect/protoreflect"

	//pref "google.golang.org/protobuf/reflect/protoreflect"
	//preg "google.golang.org/protobuf/reflect/protoregistry"
	apiPB "github.com/michaelhenkel/sub/api/proto"
	serverPB "github.com/michaelhenkel/sub/server/proto"
)

type server struct {
	serverPB.UnimplementedServerServer
}

//go:generate protoc -I../policy/proto -I$GOPATH/src/github.com/gogo/protobuf/gogoproto --proto_path=proto --gogo_out=plugins=grpc:proto --gogo_opt=paths=source_relative --include_source_info --include_imports --descriptor_set_out proto/genbyte/desc.protoset proto/server.proto

func newServer() *server {
	s := &server{}
	return s
}

type myFileDescriptor struct{}

func main() {
	fds := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(DSCByte, fds)
	if err != nil {
		fmt.Println(err)
	}

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	client := apiPB.NewApiClient(conn)

	for _, d := range fds.GetFile() {
		if *d.Name == "gogo.proto" {
			msg := &apiPB.Message{
				Fd: d,
			}
			_, err = client.RegisterType(context.Background(), msg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, d := range fds.GetFile() {
		if *d.Name == "server.proto" {
			msg := &apiPB.Message{
				Fd: d,
			}
			_, err = client.RegisterType(context.Background(), msg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	/*

		srByte, _ := sR.Descriptor()
		var fileSet descriptorpb.FileDescriptorSet
		if err := proto.Unmarshal(srByte, &fileSet); err != nil {
			fmt.Println(err)
		}
		fd, err := desc.CreateFileDescriptorFromSet(&fileSet)
		if err != nil {
			fmt.Println(err)
		}
	*/

	/*
		msgDesc, err := desc.LoadMessageDescriptorForMessage(serverPB.ServerRequest{})

		conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
		client := apiPB.NewApiClient(conn)

		var fileSet protoreflect.FileDescriptor
		fdProto := protodesc.ToFileDescriptorProto(serverPB.File_server_proto)

		msgDesc, err := desc.LoadMessageDescriptorForMessage(serverPB.ServerRequest{})

		msg := &apiPB.Message{
			Fd: fdProto,
		}
		_, err = client.RegisterType(context.Background(), msg)
		if err != nil {
			log.Println(err)
		}
	*/

}
