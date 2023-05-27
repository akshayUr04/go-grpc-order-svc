package client

import (
	"context"
	"fmt"

	"github.com/akshayUr04/go-grpc-order-svc/pkg/pb"
	"google.golang.org/grpc"
)

type CartServiceClient struct {
	Client pb.CartServiceClient
}

func InitCartServiceClient(url string) CartServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	c := CartServiceClient{
		Client: pb.NewCartServiceClient(cc),
	}
	return c
}

func (c *CartServiceClient) FindCart(userId int64) (*pb.FindCartResponse, error) {
	req := &pb.FindCartRequest{
		UserId: userId,
	}
	return c.Client.FindCart(context.Background(), req)
}

func (c *CartServiceClient) DeletCart(userId int64) (*pb.DeleteCartResponse, error) {
	req := &pb.DeleteCartRequest{
		UserId: userId,
	}
	return c.Client.DeleteCart(context.Background(), req)
}
