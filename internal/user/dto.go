package user

type UserLogin struct {
	Login    string `json:"login" valid:",required"`
	Password string `json:"password" valid:",required"`
}
