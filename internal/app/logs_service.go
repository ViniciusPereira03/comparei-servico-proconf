package app

import (
	"comparei-servico-promer/internal/domain/logs"
	logs_interface "comparei-servico-promer/internal/domain/logs/interface"
)

type LogsService struct {
	logRepo logs_interface.LogsRepository
}

func NewLogsService(logRepo logs_interface.LogsRepository) *LogsService {
	return &LogsService{logRepo: logRepo}
}

func (s *LogsService) CreateLogsConfirmacao(logs *logs.LogsConfirmacao) error {
	_, err := s.logRepo.CreateLogsConfirmacao(logs)
	return err
}
