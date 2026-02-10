package proconf

import "time"

type Proconf struct {
	ID               string     `bson:"_id,omitempty" json:"id"`
	MercadoProdutoID int        `bson:"mercado_produto_id"       json:"mercado_produto_id"`
	MercadoID        int        `bson:"mercado_id"   json:"mercado_id"`
	ProdutoID        int        `bson:"produto_id"      json:"produto_id"`
	PrecoUnitario    float32    `bson:"preco_unitario"   json:"preco_unitario"`
	NivelConfianca   int        `bson:"nivel_confianca"     json:"nivel_confianca"`
	CreatedAt        time.Time  `bson:"created_at"       json:"created_at"`
	ModifiedAt       time.Time  `bson:"modified_at"      json:"modified_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ProconfConfirmValue struct {
	MercadoProdutoID int        `bson:"mercado_produto_id"       json:"mercado_produto_id"`
	MercadoID        int        `bson:"mercado_id"   json:"mercado_id"`
	ProdutoID        int        `bson:"produto_id"      json:"produto_id"`
	PrecoUnitario    float32    `bson:"preco_unitario"   json:"preco_unitario"`
	NivelConfianca   int        `bson:"nivel_confianca"     json:"nivel_confianca"`
	CreatedAt        time.Time  `bson:"created_at"       json:"created_at"`
	ModifiedAt       time.Time  `bson:"modified_at"      json:"modified_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (p *ProconfConfirmValue) ParseToProconf() *Proconf {
	return &Proconf{
		MercadoProdutoID: p.MercadoProdutoID,
		MercadoID:        p.MercadoID,
		ProdutoID:        p.ProdutoID,
		PrecoUnitario:    p.PrecoUnitario,
		NivelConfianca:   p.NivelConfianca,
		CreatedAt:        p.CreatedAt,
		ModifiedAt:       p.ModifiedAt,
		DeletedAt:        p.DeletedAt,
	}
}
