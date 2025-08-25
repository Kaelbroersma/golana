package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetAssetRequest struct {
	JsonRPC string         `json:"jsonrpc"`
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Params  GetAssetParams `json:"params"`
}

type GetAssetParams struct {
	ID string `json:"id"`
}

func FetchToken(contract string, apiKey string) (Token, error) {
	url := fmt.Sprintf("https://mainnet.helius-rpc.com/?api-key=%v", apiKey)

	tokenReq := GetAssetRequest{
		JsonRPC: "2.0",
		ID:      "1",
		Method:  "getAsset",
		Params: GetAssetParams{
			ID: contract,
		},
	}

	reqData, err := json.Marshal(tokenReq)
	if err != nil {
		return Token{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return Token{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Token{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	fmt.Println(string(body))

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func FetchTokenSocket(contract string, apiKey string) (Token, error) {

	return Token{}, nil
}
