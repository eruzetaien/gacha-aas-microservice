package web

import (
	"gacha-master/model/domain"
)

type RarityCreateRequest struct {
	GachaSystemId int     `json:"gachaSystemId" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Chance        float32 `json:"chance" validate:"required,gte=0,lte=100"`
}

type RarityUpdateRequest struct {
	Id            int     `json:"id" validate:"required"`
	GachaSystemId int     `json:"gachaSystemId" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Chance        float32 `json:"chance" validate:"required,gte=0,lte=100"`
}

func (updateRequest *RarityUpdateRequest) UpdateRarity(rarity *domain.Rarity) {
	rarity.Name = updateRequest.Name
	rarity.Id = updateRequest.Id
	rarity.Chance = updateRequest.Chance
}
