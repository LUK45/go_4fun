package main

import (
	"context"
	"log"
	"time"
	"database/sql"
)

type Store interface {
	CreateContact(contact *Contact) error
	CreateTicket(ticket *Ticket) error
	GetContacts() ([]Contact, error)
	GetTickets() ([]Ticket, error)
	CreateTables() error
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateContact(c *Contact) error {

	//query := "INSERT INTO contacts(id, first_name, last_name, email, phone) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE first_name=?, last_name=?, email=?, phone=?"

	query := "REPLACE INTO contacts(id, first_name, last_name, email, phone) VALUES (?,?,?,?,?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := store.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, c.Id, c.Properties.FirstName, c.Properties.LastName, c.Properties.Email, c.Properties.Phone)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	return nil
}

func (store *dbStore) CreateTicket(t *Ticket) error {
	query := "REPLACE INTO tickets(id, content, hubspot_owner_id) VALUES (?,?,?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := store.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, t.Id, t.Properties.Content, t.Properties.HubspotOwnerId)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	return nil
}

func (store *dbStore) GetContacts() ([]Contact, error) {
	rows, err := store.db.Query("SELECT * from contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contacts := []Contact{}
	for rows.Next() {
		c := Contact{}
		if err := rows.Scan(&c.Id, &c.Properties.FirstName, &c.Properties.LastName, &c.Properties.Email, &c.Properties.Phone); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (store *dbStore) GetTickets() ([]Ticket, error) {
	rows, err := store.db.Query("SELECT * from tickets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tickets := []Ticket{}
	for rows.Next() {
		t := Ticket{}
		if err := rows.Scan(&t.Id, &t.Properties.Content, &t.Properties.HubspotOwnerId ); err != nil {
			return nil, err

		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

var store Store

func InitStore(s Store) {
	store = s
}

func (store *dbStore) CreateTables() error {
	query := [2]string {"create table if not exists contacts ( id serial primary key , first_name varchar(256), last_name varchar(256), email varchar(256), phone varchar(256))",
		"create table if not exists tickets ( id serial primary key , content varchar(256), hubspot_owner_id varchar(256));"}
	for _, q := range query {
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err := store.db.ExecContext(ctx, q)
		if err != nil {
			log.Printf("Error %s when creating product table", err)
			return err
		}
	}
	return nil
}
