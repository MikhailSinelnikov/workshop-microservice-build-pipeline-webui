//go:generate protoc -I echo --go_out=plugins=grpc:echo echo/echo.proto

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	aggregatorpb "github.com/kublr/workshop-microservice-build-pipeline-webui/pkg/aggregator"
)

var (
	port             = flag.Int("port", 8080, "The server port")
	aggregatorAddr   = flag.String("aggregator", "127.0.0.1:11000", "Aggregator address in the format of host:port")
	aggregatorClient aggregatorpb.AggregatorClient
)
var htmlTemplate *template.Template

func main() {
	flag.Parse()

	// prepare template
	var err error
	htmlTemplate, err = template.New("html").Parse(html)
	if err != nil {
		log.Fatalf("fail to parse template: %v", err)
	}

	// establish connection to aggregator
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*aggregatorAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial aggregator: %v", err)
	}

	// create aggregator client stub
	aggregatorClient = aggregatorpb.NewAggregatorClient(conn)

	// register web page handler
	http.HandleFunc("/", handler)

	// start listening
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

const html = `
<html>
	<body>
		Example...
		{{range $i, $r := .Ranges}}
		{{$r.Hot}}<br/>
		{{end}}
	</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	// call aggregator
	clientCtx, clientCancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer clientCancel()

	aggregateResponse, err := aggregatorClient.Aggregate(clientCtx, &aggregatorpb.AggregateRequest{})
	if err != nil {
		aggregateResponse = &aggregatorpb.AggregateResponse{
			Ranges: []*aggregatorpb.ColorRange{},
		}
	}

	// generate page
	htmlTemplate.Execute(w, aggregateResponse)
}
