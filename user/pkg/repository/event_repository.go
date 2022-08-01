// Holds EventRepository definition and it's CRUD implementation
package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const EventsTable string = "events"

type EventRepository struct {
	db *sqlx.DB
}

type Event struct {
	Id      uint   `db:"id"`
	User_id uint   `db:"user_id"`
	Title   string `db:"title"`
	Body    string `db:"body"`
}

func MakeEventRepository() EventRepository {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USERNAME"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   "mariadb:3306",
		DBName: "articles",
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return EventRepository{
		db: db,
	}
}

func (s EventRepository) Index() ([]Event, error) {
	sql := fmt.Sprintf(`SELECT * FROM %s`, EventsTable)
	articles := []Event{}
	err := s.db.Select(&articles, sql)
	return articles, err
}

func (s EventRepository) Find(id string) (Event, error) {
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, EventsTable)
	user := Event{}
	stmt, err := s.db.Preparex(sql)
	if err != nil {
		return Event{}, err
	}

	err = stmt.Get(&user, sql, id)
	return user, err
}

func (s EventRepository) Create(data Event) error {
	sql := fmt.Sprintf(`INSERT INTO %s (name) VALUES (:name)`, EventsTable)
	stmt, err := s.db.PrepareNamed(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data)
	return err
}

func (s EventRepository) Update(data Event) error {
	sql := fmt.Sprintf(`UPDATE %s SET name = :name WHERE id = :id`, EventsTable)
	stmt, err := s.db.PrepareNamed(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data)
	return err
}

func (s EventRepository) Delete(id string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, EventsTable)
	stmt, err := s.db.Preparex(sql)
	if err != nil {
		return err
	}
	stmt.Exec(id)
	return err
}
