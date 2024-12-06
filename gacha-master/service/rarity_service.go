package service

import (
	"context"
	"gacha-master/exception"
	"gacha-master/helper"
	"gacha-master/model/domain"
	"gacha-master/model/web"
	"gacha-master/repository"
	"github.com/go-playground/validator/v10"
)

type RarityService interface {
	Create(ctx context.Context, request *web.RarityCreateRequest) *web.RarityResponse
	FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *web.RarityResponse
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []web.RarityResponse
	Update(ctx context.Context, request *web.RarityUpdateRequest) *web.RarityResponse
	Delete(ctx context.Context, id int, gachaSystemId int)
}

type RarityServiceImpl struct {
	RarityRepository      repository.RarityRepository
	GachaSystemRepository repository.GachaSystemRepository
	Validate              *validator.Validate
}

func NewRarityService(
	rarityRepository repository.RarityRepository,
	gachaSystemRepository repository.GachaSystemRepository,
	validate *validator.Validate,
) RarityService {
	return &RarityServiceImpl{
		RarityRepository:      rarityRepository,
		GachaSystemRepository: gachaSystemRepository,
		Validate:              validate,
	}
}

func (service *RarityServiceImpl) Create(ctx context.Context, request *web.RarityCreateRequest) *web.RarityResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, request.GachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	existingRarity := service.RarityRepository.FindByNameAndGachaSystemId(ctx, request.Name, request.GachaSystemId)
	if existingRarity != nil {
		panic(exception.NewConflictError("Rarity with the same name already exists"))
	}

	rarity := domain.Rarity{
		Name:          request.Name,
		Chance:        request.Chance,
		GachaSystemId: request.GachaSystemId,
	}

	service.RarityRepository.Save(ctx, &rarity)

	return web.ToRarityResponse(&rarity)
}

func (service *RarityServiceImpl) FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *web.RarityResponse {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	rarity := service.RarityRepository.FindByIdAndGachaSystemId(ctx, id, gachaSystemId)
	if rarity == nil {
		panic(exception.NewNotFoundError("Rarity system not found"))
	}

	return web.ToRarityResponse(rarity)
}

func (service *RarityServiceImpl) FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []web.RarityResponse {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	rarities := service.RarityRepository.FindAllByGachaSystemId(ctx, gachaSystemId)
	if rarities == nil {
		return nil
	}

	return web.ToRaritiesResponse(rarities)
}

func (service *RarityServiceImpl) Update(ctx context.Context, request *web.RarityUpdateRequest) *web.RarityResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, request.GachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	rarity := service.RarityRepository.FindByIdAndGachaSystemId(ctx, request.Id, request.GachaSystemId)
	if rarity == nil {
		panic(exception.NewNotFoundError(helper.ErrRarityNotFound))
	}

	request.UpdateRarity(rarity)

	existingRarity := service.RarityRepository.FindByNameAndGachaSystemId(ctx, rarity.Name, rarity.GachaSystemId)
	if existingRarity != nil && existingRarity.Id != rarity.Id {
		panic(exception.NewConflictError("Rarity with the same name already exists"))
	}

	service.RarityRepository.Update(ctx, rarity)

	return web.ToRarityResponse(rarity)
}

func (service *RarityServiceImpl) Delete(ctx context.Context, id int, gachaSystemId int) {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	rarity := service.RarityRepository.FindByIdAndGachaSystemId(ctx, id, gachaSystemId)
	if rarity == nil {
		panic(exception.NewNotFoundError(helper.ErrRarityNotFound))
	}

	service.RarityRepository.Delete(ctx, id, gachaSystemId)
}
