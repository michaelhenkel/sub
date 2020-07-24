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

//go:generate protoc -I../policy/proto --proto_path=proto --gogo_out=plugins=grpc:proto --gogo_opt=paths=source_relative proto/api.proto

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
	//blabla
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
		ok := processMsg(descriptor.Message(), descriptor, value, resName, parentName)
		if !ok {
			fmt.Println("violation")
		} else {
			fmt.Println("no violation")
		}
		return ok
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
		minRange, ok := proto.GetExtension(opts, policyPB.E_RangeMin).(int32)
		if ok && minRange > 0 {
			actMinRange := value.Int()
			if actMinRange < int64(minRange) {
				fmt.Println("minRange too small")
				return false
			} else {
				fmt.Println("minRange ok")
			}
		}
		maxRange, ok := proto.GetExtension(opts, policyPB.E_RangeMax).(int32)
		if ok && maxRange > 0 {
			actMaxRange := value.Int()
			if actMaxRange < int64(maxRange) {
				fmt.Println("maxRange too small")
				return false
			} else {
				fmt.Println("maxRange ok")
			}
		}
		//fmt.Printf("resource: %v field: %v value: %v type: %v jsonName: %v parent: %v\n", resName, descriptor.Name(), value, descriptor.Kind().GoString(), descriptor.JSONName(), parentName)
	} else {
		//fmt.Printf("resource: %v field: %v type: %v jsonName: %v parent: %v\n", resName, descriptor.Name(), descriptor.Kind().GoString(), descriptor.JSONName(), parentName)
		resName = value.Message().Descriptor().Name()
		rangeOverMsg(value, resName, parentName)
	}
	return true
}

func dynamicMsgFromRequest(apiMsg *apiPB.Request) (*dynamicpb.Message, error) {
	fmt.Println(apiMsg.ProtoReflect().Descriptor().Parent().FullName())
	fullName := protoreflect.FullName(apiMsg.GetApiGroup() + "." + apiMsg.GetKind())
	d, err := protoregistry.GlobalFiles.FindDescriptorByName(fullName)
	if err != nil {
		return &dynamicpb.Message{}, err
	}
	msgDesc := d.ParentFile().Messages()
	md := msgDesc.ByName(protoreflect.Name(apiMsg.GetKind()))
	msg, err := dynamicMsg(apiMsg.GetSpec(), md)
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
