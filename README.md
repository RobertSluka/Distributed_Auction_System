# Distributed Auction System

This is a simple distributed auction system implemented in Go. It allows multiple clients to connect to a set of servers to participate in an auction, where they can place bids or request the current highest bid.

## Features
- **Multiple Clients**: Any number of clients can connect to the system and participate.
- **Bid Submission**: Clients can place bids by entering an integer.
- **Query Highest Bid**: Clients can query the current highest bid in real-time.

## How to Use

### 1. Start the Servers
To start the servers, run the following command for each server:

cd server
go run server.go

### 2. Start the Clients
To start a client, use the following command:

  cd client 
  go run client.go client_name
  
Replace client_name with a unique identifier for the client (e.g., Alice, Bob, etc.).

### 3. Client Commands
Once a client is connected, it can:
- **Place a Bid:** Enter an integer in the terminal to place a bid.
- **Request Current Highest Bid:** Type res in the terminal to request the current highest bid.
