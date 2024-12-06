package controller

import (
	"fmt"
	"gacha-pull/helper"
	"gacha-pull/model/web"
	"gacha-pull/service"
	"log"
	"net/http"
	"strings"
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
	endpoint := request.URL.String()
	log.Println(request.Host)
	log.Println(request.Proto)
	if request.URL.Scheme == "" {
		endpoint = fmt.Sprintf("%s://%s%s", strings.ToLower(request.Proto[:len(request.Proto)-4]), request.Host, request.URL.RequestURI())
	}

	selectedCharacter := controller.CharacterService.Pull(request.Context(), endpoint)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   selectedCharacter,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
