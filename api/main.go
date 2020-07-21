package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/jsonpb"
	apiPB "github.com/michaelhenkel/sub/api/proto"
	policyPB "github.com/michaelhenkel/sub/policy/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type api struct {
	apiPB.UnimplementedApiServer
}

var GlobalTypes *protoregistry.Types = new(protoregistry.Types)

//go:generate protoc -I../policy/proto --proto_path=proto --go_out=plugins=grpc:proto --go_opt=paths=source_relative proto/api.proto

func newAPI() *api {
	s := &api{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", "localhost:10002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	a := newAPI()
	apiPB.RegisterApiServer(grpcServer, a)
	grpcServer.Serve(lis)
}

func (a *api) Create(ctx context.Context, apiMsg *apiPB.Request) (*apiPB.Result, error) {
	msg, err := dynamicMsgFromRequest(apiMsg)
	if err != nil {
		fmt.Println(err)
	}
	rangeOverDynamicMsg(msg)
	m := jsonpb.Marshaler{
		Indent: "  ",
	}
	stringMsg, err := m.MarshalToString(msg)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(stringMsg)
	return &apiPB.Result{}, nil
}

func rangeOverMsg(value protoreflect.Value, resName protoreflect.Name, parentName protoreflect.Name) {
	resNameString := string(resName)
	dtypeFD := getDtypeFD(value.Message().Descriptor())
	m := value.Message().New().Mutable(dtypeFD)
	dtypeList := m.List()
	dtypeList.Append(protoreflect.ValueOf(resNameString))
	value.Message().Set(dtypeFD, protoreflect.ValueOfList(dtypeList))
	value.Message().Range(func(descriptor pref.FieldDescriptor, value pref.Value) bool {
		return processMsg(descriptor.Message(), descriptor, value, resName, parentName)
	})
}
func rangeOverDynamicMsg(msg *dynamicpb.Message) {
	resName := msg.ProtoReflect().Descriptor().Name()
	resNameString := string(resName)
	dtypeFD := getDtypeFD(msg.Descriptor())
	m := msg.New().Mutable(dtypeFD)
	dtypeList := m.List()
	dtypeList.Append(protoreflect.ValueOf(resNameString))
	msg.Set(dtypeFD, protoreflect.ValueOfList(dtypeList))
	parentName := resName
	msg.Range(func(descriptor pref.FieldDescriptor, value pref.Value) bool {
		return processMsg(descriptor.Message(), descriptor, value, resName, parentName)
	})
}

func getDtypeFD(d protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	dtypeFieldName := protoreflect.Name("dtype")
	f := d.Fields()
	return f.ByName(dtypeFieldName)
}

func processMsg(md protoreflect.MessageDescriptor, descriptor protoreflect.FieldDescriptor, value protoreflect.Value, resName, parentName protoreflect.Name) bool {
	if md == nil {
		opts := descriptor.Options().(*descriptorpb.FieldOptions)
		dg, ok := proto.GetExtension(opts, policyPB.E_Dgraph).(string)
		if ok && dg != "" {
			fmt.Println("dgraph: ", dg)
			pref.ValueOfString(dg)
			value.List().Set(0, pref.ValueOfString(dg))
		}
		fmt.Printf("resource: %v field: %v value: %v type: %v jsonName: %v parent: %v\n", resName, descriptor.Name(), value, descriptor.Kind().GoString(), descriptor.JSONName(), parentName)
	} else {
		fmt.Printf("resource: %v field: %v type: %v jsonName: %v parent: %v\n", resName, descriptor.Name(), descriptor.Kind().GoString(), descriptor.JSONName(), parentName)
		resName = value.Message().Descriptor().Name()
		rangeOverMsg(value, resName, parentName)
	}
	return true
}

func dynamicMsgFromRequest(apiMsg *apiPB.Request) (*dynamicpb.Message, error) {
	fullName := protoreflect.FullName(apiMsg.GetApiGroup() + "." + apiMsg.GetKind())
	d, err := protoregistry.GlobalFiles.FindDescriptorByName(fullName)
	if err != nil {
		return &dynamicpb.Message{}, err
	}
	msgDesc := d.ParentFile().Messages()
	md := msgDesc.ByName(protoreflect.Name(apiMsg.GetKind()))
	msg, err := dynamicMsg(apiMsg.GetMsg(), md)
	if err != nil {
		return &dynamicpb.Message{}, err
	}
	return msg, nil
}

func dynamicMsg(data []byte, desc protoreflect.MessageDescriptor) (*dynamicpb.Message, error) {
	msg := dynamicpb.NewMessage(desc)
	if err := proto.Unmarshal(data, msg); err != nil {
		return &dynamicpb.Message{}, err
	}
	return msg, nil
}

func (a *api) RegisterType(ctx context.Context, apiMsg *apiPB.Message) (*apiPB.Result, error) {
	fdProto := apiMsg.GetFd()
	fd, err := protodesc.NewFile(fdProto, protoregistry.GlobalFiles)
	if err != nil {
		fmt.Println(err)
	}
	if err := protoregistry.GlobalFiles.RegisterFile(fd); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("registry success")
	}
	return &apiPB.Result{}, nil
}
