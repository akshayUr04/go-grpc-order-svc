package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akshayUr04/go-grpc-order-svc/pkg/client"
	"github.com/akshayUr04/go-grpc-order-svc/pkg/db"
	"github.com/akshayUr04/go-grpc-order-svc/pkg/models"
	"github.com/akshayUr04/go-grpc-order-svc/pkg/pb"
)

type Server struct {
	H          db.Handler
	ProductSvc client.ProductServiceClient
	CartSvc    client.CartServiceClient

	pb.UnimplementedOrderServiceServer
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	fmt.Println("Order Service :  CreateOrder")
	product, err := s.ProductSvc.FindOne(req.ProductId)

	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest, Error: err.Error(),
		}, nil
	} else if product.Status >= http.StatusNotFound {
		return &pb.CreateOrderResponse{
			Status: product.Status,
			Error:  product.Error,
		}, nil
	} else if product.Data.Stock < req.Quantity {
		return &pb.CreateOrderResponse{
			Status: http.StatusConflict,
			Error:  "Stock too less",
		}, nil
	}

	order := models.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		Quantity:  req.Quantity,
		UserId:    req.UserId,
	}

	s.H.DB.Create(&order)

	res, err := s.ProductSvc.DecreaseStock(req.ProductId, order.Id, req.Quantity)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if res.Status == http.StatusConflict {
		s.H.DB.Delete(&models.Order{}, order.Id)

		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

	return &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}

func (s *Server) OrderFromCart(ctx context.Context, req *pb.OrderFromCartRequest) (*pb.CreateOrderResponse, error) {
	res, err := s.CartSvc.FindCart(req.UserId)
	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, err
	}
	for _, products := range res.Data {
		product, err := s.ProductSvc.FindOne(products.ProductId)
		if err != nil {
			return &pb.CreateOrderResponse{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			}, err
		}
		if product.Data.Stock < products.Qty {
			return &pb.CreateOrderResponse{
				Status: http.StatusBadRequest,
				Error:  "no stock",
			}, fmt.Errorf("no stock")
		}
		order := models.Order{
			Price:     products.Total,
			ProductId: products.ProductId,
			Quantity:  products.Qty,
			UserId:    req.UserId,
		}

		s.H.DB.Create(&order)

		_, err = s.ProductSvc.DecreaseStock(products.ProductId, order.Id, products.Qty)
		if err != nil {
			return &pb.CreateOrderResponse{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			}, err
		}
	}
	_, err = s.CartSvc.DeletCart(req.UserId)
	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, err
	}
	return &pb.CreateOrderResponse{
		Status: http.StatusOK,
	}, nil
}
