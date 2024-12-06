import { Character } from "./characterType";
import { Rarity } from "./rarityType";

export type GachaSystem = {
  id : number,
  name : string,
};

export type GachaSystemEndpoint = {
  id : number,
  endpoint : string,
  totalCharacters: number,
  message? : string
};

export type GachaSytemDetail = {
  id: number,
  name: string,
  endpoint: string,
  rarities: Rarity[],
  characters: Character[]
  message? : string
}
  