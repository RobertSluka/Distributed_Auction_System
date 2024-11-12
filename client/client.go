package main

import (
	"Auction/auction"
	proto "Auction/auction"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuctionClient(conn)

	testBid(client, 100)
	testBid(client, 200)
	testBid(client, 150)

	getAuctionResult(client)

}

func testBid(client auction.AuctionClient, amount int32) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Bid(c, &auction.BidRequest{Amount: amount})
	if err != nil {
		log.Fatalf("Error when bidding: %v", err)
	}

	outcome := "SUCCESS"
	if res.GetOutcome() == auction.BidResponse_FAIL {
		outcome = "FAIL"
	}
	fmt.Printf("Bid of %d: %s\n", amount, outcome)
}

func getAuctionResult(client auction.AuctionClient) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Result(c, &auction.ResultRequest{})
	if err != nil {
		log.Fatalf("Error when fetching result: %v", err)
	}

	fmt.Printf("Auction result: %s\n", res.GetOutcome())
}
