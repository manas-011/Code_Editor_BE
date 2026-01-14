package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type hfRequest struct {
	Inputs string `json:"inputs"`
}

func CallHuggingFace(prompt string) (string, error) {

	apiKey := os.Getenv("HF_API_KEY")
	model := os.Getenv("HF_MODEL")

	if apiKey == "" || model == "" {
		return "", errors.New("HF_API_KEY or HF_MODEL missing")
	}

	// URL-encode model path (VERY IMPORTANT)
	modelPath := url.PathEscape(model)

	url := "https://router.huggingface.co/hf-inference/models/" + modelPath

	payload, _ := json.Marshal(hfRequest{
		Inputs: prompt,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(
			"HuggingFace error: " + string(body),
		)
	}

	var result []map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", errors.New("empty response from HuggingFace")
	}

	return result[0]["generated_text"].(string), nil
}
