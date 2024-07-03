package main

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"context"
	"github.com/labstack/gommon/log"
	api "github.com/reddaemon/calendarsqlqueue/protofiles/protofiles/api"
)

var createResponse *api.CreateResponse
var updateResponse *api.UpdateResponse

//var readResponse api.UpdateResponse
//var deleteResponse api.UpdateResponse

func main() {
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)
	cc, err := grpc.NewClient("127.0.0.1:8888", grpc.WithTransportCredentials(nil))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err = cc.Close()
		if err != nil {
			log.Fatalf("unable to close connection: %v", err)
		}
	}(cc)
	c := api.NewEventServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	inCreate := &api.CreateRequest{
		Event: &api.Event{
			Title:       "Test",
			Description: "testing",
			Date:        timestamppb.New(time.Now()),
		},
	}

	inUpdate := &api.UpdateRequest{
		Id: 21,
		Event: &api.Event{
			Title:       "Test1",
			Description: "Behavior testing1",
			Date:        timestamppb.New(time.Now()),
		},
	}
	inRead := api.ReadRequest{
		Id: 21,
	}
	inDelete := api.DeleteRequest{
		Id: 20,
	}
	createResponse, err = c.Create(ctx, inCreate)
	if err != nil {
		fmt.Println("CREATE:", createResponse, err)
	}

	fmt.Printf("%v\n", createResponse)
	readResponse, err := c.Read(ctx, &inRead)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", readResponse)
	updateResponse, err = c.Update(ctx, inUpdate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", updateResponse)
	readResponse, err = c.Read(ctx, &inRead)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", readResponse)
	deleteResponse, err := c.Delete(ctx, &inDelete)
	if err != nil {
		fmt.Println("DELETEERROR: ", err)
	}
	fmt.Printf("%v\n", deleteResponse)
}
