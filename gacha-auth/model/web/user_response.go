package web

type UserResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserToken string `json:"userToken"`
}
