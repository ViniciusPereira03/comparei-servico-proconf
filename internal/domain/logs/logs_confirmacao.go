package logs

import "time"

type LogsConfirmacao struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	UsuarioID        string    `bson:"usuario_id" json:"usuario_id"`
	NivelUsuario     int       `bson:"nivel_usuario" json:"nivel_usuario"`
	MercadoProdutoID int       `bson:"mercado_produto_id"       json:"mercado_produto_id"`
	CreatedAt        time.Time `bson:"created_at"       json:"created_at"`
}
