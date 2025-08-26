package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/Kaelbroersma/golana/internal/database"
	"github.com/Kaelbroersma/golana/internal/server/auth"
	"github.com/Kaelbroersma/golana/internal/server/tokens"
	"github.com/google/uuid"
)

type CreateTradeRequest struct {
	Contract string  `json:"contract"`
	Quantity float64 `json:"quantity"`
}

func (cfg *ServerConfig) handleCreateTrade(w http.ResponseWriter, r *http.Request) {
	user, err := auth.AuthenticateWithBearer(r, cfg)
	if err != nil {
		RespondWithError(w, 401, "Unable to authenticate user", err)
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

	// Should this vvvv be handled in fetch token or create trade?? think about it...

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
	tradeQuantity := req.Quantity

	if user.BuyingPower < (tokenPrice * math.Abs(tradeQuantity)) {
		RespondWithError(w, 422, "Insufficient buying power", nil)
		return
	}

	updatedBuyingPower := user.BuyingPower - (tokenPrice * math.Abs(tradeQuantity))

	trade, err := cfg.DB.CreateTrade(r.Context(), database.CreateTradeParams{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		Contract:   req.Contract,
		Quantity:   tradeQuantity,
		OpenPrice:  sql.NullFloat64{Float64: tokenPrice, Valid: true},
		ClosePrice: sql.NullFloat64{},
	})
	if err != nil {
		RespondWithError(w, 500, "Failed to create trade", err)
		return
	}

	_, err = cfg.DB.UpdateUserBalances(r.Context(), database.UpdateUserBalancesParams{
		BuyingPower: updatedBuyingPower,
		Exposure:    user.Exposure + (tokenPrice * math.Abs(tradeQuantity)),
		ID:          user.ID,
	})
	if err != nil {
		RespondWithError(w, 500, "Failed to update user balances", err)
		return
	}

	RespondWithJSON(w, 200, trade)
}

type CloseTradeRequest struct {
	Trade   string  `json:"trade"`
	Percent float64 `json:"percent"`
}

func (cfg *ServerConfig) handleCloseTrade(w http.ResponseWriter, r *http.Request) {
	user, err := auth.AuthenticateWithBearer(r, cfg)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized", err)
		return
	}

	trade, err := cfg.DB.GetTrade(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, 500, "Failed to get trade", err)
		return
	}

	if trade.UserID != user.ID {
		RespondWithError(w, 403, "Insufficient user", nil)
		return
	}

	req := CloseTradeRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}

	_, err = cfg.DB.GetContractByID(r.Context(), req.Trade)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch contract- is your position open?", err)
		return
	}

	contract, err := tokens.FetchToken(trade.Contract, cfg.HeliusAPIKey)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch contract", err)
		return
	}

	//percentToSell := req.Percent / 100
	//contractPrice := contract.Result.TokenInfo.PriceInfo.PricePerToken
	//heldQuantity := trade.Quantity

}
