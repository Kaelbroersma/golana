package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Kaelbroersma/golana/internal/database"
	"github.com/Kaelbroersma/golana/internal/server/auth"
	"github.com/Kaelbroersma/golana/internal/server/tokens"
	"github.com/google/uuid"
)

type CreateTradeRequest struct {
	Contract    string  `json:"contract"`
	Side        string  `json:"side"`
	AmountInUSD float64 `json:"amount_usd"`
}

func (cfg *ServerConfig) handleCreateTrade(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.TokenSecret)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized: unable to validate JWT", err)
		return
	}

	user, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized: unable to find player in userbase", err)
		return
	}

	req := CreateTradeRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}

	token, err := tokens.FetchToken(req.Contract, cfg.HeliusAPIKey)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch token", err)
		return
	}

	_, err = cfg.DB.GetContractByID(r.Context(), req.Contract)
	if err != nil {
		_, err = cfg.DB.CreateContract(r.Context(), database.CreateContractParams{
			ContractID: token.Result.ID,
			Name:       token.Result.TokenInfo.Symbol,
		})
		if err != nil {
			RespondWithError(w, 500, "Failed to create contract", err)
			return
		}
	}

	tokenPrice := token.Result.TokenInfo.PriceInfo.PricePerToken
	tradeQuantity := req.AmountInUSD / tokenPrice
	fmt.Printf("buying power: %v, token price: %v, quantity to buy: %v\n", req.AmountInUSD, tokenPrice, tradeQuantity)

	if user.Balance < req.AmountInUSD {
		RespondWithError(w, 402, "Insufficient balance", err)
		return
	}

	fmt.Println("We made it to the switch statement")

	switch req.Side {
	case "long":
		_, err = cfg.DB.UpdateUserBalance(r.Context(), database.UpdateUserBalanceParams{
			Balance: user.Balance - req.AmountInUSD,
			ID:      userID,
		})
		if err != nil {
			RespondWithError(w, 500, "Failed to update user balance", err)
			return
		}

		trade, err := cfg.DB.CreateTrade(r.Context(), database.CreateTradeParams{
			ID:          uuid.New().String(),
			UserID:      userID,
			Contract:    req.Contract,
			Side:        req.Side,
			Quantity:    tradeQuantity,
			BoughtPrice: sql.NullFloat64{Float64: tokenPrice, Valid: true},
			SoldPrice:   sql.NullFloat64{Float64: 0, Valid: false},
		})
		if err != nil {
			RespondWithError(w, 500, "Failed to create trade", err)
			return
		}
		RespondWithJSON(w, http.StatusOK, trade)

	case "short":
		_, err = cfg.DB.UpdateUserBalance(r.Context(), database.UpdateUserBalanceParams{
			Balance: user.Balance + req.AmountInUSD,
			ID:      userID,
		})
		if err != nil {
			RespondWithError(w, 500, "Failed to update user balance", err)
			return
		}

		trade, err := cfg.DB.CreateTrade(r.Context(), database.CreateTradeParams{
			ID:          uuid.New().String(),
			UserID:      userID,
			Contract:    req.Contract,
			Side:        req.Side,
			Quantity:    tradeQuantity,
			BoughtPrice: sql.NullFloat64{Float64: 0, Valid: false},
			SoldPrice:   sql.NullFloat64{Float64: tokenPrice, Valid: true},
		})
		if err != nil {
			RespondWithError(w, 500, "Failed to create trade", err)
			return
		}
		RespondWithJSON(w, http.StatusOK, trade)

	default:
		RespondWithError(w, 400, "Invalid side", errors.New("pick between long or short"))
		return
	}
}
