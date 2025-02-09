package main

import (
	"context"
	"log"

	"github.com/gogo/protobuf/proto"
	apiPB "github.com/michaelhenkel/sub/api/proto"
	policyPB "github.com/michaelhenkel/sub/policy/proto"
	serverPB "github.com/michaelhenkel/sub/server/proto"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	client := apiPB.NewApiClient(conn)

	srOne := &serverPB.ServerResourceOne{
		Kind: "ServerResourceOne",
		Name: "name1",
		Spec: &serverPB.ServerResourceOneSpec{
			SrTwo: &serverPB.ServerResourceTwo{
				Kind: "ServerResourceTwo",
				Name: "srName2",
				Spec: &serverPB.ServerResourceTwoSpec{
					SpecName:   "r2specName1",
					SpecValue:  "r2specValue1",
					StartRange: 9,
					EndRange:   15,
					Pol: &policyPB.Pol{
						Poller: "blabla",
					},
				},
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
		Spec:     data,
	}

	_, err = client.Create(context.Background(), apiRequest)
	if err != nil {
		log.Println(err)
	}

}
