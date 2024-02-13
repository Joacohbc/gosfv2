package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Note holds the schema definition for the Note entity.
type Note struct {
	ent.Schema
}

// Mixin of the File.
func (Note) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Note.
func (Note) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.String("title").
			NotEmpty().
			MaxLen(255).
			StructTag(`json:"title"`),
		field.String("content").
			NotEmpty().
			StructTag(`json:"content"`),
	}
}

// Edges of the Note.
func (Note) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("notes").
			Required().
			Unique(),
		edge.To("shared_with", User.Type),
		edge.To("parent", Note.Type).
			From("children").
			Unique(),
		edge.To("files", File.Type),
	}
}
