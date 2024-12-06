package service

import (
	"context"
	"gacha-auth/exception"
	"gacha-auth/helper"
	"gacha-auth/model/domain"
	"gacha-auth/model/web"
	"gacha-auth/repository"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"time"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) web.UserResponse
	Login(ctx context.Context, request web.UserLoginRequest) web.UserResponse
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
	TokenAuth      *jwtauth.JWTAuth
}

func NewUserService(userRepository repository.UserRepository, validate *validator.Validate, tokenAuth *jwtauth.JWTAuth) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
		TokenAuth:      tokenAuth,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err, "Request validation failed")

	existingUser := service.UserRepository.FindByName(ctx, request.Name)
	if existingUser != nil {
		panic(exception.NewConflictError("User with the same name already exists"))
	}

	existingUser = service.UserRepository.FindByUsername(ctx, request.Username)
	if existingUser != nil {
		panic(exception.NewConflictError("User with the same username already exists"))
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	helper.PanicIfError(err, "Failed to hash password")

	user := domain.User{
		Name:     request.Name,
		Username: request.Username,
		Password: hashedPassword,
	}

	service.UserRepository.Save(ctx, &user)
	helper.PanicIfError(err, "Error save user")

	claims := jwt.MapClaims{
		"userId": user.Id,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*3))
	_, tokenString, err := service.TokenAuth.Encode(claims)

	return web.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		UserToken: tokenString,
	}
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) web.UserResponse {
	user := service.UserRepository.FindByUsername(ctx, request.Username)
	if user == nil {
		panic(exception.NewUserError("Invalid username or password"))
	}

	err := helper.VerifyPassword(user.Password, request.Password)
	if err != nil {
		panic(exception.NewUserError("Invalid username or password"))
	}

	claims := jwt.MapClaims{
		"userId": user.Id,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*3))
	_, tokenString, err := service.TokenAuth.Encode(claims)

	return web.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		UserToken: tokenString,
	}
}
