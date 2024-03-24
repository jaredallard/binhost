package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Pkg holds the schema definition for the Pkg entity.
type Pkg struct {
	ent.Schema
}

// Fields of the Pkg.
func (Pkg) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
	}
}
