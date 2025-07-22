package common

import "github.com/hashicorp/go-uuid"

// Function to generate UUID
func GenerateUuid() (string, error) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	return uuid, nil
}
