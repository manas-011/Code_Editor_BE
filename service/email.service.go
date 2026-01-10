package service

import (
	"log"
)

func SendEmail(email, otp string) {
	log.Printf("Sending OTP %s to %s\n", otp, email)
}