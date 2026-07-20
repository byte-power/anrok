package anrok

import (
	"fmt"
)

// ErrorType represents the Anrok API error type string returned in {"type": "..."} responses.
type ErrorType string

const (
	ErrAccountingTimeZoneNotSetForSeller ErrorType = "accountingTimeZoneNotSetForSeller"
	ErrAccountingTimeZoneNotSupported    ErrorType = "accountingTimeZoneNotSupported"
	ErrCertificateIdNotFound             ErrorType = "certificateIdNotFound"
	ErrCurrencyCodeNotSupported          ErrorType = "currencyCodeNotSupported"
	ErrCustomerAddressCouldNotResolve    ErrorType = "customerAddressCouldNotResolve"
	ErrCustomerIdNotFound                ErrorType = "customerIdNotFound"
	ErrDuplicateJurisIds                 ErrorType = "duplicateJurisIds"
	ErrExternalServiceError              ErrorType = "externalServiceError"
	ErrProductExternalIdUnknown          ErrorType = "productExternalIdUnknown"
	ErrTaxDateTooFarInFuture             ErrorType = "taxDateTooFarInFuture"
	ErrTransactionFrozenForFiling        ErrorType = "transactionFrozenForFiling"
)

// RequestError represents a client-side failure (serialization, network, deserialization).
type RequestError struct {
	Op  string // operation that failed, e.g. "marshal request", "send request"
	Err error  // underlying error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("anrok: %s: %v", e.Op, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// APIError represents a non-2xx HTTP response from the Anrok API.
// It is also the base embedded in RateLimitError and TypedError, so
// errors.As(err, &apiErr) will match any of those subtypes.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("anrok: HTTP %d: %s", e.StatusCode, e.Body)
}

// RateLimitError represents a 429 Too Many Requests response.
// RetryAfter is the number of seconds the caller should wait before retrying,
// as specified by the Retry-After response header.
type RateLimitError struct {
	APIError
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("anrok: rate limited (HTTP 429), retry after %d seconds", e.RetryAfter)
}

// As allows errors.As(err, &apiErr) to match a *RateLimitError as *APIError.
func (e *RateLimitError) As(target any) bool {
	if t, ok := target.(**APIError); ok {
		*t = &e.APIError
		return true
	}
	return false
}

// TypedError represents an Anrok business-logic error whose response body
// contains {"type": "<errorType>"}.
type TypedError struct {
	APIError
	Type ErrorType
}

func (e *TypedError) Error() string {
	return fmt.Sprintf("anrok: %s (HTTP %d)", e.Type, e.StatusCode)
}

// IsType reports whether this error matches the given ErrorType.
func (e *TypedError) IsType(t ErrorType) bool {
	return e.Type == t
}

// As allows errors.As(err, &apiErr) to match a *TypedError as *APIError.
func (e *TypedError) As(target any) bool {
	if t, ok := target.(**APIError); ok {
		*t = &e.APIError
		return true
	}
	return false
}
