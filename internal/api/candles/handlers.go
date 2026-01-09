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

func (h *handler) GetCandlesBySymbol(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandlesBySymbolParams](w, r)
	if !ok {
		return
	}
	candles, err := h.service.GetCandlesBySymbol(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting candles by symbol")
	}
	json.WriteSuccess(w, http.StatusOK, candles)
}

func (h *handler) CountCandles(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[countCandlesParams](w, r)
	if !ok {
		return
	}
	num, err := h.service.CountCandles(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting the number of candles")
	}
	json.WriteSuccess(w, http.StatusOK, num)

}

func (h *handler) CountCandlesInRange(w http.ResponseWriter, r *http.Request) {

	params, ok := json.ReadOrError[countCandlesInRangeParams](w, r)
	if !ok {
		return
	}
	num, err := h.service.CountCandlesInRange(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting candles within this range")
	}
	json.WriteSuccess(w, http.StatusOK, num)

}

func (h *handler) GetCandleStats(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandleStatsParams](w, r)
	if !ok {
		return
	}
	stats, err := h.service.GetCandleStats(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting candle stats")
	}
	json.WriteSuccess(w, http.StatusOK, stats)
}

func (h *handler) GetCandleStatsInRange(w http.ResponseWriter, r *http.Request) {

	params, ok := json.ReadOrError[getCandleStatsInRangeParams](w, r)
	if !ok {
		return
	}
	stats, err := h.service.GetCandleStatsInRange(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting candle stats within range")
	}
	json.WriteSuccess(w, http.StatusOK, stats)

}

func (h *handler) GetVolumeLeaders(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getVolumeLeadersParams](w, r)
	if !ok {
		return
	}
	volumeLeaders, err := h.service.GetVolumeLeaders(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting volume leaders")
	}
	json.WriteSuccess(w, http.StatusOK, volumeLeaders)

}

func (h *handler) UpdateCandle(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateCandleParams](w, r)
	if !ok {
		return
	}
	updatedCandle, err := h.service.UpdateCandle(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating candle")
	}
	json.WriteSuccess(w, http.StatusOK, updatedCandle)

}

func (h *handler) DeleteCandle(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	err := h.service.DeleteCandle(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error deleting candle")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteCandlesBySymbol(w http.ResponseWriter, r *http.Request) {
	sym := reponsehandlers.ParseSymbolName(r, w, "sym")
	err := h.service.DeleteCandlesBySymbol(r.Context(), sym)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting candles by Symbol")
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) DeleteCandlesByTimeFrame(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[deleteCandlesByTimeFrameParams](w, r)
	if !ok {
		return
	}
	err := h.service.DeleteCandlesByTimeFrame(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting candles by timeframe")
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *handler) DeleteOldCandles(w http.ResponseWriter, r *http.Request) {
	timestamp, ok := reponsehandlers.ParseIdResponse(r, w, "timestamp")
	if !ok {
		return
	}
	err := h.service.DeleteOldCandles(r.Context(), timestamp)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting old candles")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteCandlesInRange(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[deleteCandlesInRangeParams](w, r)
	if !ok {
		return
	}
	err := h.service.DeleteCandlesInRange(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting candles in Range")
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) GetDistinctSymbols(w http.ResponseWriter, r *http.Request) {
	sym := reponsehandlers.ParseSymbolName(r, w, "sym")
	symbols, err := h.service.GetDistinctSymbols(r.Context(), sym)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting distinct symbols")
	}
	json.WriteSuccess(w, http.StatusOK, symbols)
}

func (h *handler) GetDistinctTimeFrames(w http.ResponseWriter, r *http.Request) {

	sym := reponsehandlers.ParseSymbolName(r, w, "sym")
	timeFrames, err := h.service.GetDistinctTimeFrames(r.Context(), sym)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting distinct time frames from symbol")
	}
	json.WriteSuccess(w, http.StatusOK, timeFrames)

}

func (h *handler) GetSymbolTimeframePairs(w http.ResponseWriter, r *http.Request) {
	pairs, err := h.service.GetSymbolTimeframePairs(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting symbol and timeframe pairs")
	}
	json.WriteSuccess(w, http.StatusOK, pairs)
}

func (h *handler) CheckCandleExists(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[checkCandleExistsParams](w, r)
	if !ok {
		return
	}
	boolean, err := h.service.CheckCandleExists(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error seeing if candle exitsts")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)

}

func (h *handler) GetOldestCandle(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getOldsCandleParams](w, r)
	if !ok {
		return
	}
	oldestCandle, err := h.service.GetOldestCandle(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting oldest candle")
	}
	json.WriteSuccess(w, http.StatusOK, oldestCandle)
}

func (h *handler) GetCandleGaps(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getCandleGapsParams](w, r)
	if !ok {
		return
	}
	gaps, err := h.service.GetCandleGaps(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting gaps between candles")
	}
	json.WriteSuccess(w, http.StatusOK, gaps)
}
