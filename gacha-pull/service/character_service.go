package service

import (
	"context"
	"gacha-pull/exception"
	"gacha-pull/model/domain"
	"gacha-pull/model/web"
	"gacha-pull/repository"
	"log"
	"math/rand"
)

type CharacterService interface {
	Pull(ctx context.Context, endpoint string) *web.CharacterResponse
}

type CharacterServiceImpl struct {
	CharacterRepository   repository.CharacterRepository
	RarityRepository      repository.RarityRepository
	GachaSystemRepository repository.GachaSystemRepository
}

func NewCharacterService(
	characterRepository repository.CharacterRepository,
	rarityRepository repository.RarityRepository,
	gachaSystemRepository repository.GachaSystemRepository,
) CharacterService {
	return &CharacterServiceImpl{
		CharacterRepository:   characterRepository,
		RarityRepository:      rarityRepository,
		GachaSystemRepository: gachaSystemRepository,
	}
}

func (service *CharacterServiceImpl) Pull(ctx context.Context, endpoint string) *web.CharacterResponse {
	gachaSystem := service.GachaSystemRepository.FindByEndpoint(ctx, endpoint)
	if gachaSystem == nil {
		panic(exception.NewNotFoundError("Gacha system not found"))
	}

	rarities := service.RarityRepository.FindAllByGachaSystemId(ctx, gachaSystem.Id)
	characters := service.CharacterRepository.FindAllByGachaSystemId(ctx, gachaSystem.Id)

	if (rarities == nil || len(rarities) == 0) || (characters == nil || len(characters) == 0) {
		panic(exception.NewNotFoundError("Rarities or characters not found"))
	}

	rarityNameMap := make(map[int]string)
	for _, rarity := range rarities {
		rarityNameMap[rarity.Id] = rarity.Name
	}

	rarityCharsMap := make(map[int][]domain.Character)
	for _, character := range characters {
		rarityCharsMap[character.RarityId] = append(rarityCharsMap[character.RarityId], character)
	}

	var filteredRarities []domain.Rarity
	for _, rarity := range rarities {
		if len(rarityCharsMap[rarity.Id]) > 0 && rarity.Chance > 0 {
			filteredRarities = append(filteredRarities, rarity)
		}
	}

	if len(filteredRarities) == 0 {
		panic(exception.NewNotFoundError("Rarities or characters not found"))
	}

	selectedCharacter, rarities := SequentialRandomRarity(filteredRarities, rarityCharsMap)

	if selectedCharacter != nil {
		return web.ToCharacterResponse(selectedCharacter, rarityNameMap[selectedCharacter.RarityId])
	}

	log.Printf("Sequential random rarity failed, fallback to full random rarity")

	if rarities == nil || len(rarities) == 0 {
		rarities = filteredRarities
	}

	selectedCharacter = FullRandomRarity(rarities, rarityCharsMap)
	return web.ToCharacterResponse(selectedCharacter, rarityNameMap[selectedCharacter.RarityId])
}

func SequentialRandomRarity(rarities []domain.Rarity, rarityMap map[int][]domain.Character) (*domain.Character, []domain.Rarity) {
	selectedRarityIdx := -1
	totalChance := float32(100)

	for totalChance >= 0 && len(rarities) > 0 {
		boostChance := rand.Float32() * totalChance

		randomIndex := rand.Intn(len(rarities))
		selectedRarity := rarities[randomIndex]

		rarityTotalChance := selectedRarity.Chance + boostChance
		if rarityTotalChance < totalChance {
			if randomIndex == len(rarities)-1 {
				rarities = rarities[:randomIndex]
			} else {
				rarities = append(rarities[:randomIndex], rarities[randomIndex+1:]...)
			}
			continue
		}

		randomChance := rand.Float32() * 100
		if randomChance < rarityTotalChance {
			selectedRarityIdx = randomIndex
			break
		}

		totalChance -= rarityTotalChance

		if randomIndex == len(rarities)-1 {
			rarities = rarities[:randomIndex]
		} else {
			rarities = append(rarities[:randomIndex], rarities[randomIndex+1:]...)
		}
	}

	if selectedRarityIdx == -1 {
		return nil, rarities
	}

	selectedRarity := rarities[selectedRarityIdx]
	randomCharIndex := rand.Intn(len(rarityMap[selectedRarity.Id]))
	randomCharacter := rarityMap[selectedRarity.Id][randomCharIndex]
	return &randomCharacter, rarities
}

func FullRandomRarity(rarities []domain.Rarity, rarityMap map[int][]domain.Character) *domain.Character {
	selectedCharacters := []domain.Character{}

	for _, rarity := range rarities {
		charactersInRarity := rarityMap[rarity.Id]
		totalCharactersInRarity := len(charactersInRarity)

		randomChance := rand.Float32() * 100

		if randomChance < rarity.Chance {
			randomIndex := rand.Intn(totalCharactersInRarity)
			selectedCharacter := charactersInRarity[randomIndex]

			selectedCharacters = append(selectedCharacters, selectedCharacter)
		}
	}

	totalSelectedCharacters := len(selectedCharacters)
	if totalSelectedCharacters == 0 {
		randomIndex := rand.Intn(len(rarities))
		selectedRarity := rarities[randomIndex]

		totalCharactersInRarity := len(rarityMap[selectedRarity.Id])
		randomIndex = rand.Intn(totalCharactersInRarity)

		randomCharacter := rarityMap[selectedRarity.Id][randomIndex]
		return &randomCharacter
	}

	randomIndex := rand.Intn(totalSelectedCharacters)
	selectedCharacter := selectedCharacters[randomIndex]
	return &selectedCharacter
}
