package main

import (
	"context"
	"log"
	pb "github.com/mikedutuandu/shippy-consignment-service/proto/consignment"
	vesselProto "github.com/mikedutuandu/shippy-vessel-service/proto/vessel"
)

type handler struct {
	repository
	vesselClient vesselProto.VesselService
}


// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	println("go to create")
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	if err = s.repository.Create(req); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments -
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.repository.GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
