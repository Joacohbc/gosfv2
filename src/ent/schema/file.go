package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Mixin of the File.
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.String("filename").
			NotEmpty().
			MaxLen(255).
			StructTag(`json:"filename"`),
		field.Bool("is_dir").
			Default(false).
			StructTag(`json:"is_dir"`),
		field.Bool("is_shared").
			Default(false).
			StructTag(`json:"shared"`),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("files").
			Required().
			Unique(),
		edge.To("shared_with", User.Type),
		edge.To("parent", File.Type).
			From("children").
			Unique(),
		edge.From("notes", Note.Type).
			Ref("files").
			Unique(),
	}
}
