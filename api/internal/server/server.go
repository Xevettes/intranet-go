package server

import (
	"api/internal/config"
	"api/internal/gateways"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"api/proto/vault"
	monitoring "api/proto/zabbix"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"google.golang.org/protobuf/types/known/structpb"
)

const CacheTTL = 60 * time.Second

type Server struct {
	router         *chi.Mux
	gatewayManager *gateways.Manager
}

func NewServer(manager *gateways.Manager, cfg *config.Config) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		gatewayManager: manager,
	}

	s.router.Use(chi_middleware.RequestID)
	s.router.Use(chi_middleware.RealIP)
	s.router.Use(chi_middleware.Logger)
	s.router.Use(chi_middleware.Recoverer)
	s.router.Get("/health", s.healthCheck)

	if s.gatewayManager.VaultClient != nil {
		s.router.Route("/api/v1/secrets", func(r chi.Router) {
			r.Get("/*", s.handleReadOrListSecret)
			r.Post("/*", s.handleWriteSecret)
			r.Put("/*", s.handleWriteSecret)
			r.Patch("/*", s.handlePatchSecret)
			r.Delete("/*", s.handleSoftDeleteSecret)
		})
		slog.Info("Vault routes registered")
	}
	if s.gatewayManager.ZabbixClient != nil {
		s.router.Route("/api/v1/zabbix", func(r chi.Router) {
			r.Get("/hostgroups", s.handleListHostGroups)
			r.Get("/hosts", s.handleListHosts)
			r.Get("/items", s.handleListItems)
			r.Get("/alerts", s.handleListAlerts)
		})
		slog.Info("Zabbix routes registered")
	}

	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
func (s *Server) handleReadOrListSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vaultPath := strings.TrimPrefix(r.URL.Path, "/api/v1/secrets/")

	if strings.Contains(vaultPath, "/metadata/") {
		grpcRequest := &vault.ListSecretsRequest{Path: vaultPath}

		response, err := s.gatewayManager.VaultClient.ListSecrets(ctx, grpcRequest)
		if err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Erro ao listar segredos do Vault", err)
			return
		}
		s.respondWithJSON(w, http.StatusOK, response.GetKeys())
		return
	}

	grpcRequest := &vault.ReadSecretRequest{Path: vaultPath}

	secret, err := s.gatewayManager.VaultClient.ReadSecret(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao ler segredo do Vault", err)
		return
	}
	s.respondWithJSON(w, http.StatusOK, secret.GetData().AsMap())
}
func (s *Server) handleWriteSecret(w http.ResponseWriter, r *http.Request) {
	vaultPath := strings.TrimPrefix(r.URL.Path, "/api/v1/secrets/")
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Corpo da requisição JSON inválido", err)
		return
	}

	grpcPayload, err := structpb.NewStruct(map[string]interface{}{"data": payload})
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro interno ao converter payload", err)
		return
	}
	grpcRequest := &vault.WriteSecretRequest{Path: vaultPath, Data: grpcPayload}

	_, err = s.gatewayManager.VaultClient.WriteSecret(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao escrever segredo no Vault", err)
		return
	}
	s.respondWithJSON(w, http.StatusCreated, map[string]string{"status": "success"})
}
func (s *Server) handlePatchSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vaultPath := strings.TrimPrefix(r.URL.Path, "/api/v1/secrets/")

	readReq := &vault.ReadSecretRequest{Path: vaultPath}
	existingSecret, err := s.gatewayManager.VaultClient.ReadSecret(ctx, readReq)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao ler segredo para atualização", err)
		return
	}

	existingData := existingSecret.GetData().AsMap()["data"].(map[string]interface{})

	var patchData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Corpo da requisição JSON inválido", err)
		return
	}

	for key, value := range patchData {
		existingData[key] = value
	}

	grpcPayload, _ := structpb.NewStruct(map[string]interface{}{"data": existingData})

	writeReq := &vault.WriteSecretRequest{Path: vaultPath, Data: grpcPayload}
	_, err = s.gatewayManager.VaultClient.WriteSecret(ctx, writeReq)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao atualizar segredo no Vault", err)
		return
	}

	s.respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
func (s *Server) handleSoftDeleteSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vaultPath := strings.TrimPrefix(r.URL.Path, "/api/v1/secrets/")

	readReq := &vault.ReadSecretRequest{Path: vaultPath}
	existingSecret, err := s.gatewayManager.VaultClient.ReadSecret(ctx, readReq)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao ler segredo para deletar", err)
		return
	}

	existingData := existingSecret.GetData().AsMap()["data"].(map[string]interface{})

	existingData["hidden"] = "true"

	grpcPayload, _ := structpb.NewStruct(map[string]interface{}{"data": existingData})

	writeReq := &vault.WriteSecretRequest{Path: vaultPath, Data: grpcPayload}
	_, err = s.gatewayManager.VaultClient.WriteSecret(ctx, writeReq)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Erro ao marcar segredo como oculto", err)
		return
	}

	s.respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (s *Server) handleListHostGroups(w http.ResponseWriter, r *http.Request) {
	grpcRequest := &monitoring.ListHostGroupsRequest{}
	response, err := s.gatewayManager.ZabbixClient.ListHostGroups(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Falha ao buscar grupos de hosts do Zabbix", err)
		return
	}
	s.respondWithJSON(w, http.StatusOK, response.GetGroups())
}
func (s *Server) handleListHosts(w http.ResponseWriter, r *http.Request) {
	groupIds := r.URL.Query()["groupids"]
	if len(groupIds) == 0 {
		s.respondWithError(w, http.StatusBadRequest, "Parâmetro 'groupids' é obrigatório", nil)
		return
	}
	grpcRequest := &monitoring.ListHostsRequest{Groupids: groupIds}
	response, err := s.gatewayManager.ZabbixClient.ListHosts(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Falha ao buscar hosts do Zabbix", err)
		return
	}
	s.respondWithJSON(w, http.StatusOK, response.GetHosts())
}
func (s *Server) handleListItems(w http.ResponseWriter, r *http.Request) {
	hostIds := r.URL.Query()["hostids"]
	if len(hostIds) == 0 {
		s.respondWithError(w, http.StatusBadRequest, "Parâmetro 'hostids' é obrigatório", nil)
		return
	}
	grpcRequest := &monitoring.ListItemsRequest{Hostids: hostIds}
	response, err := s.gatewayManager.ZabbixClient.ListItems(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Falha ao buscar itens do Zabbix", err)
		return
	}
	s.respondWithJSON(w, http.StatusOK, response.GetItems())
}
func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	hostIds := r.URL.Query()["hostids"]
	if len(hostIds) == 0 {
		s.respondWithError(w, http.StatusBadRequest, "Parâmetro 'hostids' é obrigatório", nil)
		return
	}
	grpcRequest := &monitoring.ListAlertsRequest{Hostids: hostIds}
	response, err := s.gatewayManager.ZabbixClient.ListAlerts(r.Context(), grpcRequest)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Falha ao buscar alertas do Zabbix", err)
		return
	}
	s.respondWithJSON(w, http.StatusOK, response.GetAlerts())
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Falha ao serializar resposta JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func (s *Server) respondWithError(w http.ResponseWriter, code int, message string, internalErr error) {
	if internalErr != nil {
		slog.Error(message, "error", internalErr)
	}
	s.respondWithJSON(w, code, map[string]string{"error": message})
}
