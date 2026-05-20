package usecase_errors

import (
	"errors"
	"net/http"

	"github.com/BramAristyo/mawish-pos/server/internal/api/validation"
	"github.com/jackc/pgx/v5/pgconn"
)

type UseCaseError struct {
	Code    int
	Message string
}

func (e *UseCaseError) Error() string {
	return e.Message
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

type CustomFieldErrors []validation.ValidationError

func (e *CustomFieldErrors) Error() string {
	return "custom field validation failed"
}

var (
	EmailExist         = &UseCaseError{http.StatusConflict, "email already registered"}
	EmailNotFound      = &UseCaseError{http.StatusConflict, "email doesn't registered yet"}
	DuplicateEntry     = &UseCaseError{http.StatusConflict, "data already exist"}
	UserNotActive      = &UseCaseError{http.StatusForbidden, "user is not active"}
	InvalidPassword    = &UseCaseError{http.StatusUnauthorized, "invalid email or password"}
	TokenRequired      = &UseCaseError{http.StatusUnauthorized, "token required"}
	TokenExpired       = &UseCaseError{http.StatusUnauthorized, "token expired"}
	TokenInvalid       = &UseCaseError{http.StatusUnauthorized, "token invalid"}
	InvalidID          = &UseCaseError{http.StatusBadRequest, "invalid ID format"}
	NotFound           = &UseCaseError{http.StatusNotFound, "resource not found"}
	ShiftAlreadyOpen   = &UseCaseError{http.StatusBadRequest, "user already has an open shift"}
	NoOpenShift        = &UseCaseError{http.StatusNotFound, "no open shift found for this user"}
	ForbiddenAccess    = &UseCaseError{http.StatusForbidden, "you don't have access to this resource"}
	InvalidOrderItem   = &UseCaseError{http.StatusBadRequest, "order item must have either product or bundling"}
	EmptyOrderItems    = &UseCaseError{http.StatusBadRequest, "order must contain at least one item"}
	OrderAlreadyVoided = &UseCaseError{http.StatusBadRequest, "order already voided"}
	DateFilterRequired = &UseCaseError{http.StatusBadRequest, "date filter (from/to) is required"}
	LedgerRecordFailed = &UseCaseError{http.StatusInternalServerError, "transaction saved but failed to record ledger entry"}
	PhotoNotFound      = &UseCaseError{http.StatusNotFound, "photo not found, please upload again"}
	PhotoVerifyError   = &UseCaseError{http.StatusInternalServerError, "failed to verify photo upload"}
)
