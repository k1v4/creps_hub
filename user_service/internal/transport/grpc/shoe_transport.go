package grpc

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	userv1 "github.com/k1v4/protos/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/internal/models"
)

type IShoeService interface {
	AddShoe(ctx context.Context, userID int64, name, imageUrl string) (int64, error)
	GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error)
	DeleteShoe(ctx context.Context, shoeID int64) (bool, error)
	UpdateShoe(ctx context.Context, shoeId, userId int64, name, imageUrl string) (*models.Shoe, error)
	GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error)
}

type ShoeService struct {
	userv1.UnimplementedShoeServiceServer
	service IShoeService
}

func NewShoeService(service IShoeService) *ShoeService {
	return &ShoeService{service: service}
}

func (s *ShoeService) AddShoe(ctx context.Context, req *userv1.AddShoeRequest) (*userv1.AddShoeResponse, error) {
	const op = "ShoeTransport.AddShoe"

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	imageUrl := req.GetImageUrl()
	if imageUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	userId := req.GetUserId()
	if userId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong user id")
	}

	shoeId, err := s.service.AddShoe(ctx, userId, name, imageUrl)
	if err != nil {
		//TODO добавить доп проверки
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userv1.AddShoeResponse{
		ShoeId: shoeId,
	}, nil
}

func (s *ShoeService) GetShoe(ctx context.Context, req *userv1.GetShoeRequest) (*userv1.GetShoeResponse, error) {
	const op = "ShoeTransport.GetShoe"

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	shoe, err := s.service.GetShoe(ctx, shoeId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	r := pointer.Get(shoe)

	return &userv1.GetShoeResponse{
		Shoe: &userv1.Shoe{
			ShoeId:   shoeId,
			Name:     r.Name,
			ImageUrl: r.ImageUrl,
			UserId:   r.UserId,
		},
	}, nil
}

func (s *ShoeService) DeleteShoe(ctx context.Context, req *userv1.DeleteShoeRequest) (*userv1.DeleteShoeResponse, error) {
	const op = "ShoeTransport.DeleteShoe"

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	result, err := s.service.DeleteShoe(ctx, shoeId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userv1.DeleteShoeResponse{IsSuccessfully: result}, nil
}

func (s *ShoeService) UpdateShoe(ctx context.Context, req *userv1.UpdateShoeRequest) (*userv1.UpdateShoeResponse, error) {
	const op = "ShoeTransport.UpdateShoe"

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	imageUrl := req.GetImageUrl()
	if imageUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	userId := req.GetUserId()
	if userId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong user id")
	}

	shoe, err := s.service.UpdateShoe(ctx, shoeId, userId, name, imageUrl)
	if err != nil {
		return nil, err
	}

	r := pointer.Get(shoe)

	return &userv1.UpdateShoeResponse{
		Shoe: &userv1.Shoe{
			ShoeId:   r.Id,
			Name:     r.Name,
			ImageUrl: r.ImageUrl,
			UserId:   r.UserId,
		},
	}, nil
}

func (s *ShoeService) GetShoes(ctx context.Context, req *userv1.GetAllShoesRequest) (*userv1.GetAllShoesResponse, error) {
	const op = "ShoeTransport.GetAllShoes"

	userId := req.GetUserId()
	if userId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong user id")
	}

	shoes, err := s.service.GetShoes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	r := pointer.Get(shoes)
	var res []*userv1.Shoe

	for _, o := range r {
		res = append(res, &userv1.Shoe{
			ShoeId:   o.Id,
			Name:     o.Name,
			ImageUrl: o.ImageUrl,
			UserId:   o.UserId,
		})
	}

	return &userv1.GetAllShoesResponse{Shoes: res}, nil
}
