package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"service-3-gateway/internal/response"

	articlepb "github.com/thatlq1812/service-2-article/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ArticleHandler struct {
	articleClient articlepb.ArticleServiceClient
}

func NewArticleHandler(articleClient articlepb.ArticleServiceClient) *ArticleHandler {
	return &ArticleHandler{
		articleClient: articleClient,
	}
}

// extractToken gets JWT token from Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	// Bearer <token>
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// CreateArticleRequest HTTP request body
type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int32  `json:"user_id"`
}

// POST /api/v1/articles
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Extract JWT token from Authorization header
	token := extractToken(r)
	if token == "" {
		response.Unauthorized(w, "authorization token required")
		return
	}

	// Forward token in gRPC metadata
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)

	// Call gRPC Article Service
	resp, err := h.articleClient.CreateArticle(ctx, &articlepb.CreateArticleRequest{
		Title:   req.Title,
		Content: req.Content,
		UserId:  req.UserID,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			response.GRPCError(w, st.Code(), st.Message())
		} else {
			response.Error(w, err)
		}
		return
	}

	// gRPC response is already wrapped with code, message, data
	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"id":         resp.Data.Article.Id,
		"title":      resp.Data.Article.Title,
		"content":    resp.Data.Article.Content,
		"user_id":    resp.Data.Article.UserId,
		"created_at": resp.Data.Article.CreatedAt,
		"updated_at": resp.Data.Article.UpdatedAt,
	})
}

// GET /api/v1/articles/{id}
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid article id")
		return
	}

	resp, err := h.articleClient.GetArticle(context.Background(), &articlepb.GetArticleRequest{
		Id: int32(id),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			response.GRPCError(w, st.Code(), st.Message())
		} else {
			response.Error(w, err)
		}
		return
	}

	// Response already wrapped by gRPC
	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	// Format ArticleWithUser response
	articleData := map[string]interface{}{
		"id":         resp.Data.Article.Article.Id,
		"title":      resp.Data.Article.Article.Title,
		"content":    resp.Data.Article.Article.Content,
		"user_id":    resp.Data.Article.Article.UserId,
		"created_at": resp.Data.Article.Article.CreatedAt,
		"updated_at": resp.Data.Article.Article.UpdatedAt,
	}

	// Include user info if available (null if User Service unavailable)
	if resp.Data.Article.User != nil {
		articleData["user"] = map[string]interface{}{
			"id":         resp.Data.Article.User.Id,
			"name":       resp.Data.Article.User.Name,
			"email":      resp.Data.Article.User.Email,
			"created_at": resp.Data.Article.User.CreatedAt,
			"updated_at": resp.Data.Article.User.UpdatedAt,
		}
	} else {
		articleData["user"] = nil // Explicitly set null for graceful degradation
	}

	response.Success(w, articleData)
}

// UpdateArticleRequest HTTP request body
type UpdateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PUT /api/v1/articles/{id}
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid article id")
		return
	}

	var req UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	resp, err := h.articleClient.UpdateArticle(context.Background(), &articlepb.UpdateArticleRequest{
		Id:      int32(id),
		Title:   req.Title,
		Content: req.Content,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			response.GRPCError(w, st.Code(), st.Message())
		} else {
			response.Error(w, err)
		}
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"id":         resp.Data.Article.Id,
		"title":      resp.Data.Article.Title,
		"content":    resp.Data.Article.Content,
		"user_id":    resp.Data.Article.UserId,
		"created_at": resp.Data.Article.CreatedAt,
		"updated_at": resp.Data.Article.UpdatedAt,
	})
}

// DELETE /api/v1/articles/{id}
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid article id")
		return
	}

	resp, err := h.articleClient.DeleteArticle(context.Background(), &articlepb.DeleteArticleRequest{
		Id: int32(id),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			response.GRPCError(w, st.Code(), st.Message())
		} else {
			response.Error(w, err)
		}
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"success": resp.Data.Success,
	})
}

// GET /api/v1/articles?page=1&page_size=10&user_id=1
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if pageNumber < 1 {
		pageNumber = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	resp, err := h.articleClient.ListArticles(context.Background(), &articlepb.ListArticlesRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
		UserId:     int32(userID),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			response.GRPCError(w, st.Code(), st.Message())
		} else {
			response.Error(w, err)
		}
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	articles := make([]map[string]interface{}, 0, len(resp.Data.Articles))
	for _, aw := range resp.Data.Articles {
		articleData := map[string]interface{}{
			"id":         aw.Article.Id,
			"title":      aw.Article.Title,
			"content":    aw.Article.Content,
			"user_id":    aw.Article.UserId,
			"created_at": aw.Article.CreatedAt,
			"updated_at": aw.Article.UpdatedAt,
		}

		// Include user info if available
		if aw.User != nil {
			articleData["user"] = map[string]interface{}{
				"id":         aw.User.Id,
				"name":       aw.User.Name,
				"email":      aw.User.Email,
				"created_at": aw.User.CreatedAt,
				"updated_at": aw.User.UpdatedAt,
			}
		}

		articles = append(articles, articleData)
	}

	// Format list response theo mentor
	response.SuccessList(w, articles, int64(resp.Data.Total), resp.Data.Page, int32(pageSize))
}
