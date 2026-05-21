package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Signup(req dto.SignupRequest) (*models.User, error) {

	existingUser, _ := repositories.GetUserByPhoneNumber(req.PhoneNumber)
	if existingUser != nil {
		return nil, errors.New("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(hashedPassword),
	}

	err = repositories.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Login(req dto.LoginRequest) (string, error) {

	user, err := repositories.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", errors.New("invalid password")
	}

	session, err := repositories.CreateSession(user.ID)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"session_id": session.ID,
	})

	return token.SignedString(jwtSecret)
}
