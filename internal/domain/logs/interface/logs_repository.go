package logs_interface

import "comparei-servico-proconf/internal/domain/logs"

type LogsRepository interface {
	CreateLogsConfirmacao(logs *logs.LogsConfirmacao) (*logs.LogsConfirmacao, error)
}
