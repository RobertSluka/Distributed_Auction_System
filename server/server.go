package main

import (
	"Auction/auction"
	proto "Auction/auction"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

var ltime int32 = 0
var timeout int32 = 10

type auctionServer struct {
	proto.UnimplementedAuctionServer
}

func readHighestBid() (int64, string) {
	data, err := os.ReadFile("highest_bid.txt")
	ar := strings.Split(string(data), " ")

	highbid, winner := ar[0], ar[1]

	if err != nil {
		log.Fatal("Something wrong with file")
	}

	highestBid, err := strconv.ParseInt(highbid, 10, 32)
	if err != nil {
		log.Fatal("Something wrong with int")
	}

	return highestBid, winner
}

func (s *auctionServer) Bid(c context.Context, req *auction.BidRequest) (*auction.BidResponse, error) {
	if ltime >= timeout {
		return s.auctionEndedBid()
	}

	amount := req.GetAmount()
	bidClientName := req.GetBidderName()
	highestBid, _ := readHighestBid()
	bidTime := req.GetTime()

	if bidTime > ltime {
		ltime = bidTime
	}

	ltime++

	if amount > highestBid {
		winningBid := strconv.Itoa(int(amount)) + " " + bidClientName
		amount_byte := []byte(winningBid)
		os.WriteFile("highest_bid.txt", amount_byte, 0644)
		highestBid = amount
		log.Printf("New highest bid is %d by %s", amount, bidClientName)
		return &auction.BidResponse{Outcome: auction.BidResponse_SUCCESS, Time: ltime}, nil
	}
	log.Printf("Bid of %d by %s failed", amount, bidClientName)
	return &auction.BidResponse{Outcome: auction.BidResponse_FAIL, Time: ltime}, nil
}

func (s *auctionServer) Result(c context.Context, req *auction.ResultRequest) (*auction.ResultResponse, error) {
	if ltime >= timeout {
		return s.auctionEndedResult()
	}

	highbid, winner := readHighestBid()

	if highbid == 0 {
		outcome := fmt.Sprintf("No bids have been placed")
		return &auction.ResultResponse{Outcome: outcome}, nil

	}

	outcome := fmt.Sprintf("Highest bid: %d by %s", highbid, winner)
	return &auction.ResultResponse{Outcome: outcome}, nil
}

func (s *auctionServer) auctionEndedBid() (*auction.BidResponse, error) {
	return &auction.BidResponse{Outcome: auction.BidResponse_EXCEPTION}, nil
}

func (s *auctionServer) auctionEndedResult() (*auction.ResultResponse, error) {
	highbid, winner := readHighestBid()

	outcome := fmt.Sprintf("Auction ended! Highest bid: %d by %s", highbid, winner)
	return &auction.ResultResponse{Outcome: outcome}, nil
}

func main() {
	reset := []byte("0 None")
	os.WriteFile("highest_bid.txt", reset, 0644)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		lis, err = net.Listen("tcp", ":8081")
		log.Printf("Auction server listening on port 8081")
		if err != nil {
			log.Fatalf("there was an error: %s", err)
		}
	} else {
		log.Printf("Auction server listening on port 8080")
	}

	s := grpc.NewServer()

	proto.RegisterAuctionServer(s, &auctionServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
