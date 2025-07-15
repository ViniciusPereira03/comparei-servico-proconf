package proconf_interface

import "comparei-servico-promer/internal/domain/proconf"

type ProconfRepository interface {
	Create(user *proconf.Proconf) (*proconf.Proconf, error)
}
