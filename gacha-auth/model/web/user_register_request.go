package web

type UserRegisterRequest struct {
	Name     string `validate:"required,min=3,max=100" json:"name"`
	Username string `validate:"required,lowercase" json:"username"`
	Password string `validate:"required,min=8" json:"password"`
}
