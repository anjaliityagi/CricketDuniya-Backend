package services

import (
	"CricketDuniya-Backend/internal/repositories"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendForgotPasswordOTP(phone string) (string, error) {

	otp := generateOTP()

	err := repositories.SaveOTP(phone, otp)

	if err != nil {
		return "", err
	}

	// TEMPORARY
	// later replace with SMS provider

	return otp, nil
}

func VerifyOTPAndResetPassword(
	phone string,
	otp string,
	newPassword string,
) error {

	valid, err := repositories.VerifyOTP(phone, otp)

	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("invalid or expired otp")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	err = repositories.UpdatePassword(
		phone,
		string(hashedPassword),
	)

	if err != nil {
		return err
	}

	err = repositories.MarkOTPUsed(phone, otp)

	if err != nil {
		return err
	}

	return nil
}
