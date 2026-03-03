package main

import (
	"context"
	"log"
	"net"

	"crypto-check/pb"

	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"google.golang.org/grpc"
)

// server
type server struct {
	pb.UnimplementedAnalyticsServiceServer
	db *sql.DB
}

func (s *server) GetRSI(ctx context.Context, req *pb.AnalyticRequest) (*pb.AnalyticResponse, error) {

	log.Printf("[gRPC] Received a request for the symbol: %s", req.Symbol)
	// We take the last 14 prices (ORDER BY timestamp DESC will give us the most recent ones on top)
	query := `SELECT price FROM price_history WHERE symbol = ? ORDER BY timestamp DESC LIMIT ?`
	rows, err := s.db.Query(query, req.Symbol, req.Period)
	if err != nil {
		log.Printf("[ERROR] Database query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var prices []float64
	for rows.Next() {
		var p float64
		if err := rows.Scan(&p); err != nil {
			continue
		}
		prices = append(prices, p)
	}

	// If there is little data (for example, it has just been launched), the RSI cannot be calculated
	if len(prices) < 2 {
		return &pb.AnalyticResponse{
			Symbol:   req.Symbol,
			RsiValue: 50.0,
			Status:   "WAITING_FOR_DATA",
		}, nil
	}

	// For the RSI, we need [Old -> New]. Turning the slice over:
	for i, j := 0, len(prices)-1; i < j; i, j = i+1, j-1 {
		prices[i], prices[j] = prices[j], prices[i]
	}

	// Count RSI
	rsi := CalculateRSI(prices)

	status := "NEUTRAL"
	if rsi >= 70 {
		status = "OVERBOUGHT (SELL)"
	} else if rsi <= 30 {
		status = "OVERSOLD (BUY)"
	}

	return &pb.AnalyticResponse{
		Symbol:       req.Symbol,
		CurrentPrice: prices[len(prices)-1], // Last price
		RsiValue:     rsi,
		Status:       status,
	}, nil
}

func main() {
	db, err := sql.Open("sqlite", "./crypto.db")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(s, &server{db: db})

	log.Println("Analytics Service started on port :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
