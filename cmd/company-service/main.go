package main

import (
    "github.com/Piszmog/microservice-example/cmd/company-service/pb"
    "golang.org/x/net/context"
    "net"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

const (
    PORT = ":8080"
)

type IRepository interface {
    Create(company *pb.Company) (*pb.Company, error)
    GetAll() []*pb.Company
}

// The Repository -- dummy for now
type Repository struct {
    companies []*pb.Company
}

func (repo *Repository) Create(company *pb.Company) (*pb.Company, error) {
    updated := append(repo.companies, company)
    repo.companies = updated
    return company, nil
}

func (repo *Repository) GetAll() []*pb.Company {
    return repo.companies
}

type service struct {
    repo IRepository
}

func (s *service) CreateCompany(ctx context.Context, request *pb.Company) (*pb.Response, error) {
    company, err := s.repo.Create(request)
    if err != nil {
        return nil, err
    }
    return &pb.Response{Created: true, Company: company}, nil
}

func (s *service) GetCompanies(ctx context.Context, request *pb.GetRequest) (*pb.Response, error) {
    companies := s.repo.GetAll()
    return &pb.Response{Companies: companies}, nil
}

func main() {
    repo := &Repository{}

    // set up server
    lis, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    // Register with gRPC server -- uses autogen code
    pb.RegisterCompanyServiceServer(s, &service{repo})
    // register reflection service ongRPC server
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to server %v", err)
    }
}
