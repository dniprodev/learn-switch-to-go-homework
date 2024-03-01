package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(postgresURI string) (*Repository, error) {
	conn, err := pgx.Connect(context.Background(), postgresURI)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = waitForConnection(conn)

	if err != nil {
		return nil, err
	}

	repo := &Repository{conn: conn}
	return repo, nil
}

func (r *Repository) Close() {
	r.conn.Close(context.Background())
}

func (r *Repository) Initialize() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(255) PRIMARY KEY,
            name TEXT NOT NULL,
            password TEXT NOT NULL
        );`

	_, err := r.conn.Exec(context.Background(), query)

	return err
}

func (r *Repository) Save(user User) error {
	sqlStatement := `INSERT INTO users (id, name, password) VALUES ($1, $2, $3)`
	_, err := r.conn.Exec(context.Background(), sqlStatement, user.ID, user.Name, user.Password)
	return err
}

func (r *Repository) FindByUsername(name string) (User, error) {
	sqlStatement := `SELECT id, name, password FROM users WHERE name=$1`
	row := r.conn.QueryRow(context.Background(), sqlStatement, name)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("fail to find user: %w", err)
		}
		return User{}, err
	}
	return user, nil
}

func waitForConnection(db *pgx.Conn) error {
	for i := 0; i < 30; i++ {
		err := db.Ping(context.Background())
		if err == nil {
			return nil
		}
		// wait before trying again
		time.Sleep(time.Second)
	}

	return fmt.Errorf("couldn't connectt to user data basae")
}
