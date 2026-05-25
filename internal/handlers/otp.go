package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForgotPassword(c *gin.Context) {

	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
		})

		return
	}

	otp, err := services.SendForgotPasswordOTP(req.Phone)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to generate otp",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "otp generated successfully",

		// TEMP ONLY
		"otp": otp,
	})
}

func VerifyOTPAndResetPassword(c *gin.Context) {

	var req dto.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
		})

		return
	}

	err := services.VerifyOTPAndResetPassword(
		req.Phone,
		req.OTP,
		req.NewPassword,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "password updated successfully",
	})
}
