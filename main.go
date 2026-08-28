package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	db := &DB{
		Conn: conn,
	}
	if err = db.RunMigrations(ctx); err != nil {
		panic(err)
	}

	mailbox := Mailbox{
		ID:             1,
		ServerHostport: os.Getenv("IMAP_HOSTPORT"),
		Username:       os.Getenv("IMAP_USERNAME"),
		Password:       os.Getenv("IMAP_PASSWORD"),
		PrimaryInbox:   "INBOX",
	}

	// if id, err := db.InsertMailbox(ctx, nil, mailbox); err != nil {
	// 	panic(err)
	// }
	// mailbox.ID = id

	if err = syncMailbox(ctx, db, mailbox); err != nil {
		panic(err)
	}
}
