package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("repository"),
		field.String("category"),
		field.String("name"),
		field.String("version"),
		field.UUID("target_id", uuid.UUID{}),
	}
}

func (Pkg) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("repository", "category", "name", "version", "target_id").Unique(),
	}
}

func (Pkg) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("target", Target.Type).Unique().Field("target_id").Required(),
	}
}
