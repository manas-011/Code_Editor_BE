package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)



type sendEmailPayload struct {
	Sender struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}


func SendEmail(email string, otp string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	brevoURL := os.Getenv("BREVO_URL")

	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY not set")
	}

	payload := sendEmailPayload{
		Subject: "Your OTP Code",
		HtmlContent: fmt.Sprintf(`
			<h2>OTP Verification</h2>
			<p>Your OTP is:</p>
			<h1 style="letter-spacing:4px;">%s</h1>
			<p>This OTP is valid for 5 minutes.</p>
		`, otp),
	}

	payload.Sender.Email = "no-reply@yourdomain.com"
	payload.Sender.Name = "Your App"

	payload.To = append(payload.To, struct {
		Email string `json:"email"`
	}{
		Email: email,
	})

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", brevoURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send email, status: %s", resp.Status)
	}

	return nil
}
