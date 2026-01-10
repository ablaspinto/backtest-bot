package userconfigs

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

func (h *handler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[createConfigParams](w, r)
	if !ok {
		return
	}
	createdConfig, err := h.service.CreateConfig(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error creating config")
	}
	json.WriteSuccess(w, http.StatusCreated, createdConfig)
}
func (h *handler) GetConfigByID(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	conf, err := h.service.GetConfigByID(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting config by id")
	}
	json.WriteSuccess(w, http.StatusOK, conf)

}

func (h *handler) GetConfigByName(w http.ResponseWriter, r *http.Request) {
	name := reponsehandlers.ParseSymbolName(r, w, "name")
	conf, err := h.service.GetConfigByName(r.Context(), name)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting config by name")
	}
	json.WriteSuccess(w, http.StatusOK, conf)
}
func (h *handler) GetDefaultConfig(w http.ResponseWriter, r *http.Request) {
	defConf, err := h.service.GetDefaultConfig(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting default config")
	}
	json.WriteSuccess(w, http.StatusOK, defConf)
}
func (h *handler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	confs, err := h.service.ListConfigs(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error Listing user configs")
	}
	json.WriteSuccess(w, http.StatusOK, confs)

}
func (h *handler) ListConfigsPaginated(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[listConfigsPaginatedParams](w, r)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Error reading in config paginated params")
	}
	confs, err := h.service.ListConfigsPaginated(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error listing configs paginated")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetFavoriteConfigs(w http.ResponseWriter, r *http.Request) {
	confs, err := h.service.GetFavoriteConfigs(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting favorite configs")
	}

	json.WriteSuccess(w, http.StatusOK, confs)
}
func (h *handler) GetPublicConfigs(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseId32Response(r, w, "id")
	if !ok {
		return
	}
	confs, err := h.service.GetPublicConfigs(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting public configs")
	}
	json.WriteSuccess(w, http.StatusOK, confs)

}
func (h *handler) GetConfigsByStrategy(w http.ResponseWriter, r *http.Request) {
	strat := reponsehandlers.ParseSymbolName(r, w, "name")
	confs, err := h.service.GetConfigsByStrategy(r.Context(), strat)
	if err != nil {
		json.WriteError(w, http.StatusOK, "Error getting configs by strategy")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetConfigsByMarketingStrategy(w http.ResponseWriter, r *http.Request) {
	strat := reponsehandlers.ParseSymbolName(r, w, "name")
	confs, err := h.service.GetConfigsByMarketStrategy(r.Context(), strat)
	if err != nil {
		json.WriteError(w, http.StatusOK, "Error getting configs by strategy")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetConfigsByCreator(w http.ResponseWriter, r *http.Request) {
	timeStamp := reponsehandlers.ParseSymbolName(r, w, "ts")
	confs, err := h.service.GetConfigsByCreator(r.Context(), timeStamp)
	if err != nil {
		json.WriteError(w, http.StatusOK, "Error getting configs by creator")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) SearchConfigsByName(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[searchConfigsByNameParams](w, r)
	if !ok {
		return
	}
	confs, err := h.service.SearchConfigsByName(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error searching configs by name")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetConfigsWithTags(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[getConfigWithTagsParams](w, r)
	if !ok {
		return
	}
	confs, err := h.service.GetConfigsWithTags(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting configs with tags")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetRecentlyUsedConfigs(w http.ResponseWriter, r *http.Request) {
	limit, ok := reponsehandlers.ParseId32Response(r, w, "limit")
	if !ok {
		return
	}
	confs, err := h.service.GetRecentlyUsedConfigs(r.Context(), limit)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error getting recently used configs")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) GetPopularConfigs(w http.ResponseWriter, r *http.Request) {
	limit, ok := reponsehandlers.ParseId32Response(r, w, "limit")
	if !ok {
		return
	}
	confs, err := h.service.GetPopularConfigs(r.Context(), limit)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting populat configs")
	}
	json.WriteSuccess(w, http.StatusOK, confs)
}

func (h *handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateConfigParams](w, r)
	if !ok {
		return
	}
	updatedConf, err := h.service.UpdateConfig(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error updating config")
	}
	json.WriteSuccess(w, http.StatusOK, updatedConf)
}
func (h *handler) UpdateConfigPartial(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[updateConfigPartialParams](w, r)
	if !ok {
		return
	}
	updatedConf, err := h.service.UpdateConfigPartial(r.Context(), params)
	if err != nil {
		json.Error(w, http.StatusBadRequest, "Error updating config partial")
	}
	json.WriteSuccess(w, http.StatusOK, updatedConf)
}

func (h *handler) IncrementConfigUsage(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	err := h.service.IncrementConfigUsage(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error incrementing config usage")
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *handler) SetDefaultContent(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	err := h.service.SetDefaultConfig(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error setting default config")
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *handler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	boolean, err := h.service.ToggleFavorite(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error toggling favorite")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)
}

func (h *handler) ToggleConfigPublic(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	boolean, err := h.service.ToggleConfigPublic(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error toggling to public config")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)
}

func (h *handler) AddConfigTag(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[addConfigTagParams](w, r)
	if !ok {
		return
	}
	err := h.service.AddConfigTag(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error adding config tag")
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) RemoveConfigTag(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[removeConfigTagParams](w, r)
	if !ok {
		return
	}
	err := h.service.RemoveConfigTag(r.Context(), params)
	if err != nil {
		json.Error(w, http.StatusBadRequest, "Error removing config tag")
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *handler) CountConfigs(w http.ResponseWriter, r *http.Request) {
	num, err := h.service.CountConfigs(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error counting number of configs")
	}
	json.WriteSuccess(w, http.StatusOK, num)
}

func (h *handler) CountConfigsByStrategy(w http.ResponseWriter, r *http.Request) {
	strat := reponsehandlers.ParseSymbolName(r, w, "name")
	num, err := h.service.CountConfigsByStrategy(r.Context(), strat)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error counting number of configs")
	}
	json.WriteSuccess(w, http.StatusOK, num)
}

func (h *handler) GetConfigStats(w http.ResponseWriter, r *http.Request) {
	rows, err := h.service.GetConfigStats(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting config stats")
	}
	json.WriteSuccess(w, http.StatusOK, rows)
}
func (h *handler) GetStrategyDistribution(w http.ResponseWriter, r *http.Request) {
	strategyDistributionRows, err := h.service.GetStrategyDistribution(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting strategy distribution")
	}
	json.WriteSuccess(w, http.StatusOK, strategyDistributionRows)
}
func (h *handler) GetMarketBehaviorDistribution(w http.ResponseWriter, r *http.Request) {
	marketDistributionRows, err := h.service.GetMarketBehaviorDistribution(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting market behavior distribution rows")
	}

	json.WriteSuccess(w, http.StatusOK, marketDistributionRows)
}

func (h *handler) CheckConfigExists(w http.ResponseWriter, r *http.Request) {
	id, ok := reponsehandlers.ParseIdResponse(r, w, "id")
	if !ok {
		return
	}
	boolean, err := h.service.CheckConfigExists(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "error checking if config exists")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)

}
func (h *handler) CheckConfigNameExists(w http.ResponseWriter, r *http.Request) {
	name := reponsehandlers.ParseSymbolName(r, w, "name")
	boolean, err := h.service.CheckConfigNameExists(r.Context(), name)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error checking if configuration name exists")
	}
	json.WriteSuccess(w, http.StatusOK, boolean)
}
func (h *handler) GetConfigNameOrDefault(w http.ResponseWriter, r *http.Request) {
	name := reponsehandlers.ParseSymbolName(r, w, "name")
	conf, err := h.service.GetConfigNameOrDefault(r.Context(), name)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error getting config by name or default")
	}
	json.WriteSuccess(w, http.StatusOK, conf)
}

func (h *handler) DuplicateConfig(w http.ResponseWriter, r *http.Request) {
	params, ok := json.ReadOrError[duplicateConfigParams](w, r)
	if !ok {
		return
	}
	conf, err := h.service.DuplicateConfig(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Error duplicating config")
	}
	json.WriteSuccess(w, http.StatusOK, conf)
}
