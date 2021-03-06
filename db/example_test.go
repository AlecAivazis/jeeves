// Code generated by entc, DO NOT EDIT.

package db

import (
	"context"
	"log"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 DB_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleBankItem() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the bankitem's edges.

	// create bankitem vertex with its edges.
	bi := client.BankItem.
		Create().
		SetItemID("string").
		SetQuantity(1).
		SaveX(ctx)
	log.Println("bankitem created:", bi)

	// query edges.

	// Output:
}
func ExampleGuild() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the guild's edges.
	gb0 := client.GuildBank.
		Create().
		SetChannelID("string").
		SetBalance(1).
		SaveX(ctx)
	log.Println("guildbank created:", gb0)

	// create guild vertex with its edges.
	gu := client.Guild.
		Create().
		SetDiscordID("string").
		SetBank(gb0).
		SaveX(ctx)
	log.Println("guild created:", gu)

	// query edges.
	gb0, err = gu.QueryBank().First(ctx)
	if err != nil {
		log.Fatalf("failed querying bank: %v", err)
	}
	log.Println("bank found:", gb0)

	// Output:
}
func ExampleGuildBank() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the guildbank's edges.
	bi0 := client.BankItem.
		Create().
		SetItemID("string").
		SetQuantity(1).
		SaveX(ctx)
	log.Println("bankitem created:", bi0)

	// create guildbank vertex with its edges.
	gb := client.GuildBank.
		Create().
		SetChannelID("string").
		SetBalance(1).
		AddItems(bi0).
		SaveX(ctx)
	log.Println("guildbank created:", gb)

	// query edges.
	bi0, err = gb.QueryItems().First(ctx)
	if err != nil {
		log.Fatalf("failed querying items: %v", err)
	}
	log.Println("items found:", bi0)

	// Output:
}
