package Infrastructure

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) bool
}

type passwordService struct{}

func NewPasswordService() PasswordService { return &passwordService{} }

func (p *passwordService) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (p *passwordService) Compare(hashed string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
