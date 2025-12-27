package repository

import "github.com/RahulSingh9131/vector/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
