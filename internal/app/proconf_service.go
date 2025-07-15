package app

import (
	"comparei-servico-promer/internal/domain/proconf"
	proconf_interface "comparei-servico-promer/internal/domain/proconf/interface"
	"log"
)

type ProconfService struct {
	proconfRepo proconf_interface.ProconfRepository
}

func NewProconfService(proconfRepo proconf_interface.ProconfRepository) *ProconfService {
	return &ProconfService{proconfRepo: proconfRepo}
}

func (s *ProconfService) Create(proconf *proconf.Proconf) error {
	log.Println("EXEC: service.Create")
	proconf, err := s.proconfRepo.Create(proconf)
	return err
}
