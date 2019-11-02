package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/edge"
)

// GuildChannel holds the schema definition for the BankItem entity.
type GuildChannel struct {
	ent.Schema
}

// Fields of the GuildChannel.
func (GuildChannel) Fields() []ent.Field {
	return []ent.Field{
		field.String("channel"),
		field.String("role"),
	}
}

// Edges of the GuildChannel.
func (GuildChannel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).Ref("channels").Unique(),
	}			
}
