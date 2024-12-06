package controller

import (
	"errors"
	"fmt"
	"gacha-master/exception"
	"gacha-master/helper"
	"gacha-master/model/web"
	"gacha-master/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CharacterController interface {
	Create(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
	Delete(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
}

type CharacterControllerImpl struct {
	CharacterService     service.CharacterService
	ImageUploaderService service.UploaderService
}

func NewCharacterController(rarityService service.CharacterService, imageUploader service.UploaderService) CharacterController {
	return &CharacterControllerImpl{
		CharacterService:     rarityService,
		ImageUploaderService: imageUploader,
	}
}

func (controller *CharacterControllerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	validateFileTypeAndSize(request)

	imageUploadRequest := web.ToImageCharacterUploadRequest(request)
	if imageUploadRequest == nil {
		panic(exception.NewBadRequestError("Image file is required"))
	}

	characterCreateRequest := web.ToCharacterCreateRequest(request)
	characterCreateResponse := controller.CharacterService.Create(request.Context(), characterCreateRequest)

	// Upload Image
	imageUploadRequest.Id = characterCreateResponse.Id
	imageUrl := controller.ImageUploaderService.UploadCharacterImage(request.Context(), imageUploadRequest)

	// Set image url
	imageUrlUpdateRequest := &web.ImageUrlCharacterUpdateRequest{
		Id:            characterCreateResponse.Id,
		GachaSystemId: characterCreateRequest.GachaSystemId,
		ImageUrl:      imageUrl,
	}
	characterCreateResponse = controller.CharacterService.InsertImageUrl(request.Context(), imageUrlUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   characterCreateResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CharacterControllerImpl) Update(writer http.ResponseWriter, request *http.Request) {
	validateFileTypeAndSize(request)

	characterUpdateRequest := web.ToCharacterUpdateRequest(request)
	imageUploadRequest := web.ToImageCharacterUploadRequest(request)

	if imageUploadRequest != nil {
		imageUrl := controller.ImageUploaderService.UploadCharacterImage(request.Context(), imageUploadRequest)
		characterUpdateRequest.ImageUrl = imageUrl
	}

	characterUpdateResponse := controller.CharacterService.Update(request.Context(), characterUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   characterUpdateResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CharacterControllerImpl) Delete(writer http.ResponseWriter, request *http.Request) {
	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")
	gachaSystemId, _ := strconv.Atoi(gachaSystemIdStr)

	characterIdStr := chi.URLParam(request, "characterId")
	characterId, _ := strconv.Atoi(characterIdStr)

	controller.CharacterService.Delete(request.Context(), characterId, gachaSystemId)
	controller.ImageUploaderService.DeleteCharacterImage(request.Context(), characterId, gachaSystemId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Character with ID %d successfully deleted", characterId),
		},
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CharacterControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	gachaSystemIdStr := chi.URLParam(request, "gachaSystemId")
	gachaSystemId, _ := strconv.Atoi(gachaSystemIdStr)

	characterIdStr := chi.URLParam(request, "characterId")
	characterId, _ := strconv.Atoi(characterIdStr)

	characterResponse := controller.CharacterService.FindByIdAndGachaSystemId(request.Context(), characterId, gachaSystemId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   characterResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func validateFileTypeAndSize(request *http.Request) {
	// Convert max size to bytes
	maxSizeMB := 2.2
	maxSizeBytes := int64(maxSizeMB * 1024 * 1024)
	// Set the maximum allowed size for the request body
	request.Body = http.MaxBytesReader(nil, request.Body, maxSizeBytes) // Limit size

	// Parse the form data, check for size errors
	if err := request.ParseForm(); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			panic(exception.NewBadRequestError("File size exceeds the limit (" + fmt.Sprintf("%.2f", maxSizeMB) + "MB)"))
		}
		panic(exception.NewBadRequestError(err.Error()))
	}

	// Parse multipart form with the same max size
	err := request.ParseMultipartForm(maxSizeBytes)
	if err != nil {
		helper.PanicIfError(err, "Failed to parse form data")
	}
}
