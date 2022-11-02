package cns_errors

import "net/http"

func init() {
	// AUTH ERRORS
	cnsErrorInfos[ERR_AUTHORIZATION_COMMON] = CnsErrorInfo{"Unauthorized", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_JWT] = CnsErrorInfo{"Unauthorized, invalid JWT token", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_TOKEN_MISSING] = CnsErrorInfo{"Authorization header is empty", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_TOKEN_INVALID] = CnsErrorInfo{"Authorization token is invalid", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_TOKEN_EXPIRED] = CnsErrorInfo{"Token is expired", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_NOT_KYC] = CnsErrorInfo{"KYC required", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_NO_PERMISSION] = CnsErrorInfo{"No permission to access", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_NO_AUTH_HEADER] = CnsErrorInfo{"Auth header is empty", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_MISS_EXP_FIELD] = CnsErrorInfo{"Missing exp field", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_EXP_FORMAT] = CnsErrorInfo{"Exp must be float64 format", http.StatusUnauthorized}
	cnsErrorInfos[ERR_AUTHORIZATION_JWT_PAYLOAD] = CnsErrorInfo{"Jwt payload content uncorrect", http.StatusUnauthorized}

	// VALIDATION ERRORS
	cnsErrorInfos[ERR_INVALID_REQUEST_COMMON] = CnsErrorInfo{"Invalid request", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_APP_ID] = CnsErrorInfo{"Invalid app id", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_ADDRESS] = CnsErrorInfo{"Invalid address", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_CHAIN] = CnsErrorInfo{"Chain is not supported", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_CONTRACT_TYPE] = CnsErrorInfo{"Contract type is not supported", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_URL] = CnsErrorInfo{"Invalid url", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_METADATA_ID] = CnsErrorInfo{"Invalid metadataId", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_MINT_AMOUNT] = CnsErrorInfo{"Invalid mint amount, mint amount could not be 0", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_MINT_AMOUNT_721] = CnsErrorInfo{"Invalid mint amount, mint amount could not more than 1 for erc 721 contract", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_TOKEN_ID] = CnsErrorInfo{"Invalid token ID", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_CONTRACT_TYPE_UNMATCH] = CnsErrorInfo{"Contract type and contract address not match", http.StatusBadRequest}
	cnsErrorInfos[ERR_INVALID_PAGINATION] = CnsErrorInfo{"Invalid page or limit", http.StatusBadRequest}

	// CONFLICT ERRORS
	cnsErrorInfos[ERR_CONFLICT_COMMON] = CnsErrorInfo{"Conflict", http.StatusConflict}
	cnsErrorInfos[ERR_CONFLICT_COMPANY_EXISTS] = CnsErrorInfo{"Company already exists", http.StatusConflict}

	// RATELIMIT ERRORS
	cnsErrorInfos[ERR_TOO_MANY_REQUEST_COMMON] = CnsErrorInfo{"Too many requests", http.StatusTooManyRequests}

	// INTERNAL SERVER ERRORS
	cnsErrorInfos[ERR_INTERNAL_SERVER_COMMON] = CnsErrorInfo{"Internal Server error", http.StatusInternalServerError}
	cnsErrorInfos[ERR_INTERNAL_SERVER_DB] = CnsErrorInfo{"Database operation error", http.StatusInternalServerError}

	// BUSINESS ERRORS
	cnsErrorInfos[ERR_BUSINESS_COMMON] = CnsErrorInfo{"Business error", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_MINT_LIMIT_EXCEEDED] = CnsErrorInfo{"Mint limit exceeded", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_DEPLOY_LIMIT_EXCEEDED] = CnsErrorInfo{"Deploy limit exceeded", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_UPLOADE_FILE_LIMIT_EXCEEDED] = CnsErrorInfo{"Uploade file limit exceeded", HTTP_STATUS_BUSINESS_ERROR}

	cnsErrorInfos[ERR_NO_SPONSOR] = CnsErrorInfo{"Contract has no sponsor", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_NO_SPONSOR_BALANCE] = CnsErrorInfo{"Contract sponsor balance not enough", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_NO_SPONSOR_FOR_USER] = CnsErrorInfo{"Contract has no sponsor for application admin", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_NO_PERMISSION_TO_UPDATE_ADMIN] = CnsErrorInfo{"Only admin can reset admin", HTTP_STATUS_BUSINESS_ERROR}
	cnsErrorInfos[ERR_CONTRACT_NOT_OWNED_BY_APP] = CnsErrorInfo{"Contract is not belong to this application", HTTP_STATUS_BUSINESS_ERROR}
}

const (
	HTTP_STATUS_BUSINESS_ERROR = 599
)

// AUTH ERRORS
const (
	ERR_AUTHORIZATION_COMMON CnsError = http.StatusUnauthorized*100 + iota //40100
	ERR_AUTHORIZATION_JWT
	ERR_AUTHORIZATION_TOKEN_MISSING
	ERR_AUTHORIZATION_TOKEN_INVALID
	ERR_AUTHORIZATION_TOKEN_EXPIRED
	ERR_AUTHORIZATION_NOT_KYC
	ERR_AUTHORIZATION_NO_PERMISSION
	ERR_AUTHORIZATION_NO_AUTH_HEADER
	ERR_AUTHORIZATION_MISS_EXP_FIELD
	ERR_AUTHORIZATION_EXP_FORMAT
	ERR_AUTHORIZATION_JWT_PAYLOAD
)

// VALIDATION ERRORS
const (
	ERR_INVALID_REQUEST_COMMON CnsError = http.StatusBadRequest*100 + iota //40000
	ERR_INVALID_APP_ID
	ERR_INVALID_ADDRESS
	ERR_INVALID_CHAIN
	ERR_INVALID_CONTRACT_TYPE
	ERR_INVALID_URL
	ERR_INVALID_METADATA_ID
	ERR_INVALID_MINT_AMOUNT
	ERR_INVALID_MINT_AMOUNT_721
	ERR_INVALID_TOKEN_ID
	ERR_INVALID_CONTRACT_TYPE_UNMATCH
	ERR_INVALID_PAGINATION
)

// RESOURCE CONFLICT ERRORS
const (
	ERR_CONFLICT_COMMON CnsError = http.StatusConflict*100 + iota //40900
	ERR_CONFLICT_COMPANY_EXISTS
)

// RATELIMIT ERRORS
const (
	ERR_TOO_MANY_REQUEST_COMMON CnsError = http.StatusTooManyRequests*100 + iota //42900
)

// INTERNAL SERVER ERRORS
const (
	ERR_INTERNAL_SERVER_COMMON CnsError = http.StatusInternalServerError*100 + iota //50000
	ERR_INTERNAL_SERVER_DB
	ERR_INTERNAL_SERVER_DB_NOT_FOUND
)

// BUSINESS ERRORS
const (
	ERR_BUSINESS_COMMON CnsError = HTTP_STATUS_BUSINESS_ERROR*100 + iota //60000
	ERR_NO_SPONSOR
	ERR_NO_SPONSOR_BALANCE
	ERR_NO_SPONSOR_FOR_USER
	ERR_MINT_LIMIT_EXCEEDED
	ERR_DEPLOY_LIMIT_EXCEEDED
	ERR_UPLOADE_FILE_LIMIT_EXCEEDED
	ERR_NO_PERMISSION_TO_UPDATE_ADMIN
	ERR_CONTRACT_NOT_OWNED_BY_APP
)

func GetRainbowOthersErrCode(httpStatusCode int) int {
	return httpStatusCode * 100
}
