// Code generated by entc, DO NOT EDIT.

package guild

import (
	"github.com/AlecAivazis/jeeves/db/predicate"
	"github.com/facebookincubator/ent/dialect/sql"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldID), id))
		},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(ids) == 0 {
				s.Where(sql.False())
				return
			}
			v := make([]interface{}, len(ids))
			for i := range v {
				v[i] = ids[i]
			}
			s.Where(sql.In(s.C(FieldID), v...))
		},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(ids) == 0 {
				s.Where(sql.False())
				return
			}
			v := make([]interface{}, len(ids))
			for i := range v {
				v[i] = ids[i]
			}
			s.Where(sql.NotIn(s.C(FieldID), v...))
		},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldID), id))
		},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldID), id))
		},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldID), id))
		},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldID), id))
		},
	)
}

// DiscordID applies equality check predicate on the "discordID" field. It's identical to DiscordIDEQ.
func DiscordID(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDEQ applies the EQ predicate on the "discordID" field.
func DiscordIDEQ(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDNEQ applies the NEQ predicate on the "discordID" field.
func DiscordIDNEQ(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDIn applies the In predicate on the "discordID" field.
func DiscordIDIn(vs ...string) predicate.Guild {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Guild(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldDiscordID), v...))
		},
	)
}

// DiscordIDNotIn applies the NotIn predicate on the "discordID" field.
func DiscordIDNotIn(vs ...string) predicate.Guild {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Guild(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldDiscordID), v...))
		},
	)
}

// DiscordIDGT applies the GT predicate on the "discordID" field.
func DiscordIDGT(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDGTE applies the GTE predicate on the "discordID" field.
func DiscordIDGTE(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDLT applies the LT predicate on the "discordID" field.
func DiscordIDLT(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDLTE applies the LTE predicate on the "discordID" field.
func DiscordIDLTE(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDContains applies the Contains predicate on the "discordID" field.
func DiscordIDContains(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDHasPrefix applies the HasPrefix predicate on the "discordID" field.
func DiscordIDHasPrefix(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDHasSuffix applies the HasSuffix predicate on the "discordID" field.
func DiscordIDHasSuffix(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDEqualFold applies the EqualFold predicate on the "discordID" field.
func DiscordIDEqualFold(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldDiscordID), v))
		},
	)
}

// DiscordIDContainsFold applies the ContainsFold predicate on the "discordID" field.
func DiscordIDContainsFold(v string) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldDiscordID), v))
		},
	)
}

// HasChannels applies the HasEdge predicate on the "channels" edge.
func HasChannels() predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(ChannelsColumn).
						From(builder.Table(ChannelsTable)).
						Where(sql.NotNull(ChannelsColumn)),
				),
			)
		},
	)
}

// HasChannelsWith applies the HasEdge predicate on the "channels" edge with a given conditions (other predicates).
func HasChannelsWith(preds ...predicate.GuildChannel) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(ChannelsColumn).From(builder.Table(ChannelsTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasBank applies the HasEdge predicate on the "bank" edge.
func HasBank() predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(BankColumn).
						From(builder.Table(BankTable)).
						Where(sql.NotNull(BankColumn)),
				),
			)
		},
	)
}

// HasBankWith applies the HasEdge predicate on the "bank" edge with a given conditions (other predicates).
func HasBankWith(preds ...predicate.BankItem) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(BankColumn).From(builder.Table(BankTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Guild) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			for _, p := range predicates {
				p(s)
			}
		},
	)
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Guild) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			for i, p := range predicates {
				if i > 0 {
					s.Or()
				}
				p(s)
			}
		},
	)
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Guild) predicate.Guild {
	return predicate.Guild(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
