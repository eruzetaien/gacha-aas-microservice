package controller

import (
	"fmt"
	"gacha-master/helper"
	"gacha-master/model/web"
	"gacha-master/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type RarityController interface {
	Create(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
	Delete(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
}

type RarityControllerImpl struct {
	RarityService service.RarityService
}

func NewRarityController(rarityService service.RarityService) RarityController {
	return &RarityControllerImpl{
		RarityService: rarityService,
	}
}

func (controller *RarityControllerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	rarityCreateRequest := web.RarityCreateRequest{}
	helper.ReadFromRequestBody(request, &rarityCreateRequest)

	rarityResponse := controller.RarityService.Create(request.Context(), &rarityCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   rarityResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RarityControllerImpl) Update(writer http.ResponseWriter, request *http.Request) {
	rarityUpdateRequest := web.RarityUpdateRequest{}
	helper.ReadFromRequestBody(request, &rarityUpdateRequest)

	rarityResponse := controller.RarityService.Update(request.Context(), &rarityUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   rarityResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RarityControllerImpl) Delete(writer http.ResponseWriter, request *http.Request) {
	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")
	gachaSystemId, _ := strconv.Atoi(gachaSystemIdStr)

	characterIdStr := chi.URLParam(request, "rarityId")
	characterId, _ := strconv.Atoi(characterIdStr)

	controller.RarityService.Delete(request.Context(), characterId, gachaSystemId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Rarity with ID %d successfully deleted", characterId),
		},
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RarityControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")
	gachaSystemId, _ := strconv.Atoi(gachaSystemIdStr)

	raritiesResponse := controller.RarityService.FindAllByGachaSystemId(request.Context(), gachaSystemId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   raritiesResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
