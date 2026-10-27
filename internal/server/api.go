package server

import (
	"github.com/IFA-01/messenger/internal/repository/queries"
)

type ApiConfig struct {
	DB *queries.Queries
}
