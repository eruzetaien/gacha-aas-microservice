package web

import (
	"gacha-master/model/domain"
)

type CharacterResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	RarityId int    `json:"rarityId"`
}

func ToCharacterResponse(character *domain.Character) *CharacterResponse {
	return &CharacterResponse{
		Id:       character.Id,
		Name:     character.Name,
		ImageUrl: character.ImageUrl,
		RarityId: character.RarityId,
	}
}

func ToCharacterResponses(characters []domain.Character) []CharacterResponse {
	var characterResponses []CharacterResponse
	for _, character := range characters {
		characterResponse := ToCharacterResponse(&character)
		characterResponses = append(characterResponses, *characterResponse)
	}
	return characterResponses
}
