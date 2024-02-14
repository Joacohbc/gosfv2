package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.String("username").
			NotEmpty().
			Unique().
			StructTag(`json:"username"`),
		field.String("password").
			NotEmpty(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
		edge.From("shared_files", File.Type).
			Ref("shared_with"),
		edge.From("shared_notes", Note.Type).
			Ref("shared_with"),
		edge.To("notes", Note.Type),
	}
}
