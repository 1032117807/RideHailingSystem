package handler

import (
	"io"
	"net/http"
	"strings"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type KnowledgeHandler struct {
	knowledgeService *service.KnowledgeService
}

func NewKnowledgeHandler(knowledgeService *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{knowledgeService: knowledgeService}
}

func (h *KnowledgeHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can upload knowledge")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	category := strings.TrimSpace(r.FormValue("category"))

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	fileName := strings.TrimSpace(fileHeader.Filename)
	lowerName := strings.ToLower(fileName)
	if !strings.HasSuffix(lowerName, ".md") && !strings.HasSuffix(lowerName, ".txt") {
		response.Error(w, http.StatusBadRequest, "only .md and .txt are supported")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, 10<<20))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "read file failed")
		return
	}

	doc, err := h.knowledgeService.UploadDocument(r.Context(), service.UploadKnowledgeInput{
		Title:      title,
		Category:   category,
		SourceName: fileName,
		Content:    string(data),
		CreatedBy:  currentUser.ID,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, doc)
}

func (h *KnowledgeHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can list knowledge")
		return
	}

	category := strings.TrimSpace(r.URL.Query().Get("category"))
	docs, err := h.knowledgeService.ListDocuments(r.Context(), category)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, docs)
}

func (h *KnowledgeHandler) SearchDocuments(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can search knowledge")
		return
	}

	var req struct {
		Query    string `json:"query"`
		TopK     int    `json:"topK"`
		Category string `json:"category"`
	}
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := h.knowledgeService.SearchKnowledge(r.Context(), service.SearchKnowledgeInput{
		Query:    req.Query,
		TopK:     req.TopK,
		Category: req.Category,
		UserID:   currentUser.ID,
		Role:     currentUser.Role,
		Feature:  model.TokenFeatureKnowledgeSearch,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, results)
}

func (h *KnowledgeHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can view knowledge")
		return
	}

	documentID := strings.TrimSpace(r.PathValue("documentId"))
	doc, err := h.knowledgeService.GetDocument(r.Context(), documentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, doc)
}

func (h *KnowledgeHandler) UpdateDocumentStatus(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can update knowledge")
		return
	}

	documentID := strings.TrimSpace(r.PathValue("documentId"))
	var req struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Content  string `json:"content"`
		Status   string `json:"status"`
	}
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	doc, err := h.knowledgeService.UpdateDocument(r.Context(), documentID, service.UpdateKnowledgeDocumentInput{
		Title:    req.Title,
		Category: req.Category,
		Content:  req.Content,
		Status:   req.Status,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, doc)
}

func (h *KnowledgeHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can delete knowledge")
		return
	}

	documentID := strings.TrimSpace(r.PathValue("documentId"))
	if err := h.knowledgeService.DeleteDocument(r.Context(), documentID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"documentId": documentID,
		"deleted":    true,
	})
}

func (h *KnowledgeHandler) ReindexDocument(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if currentUser.Role != model.RoleAdmin {
		response.Error(w, http.StatusForbidden, "only admin can reindex knowledge")
		return
	}

	documentID := strings.TrimSpace(r.PathValue("documentId"))
	doc, err := h.knowledgeService.ReindexDocument(r.Context(), documentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, doc)
}
