package web

import "gacha-master/model/domain"

type GachaSystemDetailResponse struct {
	Id         int                 `json:"id"`
	Name       string              `json:"name"`
	Endpoint   string              `json:"endpoint"`
	Rarities   []RarityResponse    `json:"rarities"`
	Characters []CharacterResponse `json:"characters"`
}

type GachaSystemResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GachaSystemEndpointResponse struct {
	Id              int    `json:"id"`
	TotalCharacters int    `json:"totalCharacters"`
	Endpoint        string `json:"endpoint"`
}

func ToGachaSystemResponse(gachaSystem *domain.GachaSystem) *GachaSystemResponse {
	return &GachaSystemResponse{
		Id:   gachaSystem.Id,
		Name: gachaSystem.Name,
	}
}

func ToGachaSystemsResponse(rarities []domain.GachaSystem) []GachaSystemResponse {
	var gachaSystemResponses []GachaSystemResponse
	for _, gachaSystem := range rarities {
		gachaSystemResponse := ToGachaSystemResponse(&gachaSystem)
		gachaSystemResponses = append(gachaSystemResponses, *gachaSystemResponse)
	}

	return gachaSystemResponses
}
