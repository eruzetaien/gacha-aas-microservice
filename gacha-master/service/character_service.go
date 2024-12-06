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

type CharacterService interface {
	Create(ctx context.Context, request *web.CharacterCreateRequest) *web.CharacterResponse
	Update(ctx context.Context, request *web.CharacterUpdateRequest) *web.CharacterResponse
	InsertImageUrl(ctx context.Context, request *web.ImageUrlCharacterUpdateRequest) *web.CharacterResponse
	FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *web.CharacterResponse
	FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []web.CharacterResponse
	FindAllByGachaSystemIdAndUserId(ctx context.Context, gachaSystemId int, userId int) []web.CharacterResponse
	Delete(ctx context.Context, id int, gachaSystemId int)
}

type CharacterServiceImpl struct {
	CharacterRepository   repository.CharacterRepository
	RarityRepository      repository.RarityRepository
	GachaSystemRepository repository.GachaSystemRepository
	Validate              *validator.Validate
}

func NewCharacterService(
	characterRepository repository.CharacterRepository,
	rarityRepository repository.RarityRepository,
	gachaSystemRepository repository.GachaSystemRepository,
	validate *validator.Validate,
) CharacterService {
	return &CharacterServiceImpl{
		CharacterRepository:   characterRepository,
		RarityRepository:      rarityRepository,
		GachaSystemRepository: gachaSystemRepository,
		Validate:              validate,
	}
}

func (service *CharacterServiceImpl) Create(ctx context.Context, request *web.CharacterCreateRequest) *web.CharacterResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, request.GachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	rarity := service.RarityRepository.FindByIdAndGachaSystemId(ctx, request.RarityId, request.GachaSystemId)
	if rarity == nil {
		panic(exception.NewNotFoundError(helper.ErrRarityNotFound))
	}

	existingCharacter := service.CharacterRepository.FindByNameAndGachaSystemId(ctx, request.Name, request.GachaSystemId)
	if existingCharacter != nil {
		panic(exception.NewConflictError("Character with the same name already exists"))
	}

	character := domain.Character{
		Name:          request.Name,
		GachaSystemId: request.GachaSystemId,
		RarityId:      request.RarityId,
	}

	service.CharacterRepository.Save(ctx, &character)

	characterResponse := web.ToCharacterResponse(&character)

	return characterResponse
}

func (service *CharacterServiceImpl) FindAllByGachaSystemId(ctx context.Context, gachaSystemId int) []web.CharacterResponse {
	userId := helper.ExtractUserID(ctx)

	return service.FindAllByGachaSystemIdAndUserId(ctx, gachaSystemId, userId)
}

func (service *CharacterServiceImpl) FindAllByGachaSystemIdAndUserId(ctx context.Context, gachaSystemId int, userId int) []web.CharacterResponse {
	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	characters := service.CharacterRepository.FindAllByGachaSystemId(ctx, gachaSystemId)
	if characters == nil {
		return nil
	}
	return web.ToCharacterResponses(characters)
}

func (service *CharacterServiceImpl) Update(ctx context.Context, request *web.CharacterUpdateRequest) *web.CharacterResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, request.GachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	character := service.CharacterRepository.FindByIdAndGachaSystemId(ctx, request.Id, request.GachaSystemId)
	if character == nil {
		panic(exception.NewNotFoundError("Character not found"))
	}

	request.UpdateCharacter(character)

	rarity := service.RarityRepository.FindByIdAndGachaSystemId(ctx, character.RarityId, character.GachaSystemId)
	if rarity == nil {
		panic(exception.NewNotFoundError(helper.ErrRarityNotFound))
	}

	existingCharacter := service.CharacterRepository.FindByNameAndGachaSystemId(ctx, character.Name, character.GachaSystemId)
	if existingCharacter != nil && existingCharacter.Id != request.Id {
		panic(exception.NewConflictError("Character with the same name already exists"))
	}

	service.CharacterRepository.Update(ctx, character)

	characterResponse := web.ToCharacterResponse(character)

	return characterResponse
}

func (service *CharacterServiceImpl) InsertImageUrl(ctx context.Context, request *web.ImageUrlCharacterUpdateRequest) *web.CharacterResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, request.GachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	character := service.CharacterRepository.FindByIdAndGachaSystemId(ctx, request.Id, request.GachaSystemId)
	if character == nil {
		panic(exception.NewNotFoundError("Character not found"))
	}

	character.ImageUrl = request.ImageUrl

	service.CharacterRepository.InsertImageUrl(ctx, character)

	characterResponse := web.ToCharacterResponse(character)

	return characterResponse
}
func (service *CharacterServiceImpl) FindByIdAndGachaSystemId(ctx context.Context, id int, gachaSystemId int) *web.CharacterResponse {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	character := service.CharacterRepository.FindByIdAndGachaSystemId(ctx, id, gachaSystemId)
	if character == nil {
		panic(exception.NewNotFoundError(helper.ErrCharacterNotFound))
	}

	return web.ToCharacterResponse(character)
}

func (service *CharacterServiceImpl) Delete(ctx context.Context, id int, gachaSystemId int) {
	userId := helper.ExtractUserID(ctx)

	gachaSystem := service.GachaSystemRepository.FindByIdAndUserId(ctx, gachaSystemId, userId)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError(helper.ErrGachaSystemNotFound))
	}

	character := service.CharacterRepository.FindByIdAndGachaSystemId(ctx, id, gachaSystemId)
	if character == nil {
		panic(exception.NewNotFoundError(helper.ErrCharacterNotFound))
	}

	service.CharacterRepository.Delete(ctx, id, gachaSystemId)
}
