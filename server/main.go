package main

import (
	//"net"
	"context"
	"log"

	"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/reflect/protodesc"
	//"google.golang.org/protobuf/types/descriptorpb"

	//pref "google.golang.org/protobuf/reflect/protoreflect"
	//preg "google.golang.org/protobuf/reflect/protoregistry"
	apiPB "github.com/michaelhenkel/sub/api/proto"
	serverPB "github.com/michaelhenkel/sub/server/proto"
)

type server struct {
	serverPB.UnimplementedServerServer
}

//go:generate protoc -I../policy/proto --proto_path=proto --go_out=plugins=grpc:proto --go_opt=paths=source_relative proto/server.proto

func newServer() *server {
	s := &server{}
	return s
}

type myFileDescriptor struct{}

func main() {
	/*
		lis, err := net.Listen("tcp", "localhost:10001")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		s := newServer()
		serverPB.RegisterServerServer(grpcServer, s)
		grpcServer.Serve(lis)
		dynamicpb.NewMessage()
		//mt := dynamicpb.NewMessageType(message.ProtoReflect().Descriptor())
	*/

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	client := apiPB.NewApiClient(conn)
	fdProto := protodesc.ToFileDescriptorProto(serverPB.File_server_proto)
	msg := &apiPB.Message{
		Fd: fdProto,
	}
	_, err = client.RegisterType(context.Background(), msg)
	if err != nil {
		log.Println(err)
	}

}

//msgDesc := prMsg.Descriptor()
//dm := dynamicpb.NewMessageType(prMsg.Descriptor())
//msgDes := dm.Descriptor()
//protodesc.ToFileDescriptorProto()
//var fdProto *descriptorpb.FileDescriptorProto
//fdProto := fd.(*descriptorpb.FileDescriptorProto)
