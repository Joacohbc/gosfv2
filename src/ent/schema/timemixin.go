package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimeDetails holds the schema definition for the TimeDetails entity.
type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"created_at"`),
		field.Time("updated_at").
			UpdateDefault(time.Now).
			Optional().
			Nillable().
			StructTag(`json:"updated_at,omitempty"`),
	}
}
