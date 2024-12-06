package controller

import (
	"fmt"
	"gacha-master/exception"
	"gacha-master/helper"
	"gacha-master/model/web"
	"gacha-master/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"strconv"
)

type GachaSystemController interface {
	Create(writer http.ResponseWriter, request *http.Request)
	FindById(writer http.ResponseWriter, request *http.Request)
	Delete(writer http.ResponseWriter, request *http.Request)
	FindAll(writer http.ResponseWriter, request *http.Request)
	FindEndpointByNameAndUserId(writer http.ResponseWriter, request *http.Request)
}

type GachaSystemControllerImpl struct {
	GachaSystemService   service.GachaSystemService
	RarityService        service.RarityService
	CharacterService     service.CharacterService
	ImageUploaderService service.UploaderService
}

func NewGachaSystemController(
	gachaSystemService service.GachaSystemService,
	rarityService service.RarityService,
	characterService service.CharacterService,
	imageUploaderService service.UploaderService,

) GachaSystemController {
	return &GachaSystemControllerImpl{
		GachaSystemService:   gachaSystemService,
		RarityService:        rarityService,
		CharacterService:     characterService,
		ImageUploaderService: imageUploaderService,
	}
}

func (controller *GachaSystemControllerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	gachaSystemCreateRequest := web.GachaSystemCreateRequest{}
	helper.ReadFromRequestBody(request, &gachaSystemCreateRequest)

	gachaSystemResponse := controller.GachaSystemService.Create(request.Context(), &gachaSystemCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   gachaSystemResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *GachaSystemControllerImpl) FindById(writer http.ResponseWriter, request *http.Request) {

	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")

	gachaSystemId, err := strconv.Atoi(gachaSystemIdStr)
	if err != nil {
		panic(exception.NewBadRequestError("Invalid gacha system id"))
	}

	gachaSystemRarities := controller.RarityService.FindAllByGachaSystemId(request.Context(), gachaSystemId)
	gachaSyeemCharacters := controller.CharacterService.FindAllByGachaSystemId(request.Context(), gachaSystemId)

	gachaSystemResponse := controller.GachaSystemService.FindById(request.Context(), gachaSystemId)
	gachaSystemResponse.Characters = gachaSyeemCharacters
	gachaSystemResponse.Rarities = gachaSystemRarities

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   gachaSystemResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *GachaSystemControllerImpl) Delete(writer http.ResponseWriter, request *http.Request) {
	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")
	gachaSystemId, _ := strconv.Atoi(gachaSystemIdStr)

	gachaSyeemCharacters := controller.CharacterService.FindAllByGachaSystemId(request.Context(), gachaSystemId)

	if len(gachaSyeemCharacters) != 0 {
		for _, character := range gachaSyeemCharacters {
			if character.ImageUrl != "" {
				controller.ImageUploaderService.DeleteCharacterImage(request.Context(), character.Id, gachaSystemId)
			}
		}

		controller.ImageUploaderService.DeleteGachaSystemCharacterImage(request.Context(), gachaSystemId)
	}

	controller.GachaSystemService.Delete(request.Context(), gachaSystemId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Gacha system with ID %d successfully deleted", gachaSystemId),
		},
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *GachaSystemControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request) {
	userGachaSytems := controller.GachaSystemService.FindAllByUserId(request.Context())

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userGachaSytems,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *GachaSystemControllerImpl) FindEndpointByNameAndUserId(writer http.ResponseWriter, request *http.Request) {
	gachaSystemName := request.URL.Query().Get("name")
	gachaSystemName, _ = url.QueryUnescape(gachaSystemName)

	userIdStr := chi.URLParam(request, "userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		panic(exception.NewBadRequestError("Invalid user id"))
	}

	gachaSystemResponse := controller.GachaSystemService.FindByNameAndUserId(request.Context(), gachaSystemName, userId)

	gachaSystemCharacters := controller.CharacterService.FindAllByGachaSystemIdAndUserId(request.Context(), gachaSystemResponse.Id, userId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: map[string]interface{}{
			"id":              gachaSystemResponse.Id,
			"endpoint":        gachaSystemResponse.Endpoint,
			"totalCharacters": len(gachaSystemCharacters),
		},
	}

	helper.WriteToResponseBody(writer, webResponse)
}
