package domain

type ContextKey string

// AuthContextKey is the key for the users context info
var AuthContextKey ContextKey = "user"

// CorrelationIDCtxKey is the key for the correlation id
var CorrelationIDCtxKey ContextKey = "correlation_id"
