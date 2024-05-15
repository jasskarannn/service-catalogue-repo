package helpers

import (
	"time"

	"github.com/google/uuid"
)

func HashedUUIDAsTimestamp() uuid.UUID{
	currentTime := time.Now()
	// Convert the timestamp to a byte slice
	currentTimeInBytes := []byte(currentTime.Format(time.RFC3339))
	// Hash the timestamp bytes to generate a UUID
	hashedUUID := uuid.NewSHA1(uuid.Nil, currentTimeInBytes)
	
	return hashedUUID
}
