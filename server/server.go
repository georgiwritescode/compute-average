package main

import (
	computeAverage "compute-average/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) ComputeAverage(stream computeAverage.ComputeService_ComputeAverageServer) error {
	log.Println("[INFO] ComputeAverage function invoked")

	sum := 0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avg := float32(sum) / float32(count)
			return stream.SendAndClose(&computeAverage.ComputeAverageResponse{Average: avg})
		}
		if err != nil {
			log.Fatalf("[ERROR] Cannot read client stream: %v", err)
		}

		sum += int(req.GetNumber())
		count++
	}
}

func main() {
	log.Println("[INFO] Server has started ...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	computeAverage.RegisterComputeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
