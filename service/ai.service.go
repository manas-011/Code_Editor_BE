package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type hfRequest struct {
	Inputs string `json:"inputs"`
}

type hfResponse []struct {
	GeneratedText string `json:"generated_text"`
}

func CallHuggingFace(prompt string) (string, error) {

	apiKey := os.Getenv("HF_API_KEY")
	model := os.Getenv("HF_MODEL")

	if apiKey == "" || model == "" {
		return "", errors.New("huggingface env vars not set")
	}

	url := "https://api-inference.huggingface.co/models/" + model

	body, _ := json.Marshal(hfRequest{
		Inputs: prompt,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 60 * time.Second, // HF models can be slow
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("huggingface api error")
	}

	var hfRes hfResponse
	if err := json.NewDecoder(resp.Body).Decode(&hfRes); err != nil {
		return "", err
	}

	if len(hfRes) == 0 {
		return "", errors.New("empty response from huggingface")
	}

	return hfRes[0].GeneratedText, nil
}
