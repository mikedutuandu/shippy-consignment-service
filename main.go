package main

import (
	"errors"
	"context"
	"fmt"
	"github.com/micro/go-micro/server"
	pb "github.com/mikedutuandu/shippy-consignment-service/proto/consignment"
	vesselProto "github.com/mikedutuandu/shippy-vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
	"os"
	"github.com/micro/go-micro/metadata"
	userService "github.com/mikedutuandu/shippy-user-service/proto/auth"
)

const (
	defaultHost = "mongodb+srv://admin:AAC0w4Q6jvNv4r8Z@bookingcluster-tgpah.gcp.mongodb.net/test?retryWrites=true&w=majority"
)

var (
	srv micro.Service
)

func main() {
	// Set-up micro instance
	service := micro.NewService(
		micro.Name("shippy.consignment.service"),
		//micro.Version("latest"),
		//micro.WrapHandler(AuthWrapper),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}


	vesselClient := vesselProto.NewVesselService("shippy.vessel.service.cli", service.Client())
	h := &handler{repository, vesselClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(service.Server(), h)

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// AuthWrapper is a high-order function which takes a HandlerFunc
// and returns a function, which takes a context, request and response interface.
// The token is extracted from the context set in our consignment-cli, that
// token is then sent over to the user service to be validated.
// If valid, the call is passed along to the handler. If not,
// an error is returned.
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		// Really shouldn't be using a global here, find a better way
		// of doing this, since you can't pass it into a wrapper.
		authClient := userService.NewAuthService("shippy.auth.client", srv.Client())
		_, err := authClient.ValidateToken(ctx, &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
