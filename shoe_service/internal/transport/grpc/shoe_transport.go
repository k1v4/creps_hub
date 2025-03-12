package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	shoev1 "github.com/k1v4/protos/gen/shoe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/internal/models"
	"user_service/internal/service"
	"user_service/pkg/jwtpkg"
)

type IShoeService interface {
	AddShoe(ctx context.Context, userID int64, name string, imageData []byte) (int64, error)
	GetShoe(ctx context.Context, shoeId int64) (*models.Shoe, error)
	DeleteShoe(ctx context.Context, shoeID int64) (bool, error)
	UpdateShoe(ctx context.Context, shoeId, userId int64, name string, imageUrl []byte) (*models.Shoe, error)
	GetShoes(ctx context.Context, userId int64) (*[]models.Shoe, error)
}

type ShoeService struct {
	shoev1.UnimplementedShoeServiceServer
	service IShoeService
}

func NewShoeService(service IShoeService) *ShoeService {
	return &ShoeService{service: service}
}

func (s *ShoeService) AddShoe(ctx context.Context, req *shoev1.AddShoeRequest) (*shoev1.AddShoeResponse, error) {
	const op = "ShoeTransport.AddShoe"

	token := jwtpkg.ExtractToken(ctx)
	if token == "" {
		return nil, status.Error(codes.PermissionDenied, "token is empty")
	}

	userId, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token is invalid")
	}

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	imageName := req.GetImageName()
	if len([]rune(imageName)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image name is empty")
	}

	imageData := req.GetImageData()
	if len(imageData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	shoeId, err := s.service.AddShoe(ctx, userId, imageName, imageData)
	if err != nil {
		//TODO добавить доп проверки
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &shoev1.AddShoeResponse{
		ShoeId: shoeId,
	}, nil
}

func (s *ShoeService) GetShoe(ctx context.Context, req *shoev1.GetShoeRequest) (*shoev1.GetShoeResponse, error) {
	const op = "ShoeTransport.GetShoe"

	token := jwtpkg.ExtractToken(ctx)
	if token == "" {
		return nil, status.Error(codes.PermissionDenied, "token is empty")
	}

	_, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token is invalid")
	}

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	shoe, err := s.service.GetShoe(ctx, shoeId)
	if err != nil {
		if errors.Is(err, service.ErrShoeNotFound) {
			return nil, status.Error(codes.NotFound, "shoe not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	r := pointer.Get(shoe)

	return &shoev1.GetShoeResponse{
		Shoe: &shoev1.Shoe{
			ShoeId:   shoeId,
			Name:     r.Name,
			ImageUrl: r.ImageUrl,
			UserId:   r.UserId,
		},
	}, nil
}

func (s *ShoeService) DeleteShoe(ctx context.Context, req *shoev1.DeleteShoeRequest) (*shoev1.DeleteShoeResponse, error) {
	const op = "ShoeTransport.DeleteShoe"

	token := jwtpkg.ExtractToken(ctx)
	if token == "" {
		return nil, status.Error(codes.PermissionDenied, "token is empty")
	}

	_, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token is invalid")
	}

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	result, err := s.service.DeleteShoe(ctx, shoeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &shoev1.DeleteShoeResponse{IsSuccessfully: result}, nil
}

func (s *ShoeService) UpdateShoe(ctx context.Context, req *shoev1.UpdateShoeRequest) (*shoev1.UpdateShoeResponse, error) {
	const op = "ShoeTransport.UpdateShoe"

	token := jwtpkg.ExtractToken(ctx)
	if token == "" {
		return nil, status.Error(codes.PermissionDenied, "token is empty")
	}

	userId, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token is invalid")
	}

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	imageData := req.GetImageData()
	if len(imageData) == 9 {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	newFileName := req.GetNewFileName()
	if len([]rune(newFileName)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "new file name is empty")
	}

	shoeId := req.GetShoeId()
	if shoeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong shoe id")
	}

	shoe, err := s.service.UpdateShoe(ctx, shoeId, userId, newFileName, imageData)
	if err != nil {
		if errors.Is(err, service.ErrShoeNotFound) {
			return nil, status.Error(codes.NotFound, "shoe not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	r := pointer.Get(shoe)

	return &shoev1.UpdateShoeResponse{
		Shoe: &shoev1.Shoe{
			ShoeId:   r.Id,
			Name:     r.Name,
			ImageUrl: r.ImageUrl,
			UserId:   r.UserId,
		},
	}, nil
}

func (s *ShoeService) GetShoes(ctx context.Context, req *shoev1.GetAllShoesRequest) (*shoev1.GetAllShoesResponse, error) {
	const op = "ShoeTransport.GetAllShoes"

	token := jwtpkg.ExtractToken(ctx)
	if token == "" {
		return nil, status.Error(codes.PermissionDenied, "token is empty")
	}

	userId, err := jwtpkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token is invalid")
	}

	shoes, err := s.service.GetShoes(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	r := pointer.Get(shoes)
	var res []*shoev1.Shoe

	for _, o := range r {
		res = append(res, &shoev1.Shoe{
			ShoeId:   o.Id,
			Name:     o.Name,
			ImageUrl: o.ImageUrl,
			UserId:   o.UserId,
		})
	}

	return &shoev1.GetAllShoesResponse{Shoes: res}, nil
}
