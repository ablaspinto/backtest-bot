package sessions

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

func (h *handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[createSessionParams](w, r)
	if !ok {
		return
	}
	session, err := h.service.CreateSession(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error creating a session")
	}
	json.WriteSuccess(w, http.StatusCreated, session)
}

func (h *handler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "error getting user id")
	}
	session, err := h.service.GetSessionByID(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting session Id")
	}
	json.WriteSuccess(w, http.StatusOK, session)
}

func (h *handler) GetSessionByName(w http.ResponseWriter, r *http.Request) {
	sessionName := reponsehandlers.ParseSymbolName(r, w, "session_name")
	session, err := h.service.GetSessionByName(r.Context(), sessionName)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting session by Name")
	}
	json.WriteSuccess(w, http.StatusOK, session)
}

func (h *handler) ListAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.ListAllSessions(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error listing all session")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) ListSessionWithLimit(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[listSessionParams](w, r)
	if !ok {
		return
	}
	sessions, err := h.service.ListSessionsWithLimit(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error listing sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	activeSessions, err := h.service.GetActiveSessions(r.Context())
	if err != nil {
		json.WriteSuccess(w, http.StatusBadRequest, "Error getting active sessions")
	}
	json.WriteSuccess(w, http.StatusOK, activeSessions)
}

func (h *handler) GetCompletedSessions(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	sessions, err := h.service.GetCompletedSessions(r.Context(), int32(id))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting completed sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetFavoriteSessions(w http.ResponseWriter, r *http.Request) {

	sessions, err := h.service.GetFavoriteSessions(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting favorite sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetSessionsBySymbols(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getSessionBySymbolsParams](w, r)
	if !ok {
		return
	}
	sessions, err := h.service.GetSessionsBySymbol(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting sessions by symbols")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetSessionByStrategy(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getSessionByStrategyParams](w, r)
	if !ok {
		return
	}
	strategySessions, err := h.service.GetSessionByStrategy(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting sessions by strategy")
	}
	json.WriteSuccess(w, http.StatusOK, strategySessions)

}

func (h *handler) GetSessionByStatus(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getSessionByStatusParams](w, r)
	if !ok {
		return
	}
	statusSessions, err := h.service.GetSessionByStatus(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting sessions by status")
	}
	json.WriteSuccess(w, http.StatusOK, statusSessions)

}

func (h *handler) GetRecentSessions(w http.ResponseWriter, r *http.Request) {
	str := reponsehandlers.ParseSymbolName(r, w, "recent")
	sessions, err := h.service.GetRecentSessions(r.Context(), str)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)

}
func (h *handler) SearchSessionsByName(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[searchSessionByNameParams](w, r)
	if !ok {
		return
	}
	sessions, err := h.service.SearchSessionByName(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting sessions by Name")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetSessionWithTags(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getSessionWithTagsParams](w, r)
	if !ok {
		return
	}
	sessions, err := h.service.GetSessionWithTags(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting session by tags")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)

}

func (h *handler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateSessionParams](w, r)
	if !ok {
		return
	}
	updatedSession, err := h.service.UpdateSession(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating session")
	}
	json.WriteSuccess(w, http.StatusOK, updatedSession)
}

func (h *handler) UpdateSessionEnd(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateSessionEndParams](w, r)
	if !ok {
		return
	}
	updatedSession, err := h.service.UpdateSessionEnd(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating session")
	}
	json.WriteSuccess(w, http.StatusOK, updatedSession)

}

func (h *handler) UpdateSessionStatus(w http.ResponseWriter, r *http.Request) {

	params, ok := json.ReadOrError[updateSessionStatus](w, r)
	if !ok {
		return
	}
	err := h.service.UpdateSessionStatus(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating session status")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UpdateSessionPrices(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateSessionPricesParams](w, r)
	if !ok {
		return
	}
	err := h.service.UpdateSessionPrices(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating session prices")
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) ToggleSessionFavorite(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "favorite")
	if !ok {
		return
	}
	boolean, err := h.service.ToggleSessionFavorite(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error toggling session favorite")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)
}

func (h *handler) AddSessionTag(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[addSessionTagParams](w, r)
	if !ok {
		return
	}
	err := h.service.AddSessionTag(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error adding session tag")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) RemoveSessionTag(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[removeSessionTagParams](w, r)
	if !ok {
		return
	}
	err := h.service.RemoveSessionTag(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error removing session tag")
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	err := h.service.DeleteSession(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting session")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteOldSessions(w http.ResponseWriter, r *http.Request) {
	oldStr := reponsehandlers.ParseSymbolName(r, w, "old")
	err := h.service.DeleteOldSessions(r.Context(), oldStr)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting old sessions")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteSessionsByStatus(w http.ResponseWriter, r *http.Request) {
	status := reponsehandlers.ParseSymbolName(r, w, "status")
	err := h.service.DeleteSessionsByStatus(r.Context(), status)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error deleting old session by status")
	}
}

func (h *handler) CountSessions(w http.ResponseWriter, r *http.Request) {
	num, err := h.service.CountSessions(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error counting sessions")
	}
	json.WriteSuccess(w, http.StatusOK, num)
}

func (h *handler) CountSessionsByStatus(w http.ResponseWriter, r *http.Request) {
	status := reponsehandlers.ParseSymbolName(r, w, "status")
	num, err := h.service.CountSessionsByStatus(r.Context(), status)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error counting session by status")
	}
	json.WriteSuccess(w, http.StatusOK, num)
}

func (h *handler) GetSessionStats(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.GetSessionStats(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting session stats")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)

}

func (h *handler) GetSymbolStats(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.GetSymbolStats(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting symbol stats")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetStrategyPerformance(w http.ResponseWriter, r *http.Request) {

	performance, err := h.service.GetStrategyPerformance(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting strategy performance")
	}
	json.WriteSuccess(w, http.StatusOK, performance)

}

func (h *handler) GetDailySessionCount(w http.ResponseWriter, r *http.Request) {

	timestamp := reponsehandlers.ParseSymbolName(r, w, "ts")
	count, err := h.service.GetDailySessionCount(r.Context(), timestamp)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting daily session count")
	}
	json.WriteSuccess(w, http.StatusOK, count)

}

func (h *handler) CheckSessionExists(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}

	exists, err := h.service.CheckSessionExists(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error checking if sessions exists")
	}
	json.WriteSuccess(w, http.StatusOK, exists)
}

func (h *handler) GetLatestSession(w http.ResponseWriter, r *http.Request) {
	session, err := h.service.GetLatestSession(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting latest sesssion")
	}
	json.WriteSuccess(w, http.StatusOK, session)
}

func (h *handler) GetLongestSession(w http.ResponseWriter, r *http.Request) {
	num, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	sessions, err := h.service.GetLongestSession(r.Context(), int32(num))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error grabbing longest sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)

}
func (h *handler) GetMostProfitableSession(w http.ResponseWriter, r *http.Request) {

	num, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	sessions, err := h.service.GetMostProfitableSession(r.Context(), int32(num))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting most profitable sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetMostVolatileSession(w http.ResponseWriter, r *http.Request) {
	num, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	sessions, err := h.service.GetMostVolatileSession(r.Context(), int32(num))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting most volatile sessions")
	}
	json.WriteSuccess(w, http.StatusOK, sessions)
}

func (h *handler) GetSessionDuration(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "session_id")
	if !ok {
		return
	}
	duration, err := h.service.GetSessionDuration(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting session duration")
	}
	json.WriteSuccess(w, http.StatusOK, duration)
}
