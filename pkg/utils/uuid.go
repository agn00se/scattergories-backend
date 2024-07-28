package utils

import "github.com/google/uuid"

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}
