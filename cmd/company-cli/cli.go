package main

import (
    "io/ioutil"
    "encoding/json"
    "github.com/Piszmog/microservice-example/cmd/company-service/pb"
    "google.golang.org/grpc"
    "log"
    "os"
    "golang.org/x/net/context"
)

const (
    ADDRESS         = "localhost:8080"
    DEFAULTFILENAME = "C:/Users/rande/go/src/github.com/Piszmog/microservice-example/cmd/company-cli/company.json"
)

func parseFile(file string) (*pb.Company, error) {
    var company *pb.Company
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    json.Unmarshal(data, &company)
    return company, err
}

func main() {
    // Set up a connection to the server.
    conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Did not connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewCompanyServiceClient(conn)

    // Contact the server and print out its response.
    file := DEFAULTFILENAME
    if len(os.Args) > 1 {
        file = os.Args[1]
    }

    consignment, err := parseFile(file)

    if err != nil {
        log.Fatalf("Could not parse file: %v", err)
    }

    r, err := client.CreateCompany(context.Background(), consignment)
    if err != nil {
        log.Fatalf("Could not greet: %v", err)
    }
    log.Printf("Created: %t", r.Created)

    getAll, err := client.GetCompanies(context.Background(), &pb.GetRequest{})
    if err != nil {
        log.Fatalf("Could not list consignments: %v", err)
    }
    for _, v := range getAll.Companies {
        log.Println(v)
    }
}
