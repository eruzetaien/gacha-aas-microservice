package web

type GachaSystemCreateRequest struct {
	Name string `json:"name" validate:"required"`
}
