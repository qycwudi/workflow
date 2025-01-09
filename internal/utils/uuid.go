package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	// 去掉-
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
