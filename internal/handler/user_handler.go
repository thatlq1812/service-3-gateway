package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/thatlq1812/service-3-gateway/internal/response"

	userpb "github.com/thatlq1812/service-1-user/proto"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// CreateUserRequest HTTP request body
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/v1/users
// Nhận HTTP JSON từ client → gọi gRPC User Service → trả về format mentor
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Call gRPC User Service
	resp, err := h.userClient.CreateUser(context.Background(), &userpb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	// Format response theo mentor yêu cầu
	response.Success(w, map[string]interface{}{
		"id":         resp.Data.User.Id,
		"name":       resp.Data.User.Name,
		"email":      resp.Data.User.Email,
		"created_at": resp.Data.User.CreatedAt,
		"updated_at": resp.Data.User.UpdatedAt,
	})
}

// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	resp, err := h.userClient.GetUser(context.Background(), &userpb.GetUserRequest{
		Id: int32(id),
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"id":         resp.Data.User.Id,
		"name":       resp.Data.User.Name,
		"email":      resp.Data.User.Email,
		"created_at": resp.Data.User.CreatedAt,
		"updated_at": resp.Data.User.UpdatedAt,
	})
}

// UpdateUserRequest HTTP request body
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Build gRPC request with optional fields
	grpcReq := &userpb.UpdateUserRequest{
		Id: int32(id),
	}
	if req.Name != "" {
		grpcReq.Name = &req.Name
	}
	if req.Email != "" {
		grpcReq.Email = &req.Email
	}
	if req.Password != "" {
		grpcReq.Password = &req.Password
	}

	resp, err := h.userClient.UpdateUser(context.Background(), grpcReq)

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"id":         resp.Data.User.Id,
		"name":       resp.Data.User.Name,
		"email":      resp.Data.User.Email,
		"created_at": resp.Data.User.CreatedAt,
		"updated_at": resp.Data.User.UpdatedAt,
	})
}

// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	resp, err := h.userClient.DeleteUser(context.Background(), &userpb.DeleteUserRequest{
		Id: int32(id),
	})

	if err != nil {
		response.Error(w, err)
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

// GET /api/v1/users?page=1&page_size=10
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	resp, err := h.userClient.ListUsers(context.Background(), &userpb.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	users := make([]map[string]interface{}, 0, len(resp.Data.Users))
	for _, u := range resp.Data.Users {
		users = append(users, map[string]interface{}{
			"id":         u.Id,
			"name":       u.Name,
			"email":      u.Email,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		})
	}

	// Format list response theo mentor: {"code":"0", "message":"success", "data":{"items":[...], "total":...}}
	response.SuccessList(w, users, resp.Data.Total, resp.Data.Page, resp.Data.Size)
}

// LoginRequest HTTP request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/v1/auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	resp, err := h.userClient.Login(context.Background(), &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"access_token":  resp.Data.AccessToken,
		"refresh_token": resp.Data.RefreshToken,
	})
}

// ValidateTokenRequest HTTP request body
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// POST /api/v1/auth/validate
func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var req ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	resp, err := h.userClient.ValidateToken(context.Background(), &userpb.ValidateTokenRequest{
		Token: req.Token,
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"valid":   resp.Data.Valid,
		"user_id": resp.Data.UserId,
		"email":   resp.Data.Email,
	})
}

// RefreshTokenRequest HTTP request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// POST /api/v1/auth/refresh
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	resp, err := h.userClient.RefreshToken(context.Background(), &userpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	if resp.Code != "000" {
		response.CustomError(w, resp.Code, resp.Message)
		return
	}

	response.Success(w, map[string]interface{}{
		"access_token":  resp.Data.AccessToken,
		"refresh_token": resp.Data.RefreshToken,
	})
}

// LogoutRequest HTTP request body
type LogoutRequest struct {
	Token        string `json:"token"`         // Access token (required)
	RefreshToken string `json:"refresh_token"` // Refresh token (optional but recommended)
}

// POST /api/v1/auth/logout
// Blacklists both access and refresh tokens for complete logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	resp, err := h.userClient.Logout(context.Background(), &userpb.LogoutRequest{
		Token:        req.Token,
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		response.Error(w, err)
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
