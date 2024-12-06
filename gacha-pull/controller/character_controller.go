package controller

import (
	"gacha-pull/helper"
	"gacha-pull/model/web"
	"gacha-pull/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CharacterController interface {
	Pull(writer http.ResponseWriter, request *http.Request)
}

type CharacterControllerImpl struct {
	CharacterService service.CharacterService
}

func NewCharacterController(rarityService service.CharacterService) CharacterController {
	return &CharacterControllerImpl{
		CharacterService: rarityService,
	}
}

func (controller *CharacterControllerImpl) Pull(writer http.ResponseWriter, request *http.Request) {
	endpointId := chi.URLParam(request, "endpointId")

	selectedCharacter := controller.CharacterService.Pull(request.Context(), endpointId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   selectedCharacter,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
