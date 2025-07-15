package users

import "time"

type User struct {
	ID         string     `bson:"_id,omitempty" json:"proconf_user_id"`
	UsuarioID  string     `bson:"id_usuario" json:"id"`
	Level      int        `bson:"level" json:"level"`
	CreatedAt  time.Time  `bson:"created_at"       json:"created_at"`
	ModifiedAt time.Time  `bson:"modified_at"      json:"modified_at"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
