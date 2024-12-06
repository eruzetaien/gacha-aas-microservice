package web

import (
	"gacha-pull/model/domain"
)

type CharacterResponse struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Rarity   string `json:"rarity"`
}

func ToCharacterResponse(character *domain.Character, rarityName string) *CharacterResponse {
	return &CharacterResponse{
		Name:     character.Name,
		ImageUrl: character.ImageUrl,
		Rarity:   rarityName,
	}
}
