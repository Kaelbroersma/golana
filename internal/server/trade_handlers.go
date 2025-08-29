package server

import (
	"database/sql"
	"encoding/json"
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

func (cfg *Config) handleCreateTrade(w http.ResponseWriter, r *http.Request) {
	user, err := auth.AuthenticateWithBearer(r, cfg.TokenSecret, cfg.DB)
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

	token, err := tokens.FetchTokenData(req.Contract, cfg.HeliusAPIKey)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch token", err)
		return
	}

	liveTokenPrice, err := tokens.FetchTokenPrice(req.Contract)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch token price", err)
		return
	}

	// Should this vvvv be handled in fetch token or create trade?

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

	tokenPrice := liveTokenPrice.UsdPrice
	tradeQuantity := req.Quantity

	if user.BuyingPower < (tokenPrice * math.Abs(tradeQuantity)) {
		RespondWithError(w, 422, "Insufficient buying power", nil)
		return
	}

	updatedBuyingPower := user.BuyingPower - (tokenPrice * math.Abs(tradeQuantity))

	trade, err := cfg.DB.CreateTrade(r.Context(), database.CreateTradeParams{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		Contract:     req.Contract,
		OpenQuantity: tradeQuantity,
		OpenPrice:    sql.NullFloat64{Float64: tokenPrice, Valid: true},
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

func (cfg *Config) handleCloseTrade(w http.ResponseWriter, r *http.Request) {
	var averageSoldPrice float64
	var profit float64

	user, err := auth.AuthenticateWithBearer(r, cfg.TokenSecret, cfg.DB)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized", err)
		return
	}

	req := CloseTradeRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}

	trade, err := cfg.DB.GetTrade(r.Context(), req.Trade)
	if err != nil {
		RespondWithError(w, 500, "Failed to get trade", err)
		return
	}

	if trade.UserID != user.ID {
		RespondWithError(w, 403, "Insufficient user", nil)
		return
	}

	_, err = cfg.DB.GetContractByID(r.Context(), trade.Contract)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch contract- is your position open?", err)
		return
	}

	contract, err := tokens.FetchTokenPrice(trade.Contract)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch contract", err)
		return
	}

	percentToSell := req.Percent / 100
	contractPrice := contract.UsdPrice
	closedQuantity := trade.ClosedQuantity.Float64 + (trade.OpenQuantity * percentToSell)
	openQuantity := trade.OpenQuantity - (trade.OpenQuantity * percentToSell)

	if trade.ClosedQuantity.Valid {
		averageSoldPrice = ((trade.AverageClosePrice.Float64 * trade.ClosedQuantity.Float64) + ((trade.OpenQuantity * percentToSell) * contractPrice)) / closedQuantity
	} else {
		averageSoldPrice = contractPrice
	}

	if trade.RealizedProfit.Valid {
		profit = trade.RealizedProfit.Float64 + (closedQuantity * (contractPrice - trade.OpenPrice.Float64))
	} else {
		profit = closedQuantity * (contractPrice - trade.OpenPrice.Float64)
	}

	closedTrade, err := cfg.DB.UpdateTrade(r.Context(), database.UpdateTradeParams{
		ID:                req.Trade,
		OpenQuantity:      openQuantity,
		ClosedQuantity:    sql.NullFloat64{Float64: closedQuantity, Valid: true},
		AverageClosePrice: sql.NullFloat64{Float64: averageSoldPrice, Valid: true},
		RealizedProfit:    sql.NullFloat64{Float64: profit, Valid: true},
	})
	if err != nil {
		RespondWithError(w, 500, "Failed to close trade", err)
		return
	}

	RespondWithJSON(w, 200, closedTrade)
}

type ListOfTrades map[string][]TradeResponse

type TradeResponse struct {
	Trade            string  `json:"trade"`
	Contract         string  `json:"contract"`
	OpenQuantity     float64 `json:"open_quantity"`
	ClosedQuantity   float64 `json:"closed_quantity"`
	UnrealizedProfit float64 `json:"unrealized_profit"`
	RealizedProfit   float64 `json:"realized_profit"`
}

func (cfg *Config) handleGetTrades(w http.ResponseWriter, r *http.Request) {
	tradeList := make(ListOfTrades)
	user, err := auth.AuthenticateWithBearer(r, cfg.TokenSecret, cfg.DB)
	if err != nil {
		RespondWithError(w, 401, "Could not find user. How are we supposed to display your trades?", err)
		return
	}

	userTrades, err := cfg.DB.GetTradesForUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to get open trades for user %v", user.Name), err)
		return
	}

	if len(userTrades) == 0 {
		RespondWithJSON(w, 200, tradeList)
		return
	}

	for _, trade := range userTrades {
		var status string
		if trade.OpenQuantity > 0 {
			status = "Open"
		}
		if trade.OpenQuantity == 0 {
			status = "Closed"
		}

		tradeList[status] = append(tradeList[status], TradeResponse{
			Trade:            trade.ID,
			Contract:         trade.Contract,
			OpenQuantity:     trade.OpenQuantity,
			ClosedQuantity:   trade.ClosedQuantity.Float64,
			UnrealizedProfit: trade.UnrealizedProfit.Float64,
			RealizedProfit:   trade.RealizedProfit.Float64,
		})
	}

	RespondWithJSON(w, 200, tradeList)

}
