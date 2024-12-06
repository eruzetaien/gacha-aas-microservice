package web

import "gacha-master/model/domain"

type RarityResponse struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Chance float32 `json:"chance"`
}

func ToRarityResponse(rarity *domain.Rarity) *RarityResponse {
	return &RarityResponse{
		Id:     rarity.Id,
		Name:   rarity.Name,
		Chance: rarity.Chance,
	}
}

func ToRaritiesResponse(rarities []domain.Rarity) []RarityResponse {
	var rarityResponses []RarityResponse
	for _, rarity := range rarities {
		rarityResponse := ToRarityResponse(&rarity)
		rarityResponses = append(rarityResponses, *rarityResponse)
	}

	return rarityResponses
}
