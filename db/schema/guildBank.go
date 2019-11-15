package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/edge"
)

// GuildBank holds the schema definition for the BankItem entity.
type GuildBank struct {
	ent.Schema
}

// Fields of the GuildBank.
func (GuildBank) Fields() []ent.Field {
	return []ent.Field{
		field.String("channelID"),
		field.String("displayMessageID").Unique(),
	}
}

// Edges of the GuildBank.
func (GuildBank) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", BankItem.Type),
		edge.From("guild", Guild.Type).Ref("bank").Unique(),
	}
}
