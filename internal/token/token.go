package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type ApiResponse struct {
	IamToken  string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Token struct {
	token   string
	expired time.Time
	mu      sync.Mutex
}

const (
	yaCloudTokenApiUrl      = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
	yaCloudTokenApiOauthKey = "yandexPassportOauthToken"
)

func getIamTokenFromApi() (ApiResponse, error) {
	requestData := map[string]string{yaCloudTokenApiOauthKey: os.Getenv("YANDEX_OAUTH_TOKEN")}

	requestJsonData, _ := json.Marshal(requestData)

	resp, err := http.Post(yaCloudTokenApiUrl, "application/json", bytes.NewBuffer(requestJsonData))

	if err != nil {
		return ApiResponse{}, err
	}

	if resp.StatusCode != 200 {
		return ApiResponse{}, errors.New("iam token could not be retrieved")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return ApiResponse{}, err
	}

	var tokenResponse ApiResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return ApiResponse{}, err
	}

	return tokenResponse, nil
}

func (t *Token) Get() (string, error) {
	durationUntilExpired := time.Until(t.expired)
	if t.token == "" || durationUntilExpired <= 6*time.Hour {
		apiResponse, err := getIamTokenFromApi()
		if err != nil {
			return "", err
		}
		t.mu.Lock()
		t.token = apiResponse.IamToken
		t.expired = apiResponse.ExpiresAt
		t.mu.Unlock()
		return t.token, nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.token, nil
}
