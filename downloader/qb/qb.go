package qb

import (
	"context"
)

type (
	Config struct {
		Ctx      context.Context
		Host     string
		Username string
		Password string
	}
	Client struct {
		config Config
		ck     string
	}
)
