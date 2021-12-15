package handle

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	BadRequestErrorResponse           = HTTPErrorResponse{HTTPStatusCode: 400, ErrorCode: "400", Message: "Requisição mal formatada"}
	UnauthorizedErrorResponse         = HTTPErrorResponse{HTTPStatusCode: 401, ErrorCode: "401", Message: "Necessita autenticação"}
	UnauthorizedLoginResponse         = HTTPErrorResponse{HTTPStatusCode: 401, ErrorCode: "401", Message: "invalid username or password"}
	ExpiredTokenErrorResponse         = HTTPErrorResponse{HTTPStatusCode: 401, ErrorCode: "401", Message: "Token expirado"}
	ForbiddenErrorResponse            = HTTPErrorResponse{HTTPStatusCode: 403, ErrorCode: "403", Message: "Usuário não autorizado"}
	NotFoundErrorResponse             = HTTPErrorResponse{HTTPStatusCode: 404, ErrorCode: "404", Message: "Recurso não encontrado"}
	MethodNotAllowedErrorResponse     = HTTPErrorResponse{HTTPStatusCode: 405, ErrorCode: "405", Message: "Método não suportado"}
	ConflictRequestErrorResponse      = HTTPErrorResponse{HTTPStatusCode: 409, ErrorCode: "409", Message: "Requisição com conflito"}
	UnsupportedMediaTypeErrorResponse = HTTPErrorResponse{HTTPStatusCode: 415, ErrorCode: "415", Message: "Tipo de mídia não suportado"}
	UnprocessableEntityErrorResponse  = HTTPErrorResponse{HTTPStatusCode: 422, ErrorCode: "422", Message: "Não foi possível processar as instruções contidas na requisição"}
	InternalServerErrorResponse       = HTTPErrorResponse{HTTPStatusCode: 500, ErrorCode: "500", Message: "Erro interno ao servidor"}
	ServiceUnavailableErrorResponse   = HTTPErrorResponse{HTTPStatusCode: 503, ErrorCode: "503", Message: "Serviço temporariamente indisponível"}
	GatewayUnavailableErrorResponse   = HTTPErrorResponse{HTTPStatusCode: 504, ErrorCode: "504", Message: "Serviço temporariamente indisponível"}
)

type RestrictionType string

const (
	RequiredRestriction  RestrictionType = "REQUIRED"
	InvalidRestriction   RestrictionType = "INVALID"
	BusinessRestriction  RestrictionType = "BUSINESS"
	MinLengthRestriction RestrictionType = "MIN_LENGTH"
	MaxLengthRestriction RestrictionType = "MAX_LENGTH"
)

type ValidationError struct {
	FieldName       string          `json:"fieldName"`
	RestrictionType RestrictionType `json:"restrictionType"`
	Message         string          `json:"message"`
}

func (v ValidationError) String() string {
	return fmt.Sprintf("validationError: %s:%s:%s", v.FieldName, v.RestrictionType, v.Message)
}

type AdditionalInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (a AdditionalInfo) String() string {
	return fmt.Sprintf("additionalInfo: %s:%s", a.Key, a.Value)
}

// HTTPErrorResponse
type HTTPErrorResponse struct {
	HTTPStatusCode   int               `json:"httpStatusCode"`
	ErrorCode        string            `json:"errorCode,omitempty"`
	Message          string            `json:"message,omitempty"`
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`
	AdditionalInfo   []AdditionalInfo  `json:"additionalInfo,omitempty"`
}

func (resp *HTTPErrorResponse) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("statusCode:%d, ", resp.HTTPStatusCode))
	s.WriteString(fmt.Sprintf("errorCode:%s, ", resp.ErrorCode))
	s.WriteString(fmt.Sprintf("message:%s, ", resp.Message))
	for _, validationError := range resp.ValidationErrors {
		s.WriteString(fmt.Sprintf("%s, ", validationError.String()))
	}
	for _, info := range resp.AdditionalInfo {
		s.WriteString(fmt.Sprintf("%s, ", info.String()))
	}
	return s.String()
}

func NewBadRequestResponse(errMsg string) HTTPErrorResponse {
	return HTTPErrorResponse{HTTPStatusCode: 400, ErrorCode: "400", Message: errMsg}
}

func NewNotFoundResponse(err string) HTTPErrorResponse {
	return HTTPErrorResponse{HTTPStatusCode: 404, ErrorCode: "404", Message: err}
}

func NewConflictResponse(err string) HTTPErrorResponse {
	return HTTPErrorResponse{HTTPStatusCode: 409, ErrorCode: "409", Message: err}
}

func NewInternalServerResponse(err string) HTTPErrorResponse {
	return HTTPErrorResponse{HTTPStatusCode: 500, ErrorCode: "500", Message: err}
}

func UnprocessableEntityWithError(fieldName, message string, restrictionType RestrictionType) HTTPErrorResponse {
	response := UnprocessableEntityErrorResponse
	response.AddError(fieldName, message, restrictionType)
	return response
}
func (resp *HTTPErrorResponse) AddError(fieldName, message string, restrictionType RestrictionType) {
	err := ValidationError{
		FieldName:       fieldName,
		RestrictionType: restrictionType,
		Message:         message,
	}
	resp.ValidationErrors = append(resp.ValidationErrors, err)
}

func (resp *HTTPErrorResponse) HasError() bool {
	return len(resp.ValidationErrors) > 0
}

func (resp *HTTPErrorResponse) ToJSON() []byte {
	if resp == nil {
		return []byte{}
	}
	response, _ := json.Marshal(resp)
	return response
}
