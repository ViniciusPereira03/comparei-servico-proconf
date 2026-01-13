package logs_interface

import "comparei-servico-promer/internal/domain/logs"

type LogsRepository interface {
	CreateLogsConfirmacao(logs *logs.LogsConfirmacao) (*logs.LogsConfirmacao, error)
}
