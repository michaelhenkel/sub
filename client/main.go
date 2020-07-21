package main

import (
	"context"
	"log"

	apiPB "github.com/michaelhenkel/sub/api/proto"
	serverPB "github.com/michaelhenkel/sub/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	client := apiPB.NewApiClient(conn)

	srOne := &serverPB.ServerResourceOne{
		Kind: "ServerResourceOne",
		Name: "name1",
		SrTwo: &serverPB.ServerResourceTwo{
			Kind: "ServerResourceTwo",
			Name: "srName2",
			Spec: &serverPB.ResourceTwoSpec{
				SpecName:  "r2specName1",
				SpecValue: "r2specValue1",
			},
		},
	}
	data, err := proto.Marshal(srOne)
	if err != nil {
		log.Println(err)
	}
	apiRequest := &apiPB.Request{
		Kind:     srOne.GetKind(),
		ApiGroup: "server",
		Msg:      data,
	}

	_, err = client.Create(context.Background(), apiRequest)
	if err != nil {
		log.Println(err)
	}

}
