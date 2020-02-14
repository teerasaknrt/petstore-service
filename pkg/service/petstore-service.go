package v1

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	v1 "petstore-service/pkg/api/v1"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// petstoreServiceServer is implementation of v1.ToDoServiceServer proto interface
type petStoreServiceServer struct {
	db *firestore.Client
}

// NewPatstoreServiceServer creates PatStore service
func NewPetStoreServiceServer(db *firestore.Client) v1.PetStoreServiceServer {
	return &petStoreServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *petStoreServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Add a new pet to the store
func (s *petStoreServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	s.checkAPI(req.Api)
	id := randStringRunes(20)
	//add data to db
	states := s.db.Collection("petstore")
	ny := states.Doc(id)
	_, err := ny.Create(ctx, v1.Data{
		Id: id,
		Category: &v1.Category{
			Id:   req.Data.Category.Id,
			Name: req.Data.Name,
		},
		Name:      req.Data.Name,
		PhotoUrls: req.Data.PhotoUrls,
		Tags:      req.Data.Tags,
		Status:    req.Data.Status,
		IsDelete:  false,
	})
	if err != nil {
		log.Println("connot create a new petstore")
	}
	res := &v1.CreateResponse{
		Api:     apiVersion,
		Message: "success",
	}
	log.Println("success : create a new petstore")
	return res, nil
}

// Update an pet
func (s *petStoreServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	s.checkAPI(req.Api)
	states := s.db.Collection("petstore")
	ny := states.Doc(req.Data.Id)
	_, err := ny.Update(ctx, []firestore.Update{{Path: "Status", Value: req.Data.Status}})
	if err != nil {
		log.Println("cannot update petsore")
		return nil, err
	}
	res := &v1.UpdateResponse{
		Api:     "v1",
		Message: "success",
	}
	return res, nil
}

// Find pet by ID
func (s *petStoreServiceServer) Find(ctx context.Context, req *v1.FindRequest) (*v1.FindResponse, error) {
	data := &v1.Data{}
	states := s.db.Collection("petstore")
	ny := states.Doc(req.Id)
	docsnap, err := ny.Get(ctx)
	if err != nil {
		fmt.Println("cannot get petstore")
		return nil, err

	}
	mapstructure.Decode(docsnap.Data(), &data)

	res := &v1.FindResponse{
		Api:  apiVersion,
		Data: data,
	}
	log.Println("success : get a new petstore")

	return res, nil
}

// Deletes a pet
func (s *petStoreServiceServer) Deletes(ctx context.Context, req *v1.DeletesRequest) (*v1.DeletesResponse, error) {
	states := s.db.Collection("petstore")
	ny := states.Doc(req.Id)
	_, err := ny.Update(ctx, []firestore.Update{{Path: "IsDelete", Value: true}})
	if err != nil {
		log.Println("cannot deletes petsore")
		return nil, err
	}
	res := &v1.DeletesResponse{
		Api:     "v1",
		Message: "success",
	}
	log.Println("success : delete a petstore")
	return res, nil
}

var letterRunes = []rune("_-0987654321abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
