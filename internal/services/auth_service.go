package services

type AuthServiceI interface {
	Authenticate(username, pass string) bool
}

type AuthService struct {
}

func NewAuthService() AuthService {
	return AuthService{}
}

const AuthUser = "user"
const AuthPass = "pass"

func (s AuthService) Authenticate(username string, pass string) bool {
	return username == AuthUser && pass == AuthPass
}
