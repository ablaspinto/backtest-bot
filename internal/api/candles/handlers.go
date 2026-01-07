package candles

import (
	"cmd/internal/json"
	reponsehandlers "cmd/internal/response_handlers"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateCandle(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[createCandleParams](w, r)
	if !ok {
		return
	}
	candle, err := h.service.CreateCandle(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, "error creating the candle")
		return
	}
	json.WriteSuccess(w, http.StatusCreated, candle)

}

func (h *handler) GetCandleByID(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	c, err := h.service.GetCandleByID(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting this candle")
	}
	json.WriteSuccess(w, http.StatusOK, c)
}

func (h *handler) GetCandle(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandleParams](w, r)
	if !ok {
		return
	}
	c, err := h.service.GetCandle(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting candle")
	}
	json.WriteSuccess(w, http.StatusOK, c)

}

func (h *handler) GetLatestCandle(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getLatestCandleParams](w, r)
	if !ok {
		return
	}
	c, err := h.service.GetLatestCandle(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting latest candle")
	}
	json.WriteSuccess(w, http.StatusOK, c)
}

func (h *handler) GetRecentCandles(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getRecentCandleParams](w, r)
	if !ok {
		return
	}
	c, err := h.service.GetRecentCandles(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting recent candles")
	}
	json.WriteSuccess(w, http.StatusOK, c)
}

func (h *handler) GetCandlesInRange(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandlesInRangeParams](w, r)
	if !ok {
		return
	}
	c, err := h.service.GetCandlesInRange(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting candles within range")
	}
	json.WriteSuccess(w, http.StatusOK, c)

}

func (h *handler) GetCandlesAfterTimestamp(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandlesAfterTimestampParams](w, r)
	if !ok {
		return
	}
	cAfter, err := h.service.GetCandlesAfterTimestamp(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting candles after timestamp")
	}
	json.WriteSuccess(w, http.StatusOK, cAfter)
}

func (h *handler) GetCandlesBeforeTimestamp(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandlesBeforeTimestampParams](w, r)
	if !ok {
		return
	}
	cBefore, err := h.service.GetCandlesBeforeTimestamp(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting getting candles before timestamp")
	}
	json.WriteSuccess(w, http.StatusOK, cBefore)

}

func (h *handler) ListCandlesPaginated(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[listCandlesPaginatedParams](w, r)
	if !ok {
		return
	}
	c, err := h.service.ListCandlesPaginated(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error listing candles paginated")
	}
	json.WriteSuccess(w, http.StatusOK, c)

}

//func (h *handler) GetCandlesBySymbol(w http.ResponseWriter, r *http.Request)
//func (h *handler) CountCandles(w http.ResponseWriter, r *http.Request)
//func (h *handler) CountCandlesInRange(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetCandleStats(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetCandleStatsInRange(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetVolumeLeaders(w http.ResponseWriter, r *http.Request)
//func (h *handler) UpdateCandle(w http.ResponseWriter, r *http.Request)
//func (h *handler) DeleteCandle(w http.ResponseWriter, r *http.Request)
//func (h *handler) DeleteCandlesBySymbol(w http.ResponseWriter, r *http.Request)
//func (h *handler) DeleteCandlesByTimeFrame(w http.ResponseWriter, r *http.Request)
//func (h *handler) DeleteOldCandles(w http.ResponseWriter, r *http.Request)
//func (h *handler) DeleteCandlesInRange(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetDistinctSymbols(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetDistinctTimeFrames(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetSymbolTimeframePairs(w http.ResponseWriter, r *http.Request)
//func (h *handler) CheckCandleExists(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetOldestCandle(w http.ResponseWriter, r *http.Request)
//func (h *handler) GetCandleGaps(w http.ResponseWriter, r *http.Request)
