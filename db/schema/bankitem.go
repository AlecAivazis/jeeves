package schema

import "github.com/facebookincubator/ent"

// BankItem holds the schema definition for the BankItem entity.
type BankItem struct {
	ent.Schema
}

// Fields of the BankItem.
func (BankItem) Fields() []ent.Field {
	return nil
}

// Edges of the BankItem.
func (BankItem) Edges() []ent.Edge {
	return nil
}
