package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const hfURL = "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.2"

type HFRequest struct {
	Inputs string `json:"inputs"`
}

type HFResponse []struct {
	GeneratedText string `json:"generated_text"`
}

func AskAI(question string) (string, error) {
	body, _ := json.Marshal(HFRequest{
		Inputs: question,
	})

	req, _ := http.NewRequest("POST", hfURL, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("HF_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI API failed with status %d", resp.StatusCode)
	}

	var result HFResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result[0].GeneratedText, nil
}
