package response

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIResponse standard format
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListData for paginated responses
type ListData struct {
	Items   interface{} `json:"items"`
	Total   int64       `json:"total"`
	Page    int32       `json:"page"`
	Size    int32       `json:"size"`
	HasMore bool        `json:"has_more"`
}

// gRPC status codes mapping to string format
// Theo chuẩn gRPC: https://grpc.io/docs/guides/status-codes/
const (
	CodeOK                 = "000" // OK
	CodeCancelled          = "001" // CANCELLED
	CodeUnknown            = "002" // UNKNOWN
	CodeInvalidArgument    = "003" // INVALID_ARGUMENT
	CodeDeadlineExceeded   = "004" // DEADLINE_EXCEEDED
	CodeNotFound           = "005" // NOT_FOUND
	CodeAlreadyExists      = "006" // ALREADY_EXISTS
	CodePermissionDenied   = "007" // PERMISSION_DENIED
	CodeResourceExhausted  = "008" // RESOURCE_EXHAUSTED
	CodeFailedPrecondition = "009" // FAILED_PRECONDITION
	CodeAborted            = "010" // ABORTED
	CodeOutOfRange         = "011" // OUT_OF_RANGE
	CodeUnimplemented      = "012" // UNIMPLEMENTED
	CodeInternal           = "013" // INTERNAL
	CodeUnavailable        = "014" // UNAVAILABLE
	CodeDataLoss           = "015" // DATA_LOSS
	CodeUnauthenticated    = "016" // UNAUTHENTICATED
)

// MapGRPCCodeToString converts gRPC status code to string format
func MapGRPCCodeToString(code codes.Code) string {
	switch code {
	case codes.OK:
		return CodeOK
	case codes.Canceled:
		return CodeCancelled
	case codes.Unknown:
		return CodeUnknown
	case codes.InvalidArgument:
		return CodeInvalidArgument
	case codes.DeadlineExceeded:
		return CodeDeadlineExceeded
	case codes.NotFound:
		return CodeNotFound
	case codes.AlreadyExists:
		return CodeAlreadyExists
	case codes.PermissionDenied:
		return CodePermissionDenied
	case codes.ResourceExhausted:
		return CodeResourceExhausted
	case codes.FailedPrecondition:
		return CodeFailedPrecondition
	case codes.Aborted:
		return CodeAborted
	case codes.OutOfRange:
		return CodeOutOfRange
	case codes.Unimplemented:
		return CodeUnimplemented
	case codes.Internal:
		return CodeInternal
	case codes.Unavailable:
		return CodeUnavailable
	case codes.DataLoss:
		return CodeDataLoss
	case codes.Unauthenticated:
		return CodeUnauthenticated
	default:
		return CodeUnknown
	}
}

// MapGRPCCodeToHTTPStatus converts gRPC code to HTTP status code
func MapGRPCCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499 // Client Closed Request
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Success returns success response với code "0"
func Success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeOK,
		Message: "success",
		Data:    data,
	})
}

// SuccessList returns paginated success response
func SuccessList(w http.ResponseWriter, items interface{}, total int64, page, size int32) {
	hasMore := int64(page*size) < total
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeOK,
		Message: "success",
		Data: ListData{
			Items:   items,
			Total:   total,
			Page:    page,
			Size:    size,
			HasMore: hasMore,
		},
	})
}

// Error converts gRPC error to API error response
// Đây là hàm chính để convert từ gRPC status sang format mentor yêu cầu
func Error(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error - treat as internal error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Code:    CodeInternal,
			Message: err.Error(),
		})
		return
	}

	// Convert gRPC code to API code (string format)
	apiCode := MapGRPCCodeToString(st.Code())
	httpStatus := MapGRPCCodeToHTTPStatus(st.Code())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    apiCode,
		Message: st.Message(),
	})
}

// BadRequest returns invalid argument error (code "3")
func BadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeInvalidArgument,
		Message: message,
	})
}

// Unauthorized returns unauthenticated error (code "16")
func Unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeUnauthenticated,
		Message: message,
	})
}

// NotFound returns not found error (code "5")
func NotFound(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeNotFound,
		Message: message,
	})
}

// ServiceUnavailable returns service unavailable error (code "14")
// Used when circuit breaker is open or service is down
func ServiceUnavailable(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeUnavailable,
		Message: message,
	})
}

// InternalError returns internal error (code "13")
func InternalError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodeInternal,
		Message: message,
	})
}

// Forbidden returns permission denied error (code "7")
func Forbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    CodePermissionDenied,
		Message: message,
	})
}

// GRPCError converts gRPC code and message to API error response
func GRPCError(w http.ResponseWriter, code codes.Code, message string) {
	apiCode := MapGRPCCodeToString(code)
	httpStatus := MapGRPCCodeToHTTPStatus(code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    apiCode,
		Message: message,
	})
}

// CustomError returns custom error with specific code
func CustomError(w http.ResponseWriter, code string, message string) {
	// Map string code to HTTP status
	var httpStatus int
	switch code {
	case CodeOK:
		httpStatus = http.StatusOK
	case CodeInvalidArgument:
		httpStatus = http.StatusBadRequest
	case CodeNotFound:
		httpStatus = http.StatusNotFound
	case CodeAlreadyExists:
		httpStatus = http.StatusConflict
	case CodePermissionDenied:
		httpStatus = http.StatusForbidden
	case CodeUnauthenticated:
		httpStatus = http.StatusUnauthorized
	case CodeUnavailable:
		httpStatus = http.StatusServiceUnavailable
	default:
		httpStatus = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    code,
		Message: message,
	})
}
