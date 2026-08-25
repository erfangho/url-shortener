package config

import (
	"os"
	"time"
)

type JWT struct {
}

func (j *JWT) JWTSecret() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}

func (j *JWT) TokenExpiry() time.Duration {
	return 24 * time.Hour
}
