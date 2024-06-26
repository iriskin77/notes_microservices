package auth

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID     string `json:"UUID,omitempty"`
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

type CreateUserDTO struct {
	Email          string `json:"Email,omitempty"`
	Password       string `json:"Password,omitempty"`
	RepeatPassword string `json:"RepeatPassword,omitempty"`
}

func NewUser(dto CreateUserDTO) User {
	return User{
		UUID:     uuid.New().String(),
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type LoginUserDTO struct {
	Email    string
	Password string
}

func LoginUser(dto LoginUserDTO) User {
	return User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (u *User) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}

type RefreshTokensInput struct {
	Token string
}

// func (u *User) CheckPassword(password string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
// 	if err != nil {
// 		return fmt.Errorf("password does not match")
// 	}
// 	return nil
// }
