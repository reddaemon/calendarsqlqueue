package integr

import (
	"fmt"
	"os"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/golang/protobuf/ptypes/timestamp"
	eventpb "github.com/reddaemon/calendarsqlqueue/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var eventServer = os.Getenv("EVENT_SERVICE")
var createID uint32
var createResponse *eventpb.CreateResponse
var updateResponse *eventpb.UpdateResponse
var respErr error

func init() {
	if eventServer == "" {
		eventServer = "localhost:8080"
	}
}

func iCallGrpcEventMethodCreate() error {
	conn, err := grpc.Dial(eventServer, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := eventpb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	inCreate := &eventpb.CreateRequest{
		Event: &eventpb.Event{
			Title:       "Create event test",
			Description: "Create event test",
			Date: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}
	createResponse, respErr = cli.Create(ctx, inCreate)
	return nil
}

func theErrorShouldBeNil() error {
	return respErr
}

func theCreateResponseSuccessShouldBeTrue() error {
	if !createResponse.Success {
		return fmt.Errorf("Create response not success")
	}
	createID = createResponse.Id
	return nil
}

func iCallGrpcEventMethodUpdate() error {
	conn, err := grpc.Dial(eventServer, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}
	defer conn.Close()

	cli := eventpb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	inUpdate := &eventpb.UpdateRequest{
		Id: createID,
		Event: &eventpb.Event{
			Title:       "Update event",
			Description: "Update event",
			Date: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}

	updateResponse, respErr = cli.Update(ctx, inUpdate)

	return nil
}

func theUpdateResponseSuccessShouldBeTrue() error {
	if !updateResponse.Success {
		return fmt.Errorf("Update response not success")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I call grpc event method Create$`, iCallGrpcEventMethodCreate)
	s.Step(`^The error should be nil$`, theErrorShouldBeNil)

	s.Step(`^The create response success should be true$`, theCreateResponseSuccessShouldBeTrue)
	s.Step(`^I call grpc event method Update$`, iCallGrpcEventMethodUpdate)
	s.Step(`^The update response success should be true$`, theUpdateResponseSuccessShouldBeTrue)
}
