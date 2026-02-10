package app

import (
	"comparei-servico-proconf/internal/domain/proconf"
	proconf_interface "comparei-servico-proconf/internal/domain/proconf/interface"
	"log"
)

type ProconfService struct {
	proconfRepo proconf_interface.ProconfRepository
}

func NewProconfService(proconfRepo proconf_interface.ProconfRepository) *ProconfService {
	return &ProconfService{proconfRepo: proconfRepo}
}

func (s *ProconfService) Create(proconf *proconf.Proconf) (*proconf.Proconf, error) {
	log.Println("EXEC: service.Create")
	proconf, err := s.proconfRepo.Create(proconf)
	return proconf, err
}

func (s *ProconfService) GetMercadoProdutoByID(id int) (*proconf.Proconf, error) {
	log.Println("EXEC: service.GetMercadoProdutoByID")
	proconf, err := s.proconfRepo.GetMercadoProdutoByID(id)
	return proconf, err
}

}
