package main

import (
	"Auction/auction"
	proto "Auction/auction"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type auctionServer struct {
	proto.UnimplementedAuctionServer
	highestBid int32
	winner     string
}

func (s *auctionServer) Bid(c context.Context, req *auction.BidRequest) (*auction.BidResponse, error) {
	amount := req.GetAmount()

	if amount > s.highestBid {
		s.highestBid = amount
		s.winner = "current_client" // needs to be modified with an actual identifier

		return &auction.BidResponse{Outcome: auction.BidResponse_SUCCESS}, nil
	}
	return &auction.BidResponse{Outcome: auction.BidResponse_FAIL}, nil
}

func (s *auctionServer) Result(c context.Context, req *auction.ResultRequest) (*auction.ResultResponse, error) {
	if s.highestBid == 0 {
		return nil, status.Error(codes.NotFound, "No bids have been placed.")
	}

	outcome := fmt.Sprintf("Highest bid: %d by %s", s.highestBid, s.winner)
	return &auction.ResultResponse{Outcome: outcome}, nil
}

// The main function starts the server and listens on port 8080.
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Auction server listening on port 8080")

	s := grpc.NewServer()

	proto.RegisterAuctionServer(s, &auctionServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
