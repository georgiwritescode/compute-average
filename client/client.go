package main

import (
	computeAverage "compute-average/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	log.Println("[INFO] Client has started ...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[ERROR] Failed to Dial: %v", err)
	}
	defer cc.Close()

	c := computeAverage.NewComputeServiceClient(cc)

	doClientStreamingComputeAverage(c)
}

func doClientStreamingComputeAverage(c computeAverage.ComputeServiceClient) {
	log.Println("[INFO] doClientStreamingComputeAverage invoked ...")

	requests := []*computeAverage.ComputeAverageRequest{
		{Number: 5},
		{Number: 3},
		{Number: 12},
		{Number: 97},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("[ERROR] Failed call Compute Average: %v", err)
	}

	for _, req := range requests {
		log.Printf("[INFO] Sending request: %v", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[ERROR] Could not receive response from Compute Average: %v", err)
	}
	log.Printf("[INFO] Response: %v", res)
}
