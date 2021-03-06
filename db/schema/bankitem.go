package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// BankItem holds the schema definition for the BankItem entity.
type BankItem struct {
	ent.Schema
}

// Fields of the BankItem.
func (BankItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("itemID"),
		field.Int("quantity").Positive(),
	}
}

// Edges of the BankItem.
func (BankItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bank", GuildBank.Type).Ref("items").Unique(),
	}
}
