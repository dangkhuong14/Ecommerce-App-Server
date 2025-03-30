package grpcservice

import (
	"context"
	"ecommerce/common"
	"ecommerce/gen/category"
	"ecommerce/module/category/query"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
)

type gRPCCategoryServer struct {
	port int
	sctx sctx.ServiceContext
}

func NewGRPCCategoryServer(sctx sctx.ServiceContext, port int) *gRPCCategoryServer {
	return &gRPCCategoryServer{sctx: sctx, port: port}
}

type categoryServiceServer struct {
	sctx sctx.ServiceContext
	category.UnimplementedCategoryServiceServer
}

func NewCategoryServiceServer(sctx sctx.ServiceContext) *categoryServiceServer {
	return &categoryServiceServer{sctx: sctx}
}

func (s *gRPCCategoryServer) Start() error {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	// Create a gRPC server object
	server := grpc.NewServer()
	// Attach the Greeter service to the server
	category.RegisterCategoryServiceServer(server, NewCategoryServiceServer(s.sctx))
	// Serve gRPC Server
	log.Printf("Serving gRPC on 0.0.0.0:%d", s.port)

	// Run server and return error if exists
	if err := server.Serve(lis); err != nil {
		log.Printf("gRPC server error: %v", err)
		return err
	}

	return nil
}

func (s *categoryServiceServer) FindCategoriesByIDs(ctx context.Context, req *category.FindCategoriesReq) (*category.FindCategoriesRes, error) {

	// Convert request's IDs from string to UUID
	uuidIDs := make([]common.UUID, len(req.Ids))

	for i := range uuidIDs {
		uuidIDs[i] = common.UUID(uuid.MustParse(req.Ids[i]))
	}

	// Call use case to get categories by IDs
	findCtgsByIDsQuery := query.NewFindCategoriesByIDsQuery(s.sctx)
	categories, err := findCtgsByIDsQuery.Execute(ctx, uuidIDs)
	if err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list categories").WithDebug(err.Error())
	}

	// Convert CategoryDTO to gRPC generated type
	results := make([]*category.CategoryDTO, len(categories))
	for i := range categories {
		results[i] = &category.CategoryDTO{
			Id:    (categories[i].Id).String(),
			Title: categories[i].Title,
		}
	}
	return &category.FindCategoriesRes{Data: results}, nil
}
