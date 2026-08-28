package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB struct {
	Conn *pgx.Conn
}

func (db *DB) RunMigrations(ctx context.Context) error {
	schema := `CREATE TABLE IF NOT EXISTS mailboxes (
        id SERIAL PRIMARY KEY,
        server_hostport TEXT NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        primary_inbox_name TEXT NOT NULL,
        last_uid INTEGER NOT NULL DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS emails (
        id SERIAL PRIMARY KEY,
        mailbox_id INT NOT NULL REFERENCES mailboxes(id),
        uid INTEGER NOT NULL,
        message_id TEXT NOT NULL,
        subject TEXT NOT NULL,
        from_address TEXT NOT NULL,
        received_at TIMESTAMPTZ NOT NULL,
        body TEXT,
        UNIQUE(mailbox_id, uid)
    );
    `
	_, err := db.Conn.Exec(ctx, schema)
	return err
}

// get last synced uid
func (db *DB) GetLastSyncedUID(ctx context.Context, tx *pgx.Tx, mailboxID int) (uint32, error) {
	var lastUID uint32
	err := connOrTx(db, tx).QueryRow(ctx, "SELECT last_uid FROM mailboxes WHERE id = $1", mailboxID).Scan(&lastUID)
	if err != nil {
		return 0, err
	}
	return lastUID, nil
}

// update last synced uid (pushing up the marker)
func (db *DB) UpdateLastSyncedUID(ctx context.Context, tx *pgx.Tx, mailboxID int, lastUID uint32) error {
	_, err := connOrTx(db, tx).Exec(ctx, "UPDATE mailboxes SET last_uid = $1 WHERE id = $2", lastUID, mailboxID)
	if err != nil {
		return err
	}
	return nil
}

// add an email to the db; after this u probably want to update the last synced uid
func (db *DB) InsertEmail(ctx context.Context, tx *pgx.Tx, email Email) error {
	_, err := connOrTx(db, tx).Exec(ctx,
		"INSERT INTO emails (mailbox_id, uid, message_id, subject, from_address, received_at, body) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		email.Mailbox.ID,
		email.UID,
		email.MessageID,
		email.Subject,
		email.FromAddress,
		email.RecievedAt,
		email.Body,
	)
	return err
}

func (db *DB) InsertEmailandUpdateLastSyncedUID(ctx context.Context, email Email) error {
	// start a transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// insert the email
	err = db.InsertEmail(ctx, &tx, email)
	if err != nil {
		return err
	}

	// update the last synced UID
	err = db.UpdateLastSyncedUID(ctx, &tx, email.Mailbox.ID, email.UID)
	if err != nil {
		return err
	}

	// commit the transaction
	return tx.Commit(ctx)
}

func (db *DB) InsertMailbox(ctx context.Context, tx *pgx.Tx, mailbox Mailbox) (int, error) {
	var id int
	err := connOrTx(db, tx).QueryRow(ctx,
		"INSERT INTO mailboxes (server_hostport, username, password, primary_inbox_name) VALUES ($1, $2, $3, $4) RETURNING id",
		mailbox.ServerHostport,
		mailbox.Username,
		mailbox.Password,
		mailbox.PrimaryInbox,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DB) ListMailboxes(ctx context.Context, tx *pgx.Tx) ([]Mailbox, error) {
	rows, err := connOrTx(db, tx).Query(ctx, "SELECT id, server_hostport, username, password, primary_inbox_name FROM mailboxes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mailboxes []Mailbox
	for rows.Next() {
		var mailbox Mailbox
		err := rows.Scan(&mailbox.ID, &mailbox.ServerHostport, &mailbox.Username, &mailbox.Password, &mailbox.PrimaryInbox)
		if err != nil {
			return nil, err
		}
		mailboxes = append(mailboxes, mailbox)
	}
	return mailboxes, nil
}

func (db *DB) GetMailbox(ctx context.Context, tx *pgx.Tx, id int) (*Mailbox, error) {
	var mailbox Mailbox
	err := connOrTx(db, tx).QueryRow(ctx,
		"SELECT id, server_hostport, username, password, primary_inbox_name FROM mailboxes WHERE id = $1",
		id,
	).Scan(&mailbox.ID, &mailbox.ServerHostport, &mailbox.Username, &mailbox.Password, &mailbox.PrimaryInbox)
	if err != nil {
		return nil, err
	}
	return &mailbox, nil
}

// list emails for a given mailbox without their body
func (db *DB) ListEmails(ctx context.Context, tx *pgx.Tx, mailboxID int) ([]Email, error) {
	rows, err := connOrTx(db, tx).Query(ctx, "SELECT id, mailbox_id, uid, message_id, subject, from_address, received_at FROM emails WHERE mailbox_id = $1", mailboxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.ID, &email.Mailbox.ID, &email.UID, &email.MessageID, &email.Subject, &email.FromAddress, &email.RecievedAt)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

// get email with body
func (db *DB) GetEmail(ctx context.Context, tx *pgx.Tx, mailboxID int, messageID string) (*Email, error) {
	var email Email
	err := connOrTx(db, tx).QueryRow(ctx,
		"SELECT id, mailbox_id, uid, message_id, subject, from_address, received_at, body FROM emails WHERE mailbox_id = $1 AND message_id = $2",
		mailboxID,
		messageID,
	).Scan(&email.ID, &email.Mailbox.ID, &email.UID, &email.MessageID, &email.Subject, &email.FromAddress, &email.RecievedAt, &email.Body)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

type QueryAndExec interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// selects either the transaction if not nil, otherwise default to the database connection
func connOrTx(db *DB, tx *pgx.Tx) QueryAndExec {
	if tx != nil {
		return *tx
	}
	return db.Conn
}
