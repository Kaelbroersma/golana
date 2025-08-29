package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func FetchTokenPrice(contract string) (TokenPriceInfo, error) {
	endpoint := fmt.Sprintf("https://lite-api.jup.ag/price/v3?ids=%v", contract)

	u, err := url.Parse(endpoint)
	if err != nil {
		return TokenPriceInfo{}, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return TokenPriceInfo{}, err
	}

	if responseCode := resp.StatusCode; responseCode != 200 {
		return TokenPriceInfo{}, fmt.Errorf("received response code %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenPriceInfo{}, err
	}
	defer resp.Body.Close()

	var token TokenPriceResponse

	err = json.Unmarshal(body, &token)
	if err != nil {
		return TokenPriceInfo{}, err
	}

	return token[contract], nil
}

func FetchTokenData(contract string, apiKey string) (TokenData, error) {
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
		return TokenData{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return TokenData{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return TokenData{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenData{}, err
	}
	defer resp.Body.Close()

	var token TokenData
	err = json.Unmarshal(body, &token)
	if err != nil {
		return TokenData{}, err
	}

	fmt.Printf("Successfully data for token: %v with CA: %v", token.Result.TokenInfo.Symbol, contract)

	return token, nil
}

func FetchTokenSocket(contract string, apiKey string) (TokenData, error) {

	return TokenData{}, nil
}
