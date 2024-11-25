package main

import (
	"Auction/auction"
	proto "Auction/auction"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ltime int32 = 0

func main() {

	clientName := os.Args[1]
	log.Printf("%s joined the auction!", clientName)

	conn1, err := grpc.NewClient("0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn1.Close()

	replica := proto.NewAuctionClient(conn1)

	conn2, err := grpc.NewClient("0.0.0.0:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn2.Close()

	replica2 := proto.NewAuctionClient(conn2)

	for {
		var input string
		fmt.Scan(&input)

		if strings.Contains(input, "res") {
			getAuctionResult(replica, replica2)
		} else {
			amount, err := strconv.Atoi(input)
			testBid(clientName, replica, replica2, int64(amount))
			if err != nil {
				log.Panic("oh no")
			}

		}

	}

}

func testBid(clientName string, client1 auction.AuctionClient, client2 auction.AuctionClient, amount int64) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ltime++
	res, err := client1.Bid(c, &auction.BidRequest{Amount: amount, BidderName: clientName, Time: ltime})
	if err != nil {
		log.Printf("Failed to bid, switching to replica")
		res, err = client2.Bid(c, &auction.BidRequest{Amount: amount, BidderName: clientName, Time: ltime})
		if err != nil {
			log.Fatalf("Error when bidding: %v", err)
		}
	}

	outcome := "SUCCESS"
	if res.GetOutcome() == auction.BidResponse_FAIL {
		outcome = "FAIL"
	}

	if res.GetTime() > ltime {
		ltime = res.GetTime()
	}

	fmt.Printf("Bid of %d: %s\n", amount, outcome)
}

func getAuctionResult(client1 auction.AuctionClient, client2 auction.AuctionClient) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client1.Result(c, &auction.ResultRequest{})
	if err != nil {
		log.Printf("Error, connecting to replica")
		res, err = client2.Result(c, &auction.ResultRequest{})
		if err != nil {
			log.Fatalf("Error when fetching result: %v", err)
		}
	}

	fmt.Printf("Auction result: %s\n", res.GetOutcome())
}
