export type Character = {
  id : number,
  name : string,
  imageUrl : string,
  rarityId : number,
  message? : string
}

export type CharacterInput = {
  id : number,
  name : string,
  image : File | null,
  rarityId : number,
  gachaSystemId : number
}

export type CharacterGacha = {
  id : number,
  name : string,
  imageUrl : string,
  rarity : string,
  message? : string
}
  
  