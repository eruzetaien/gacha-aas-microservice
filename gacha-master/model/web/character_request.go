package web

import (
	"gacha-master/exception"
	"gacha-master/model/domain"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type CharacterCreateRequest struct {
	GachaSystemId int    `form:"gachaSystemId" validate:"required"`
	Name          string `form:"name" validate:"required"`
	RarityId      int    `form:"rarityId" validate:"required"`
}

type ImageCharacterUploadRequest struct {
	Id             int            `validate:"required"`
	CharacterImage multipart.File `validate:"required"`
	GachaSystemId  int            `validate:"required"`
}

type ImageUrlCharacterUpdateRequest struct {
	Id            int    `validate:"required"`
	GachaSystemId int    `validate:"required"`
	ImageUrl      string `validate:"required"`
}

type CharacterUpdateRequest struct {
	Id            int    `form:"id" validate:"required"`
	GachaSystemId int    `form:"gachaSystemId" validate:"required"`
	Name          string `form:"name"`
	RarityId      int    `form:"rarityId"`
	ImageUrl      string
}

func ToCharacterUpdateRequest(request *http.Request) *CharacterUpdateRequest {
	id, err := strconv.Atoi(request.FormValue("id"))
	gachaSystemId, _ := strconv.Atoi(request.FormValue("gachaSystemId"))

	rarityIdStr := request.FormValue("rarityId")
	rarityId, err := strconv.Atoi(rarityIdStr)
	if err != nil || rarityIdStr == "" {
		rarityId = -1
	}

	return &CharacterUpdateRequest{
		Id:            id,
		Name:          request.FormValue("name"),
		GachaSystemId: gachaSystemId,
		RarityId:      rarityId,
	}
}

func (updateRequest *CharacterUpdateRequest) UpdateCharacter(character *domain.Character) {
	if updateRequest.Name != "" {
		character.Name = updateRequest.Name
	}
	if updateRequest.RarityId != -1 {
		character.RarityId = updateRequest.RarityId
	}
	if updateRequest.ImageUrl != "" {
		character.ImageUrl = updateRequest.ImageUrl
	}
}

func ToImageCharacterUploadRequest(request *http.Request) *ImageCharacterUploadRequest {
	imageFile, fileHeader, err := request.FormFile("image")
	if err != nil || fileHeader == nil || imageFile == nil {
		log.Printf("Error while getting image file: %v", err)
		return nil
	}

	gachaSystemId, _ := strconv.Atoi(request.FormValue("gachaSystemId"))

	characterIdStr := request.FormValue("id")
	characterId, err := strconv.Atoi(characterIdStr)
	if err != nil || characterIdStr == "" {
		characterId = 0
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/png" {
		panic(exception.NewBadRequestError("only PNG images are allowed"))
	}

	return &ImageCharacterUploadRequest{
		Id:             characterId,
		CharacterImage: imageFile,
		GachaSystemId:  gachaSystemId,
	}
}

func ToCharacterCreateRequest(request *http.Request) *CharacterCreateRequest {
	gachaSystemId, _ := strconv.Atoi(request.FormValue("gachaSystemId"))
	rarityId, _ := strconv.Atoi(request.FormValue("rarityId"))

	return &CharacterCreateRequest{
		Name:          request.FormValue("name"),
		GachaSystemId: gachaSystemId,
		RarityId:      rarityId,
	}
}
