package proconf_interface

import "comparei-servico-proconf/internal/domain/proconf"

type ProconfRepository interface {
	Create(user *proconf.Proconf) (*proconf.Proconf, error)
	GetMercadoProdutoByID(id int) (*proconf.Proconf, error)
}
