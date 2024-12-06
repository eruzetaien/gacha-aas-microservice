package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"gacha-master/exception"
	"gacha-master/helper"
	"gacha-master/model/domain"
	"gacha-master/model/web"
	"gacha-master/repository"
	"github.com/go-playground/validator/v10"
	"os"
)

type GachaSystemService interface {
	Create(ctx context.Context, request *web.GachaSystemCreateRequest) *web.GachaSystemDetailResponse
	FindById(ctx context.Context, id int) *web.GachaSystemDetailResponse
	Delete(ctx context.Context, id int)
	FindAllByUserId(ctx context.Context) []web.GachaSystemResponse
	FindByNameAndUserId(ctx context.Context, name string, userId int) *web.GachaSystemDetailResponse
}

type GachaSystemServiceImpl struct {
	GachaSystemRepository repository.GachaSystemRepository
	RarityRepository      repository.RarityRepository
	CharacterRepository   repository.CharacterRepository
	Validate              *validator.Validate
}

func NewGachaSystemService(
	gachaSystemRepository repository.GachaSystemRepository,
	rarityRepository repository.RarityRepository,
	characterRepository repository.CharacterRepository,
	validate *validator.Validate,
) GachaSystemService {
	return &GachaSystemServiceImpl{
		GachaSystemRepository: gachaSystemRepository,
		RarityRepository:      rarityRepository,
		CharacterRepository:   characterRepository,
		Validate:              validate,
	}
}

func (service *GachaSystemServiceImpl) Create(ctx context.Context, request *web.GachaSystemCreateRequest) *web.GachaSystemDetailResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	existingGachaSystem := service.GachaSystemRepository.FindByNameAndUserId(ctx, request.Name, userId)
	if existingGachaSystem != nil {
		panic(exception.NewConflictError("Gacha system with the same name already exists"))
	}

	gachaSystem := domain.GachaSystem{
		Name:     request.Name,
		UserId:   userId,
		Endpoint: createEndpoint(request.Name, userId),
	}

	service.GachaSystemRepository.Save(ctx, &gachaSystem)

	return &web.GachaSystemDetailResponse{
		Id:       gachaSystem.Id,
		Name:     gachaSystem.Name,
		Endpoint: gachaSystem.Endpoint,
	}
}

func (service *GachaSystemServiceImpl) FindById(ctx context.Context, id int) *web.GachaSystemDetailResponse {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, id, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	return &web.GachaSystemDetailResponse{
		Id:       gachaSystem.Id,
		Name:     gachaSystem.Name,
		Endpoint: gachaSystem.Endpoint,
	}
}

func (service *GachaSystemServiceImpl) Delete(ctx context.Context, id int) {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, id, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	service.GachaSystemRepository.Delete(ctx, id)

}

func (service *GachaSystemServiceImpl) FindAllByUserId(ctx context.Context) []web.GachaSystemResponse {
	userId := helper.ExtractUserID(ctx)

	gachaSystems := service.GachaSystemRepository.FindAllByUserId(ctx, userId)
	if gachaSystems == nil {
		return nil
	}

	return web.ToGachaSystemsResponse(gachaSystems)
}

func (service *GachaSystemServiceImpl) FindByNameAndUserId(ctx context.Context, name string, userId int) *web.GachaSystemDetailResponse {
	gachaSystem := service.GachaSystemRepository.FindByNameAndUserId(ctx, name, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	return &web.GachaSystemDetailResponse{
		Id:       gachaSystem.Id,
		Name:     gachaSystem.Name,
		Endpoint: gachaSystem.Endpoint,
	}
}

func shiftString(input string, shiftCount int) string {
	n := len(input)
	shiftCount = shiftCount % n
	return input[shiftCount:] + input[:shiftCount]
}

func createEndpoint(gachaName string, userId int) string {
	shiftCount := userId % 6
	shiftedGachaName := shiftString(gachaName, shiftCount)
	combined := fmt.Sprintf("%s%d", shiftedGachaName, userId)

	var sha = sha1.New()
	sha.Write([]byte(combined))
	var encrypted = sha.Sum(nil)
	gachaEndpoint := fmt.Sprintf("%s/%x", os.Getenv("GACHA_PULL_URL"), encrypted)

	return gachaEndpoint
}
