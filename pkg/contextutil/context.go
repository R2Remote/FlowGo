package contextutil

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
)

// GetUserID retrieves the UserID from the context.
// It checks both gin.Context (Keys) and standard context.Context (Values).
func GetUserID(ctx context.Context) (uint64, error) {
	// 1. If it's a *gin.Context, check the Keys map
	if c, ok := ctx.(*gin.Context); ok {
		if val, exists := c.Get(UserIDKey); exists {
			if id, ok := val.(uint64); ok {
				return id, nil
			}
		}
	}

	// 2. Fallback: Check standard context.Value
	val := ctx.Value(UserIDKey)
	if val != nil {
		if id, ok := val.(uint64); ok {
			return id, nil
		}
		// Also handle the string key if context.WithValue was used with the string "user_id"
		// (Though best practice is to use a custom type key, for simple integrations string is common)
		if id, ok := ctx.Value("user_id").(uint64); ok {
			return id, nil
		}
	}

	return 0, errors.New("user id not found in context")
}
