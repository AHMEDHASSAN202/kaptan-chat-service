package redis

import (
	"context"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/utils"
	"kaptan/pkg/validators"
	"net/http"
	"time"
)

func Lock(ctx context.Context, redisClient *RedisClient, log logger.ILogger, lockKey string, lockDuration time.Duration) validators.ErrorResponse {
	var success bool
	var errLock error

	// Try to acquire the lock
	unexpectedError := utils.TryCatch(func() {
		success, errLock = redisClient.Lock(lockKey, "locked", lockDuration)
		if errLock != nil {
			// Log the error for monitoring purposes
			log.Errorf("Failed to acquire lock for key %s: %s", lockKey, errLock)
		}
	})()

	// Handle unexpected errors during lock acquisition
	if unexpectedError != nil {
		// Log the error for monitoring purposes
		log.Errorf("UnexpectedError Failed to acquire lock for key %s: %s", lockKey, unexpectedError.Error())
	}
	if errLock != nil || unexpectedError != nil {
		return validators.ErrorResponse{}
	}

	// If lock was not acquired successfully, return an error response
	if !success {
		// The lock was not acquired
		return validators.GetErrorResponse(&ctx, localization.LockError, nil, utils.GetAsPointer(http.StatusTooManyRequests))
	}

	return validators.ErrorResponse{}
}

func UnLock(ctx context.Context, redisClient *RedisClient, log logger.ILogger, lockKey string) {
	// Ensure the lock is released after handling the request
	unexpectedErrorUnlock := utils.TryCatch(func() {
		errUnlock := redisClient.Unlock(lockKey)
		if errUnlock != nil {
			// Log the error for monitoring purposes
			log.Errorf("Failed to release lock for key %s: %s", lockKey, errUnlock.Error())
		}
	})()
	// Handle unexpected errors during lock release
	if unexpectedErrorUnlock != nil {
		// Log the error for monitoring purposes
		log.Errorf("UnexpectedError Failed to release lock for key %s: %s", lockKey, unexpectedErrorUnlock.Error())
	}
}
