package service

import (
	"log"
)

func SendOTPEmail(email, otp string) {
	log.Printf("Sending OTP %s to %s\n", otp, email)
}