package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func addUser(conn *pgx.Conn, email string, plainPassword string, isStaff bool, firstName string, lastName string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	password := string(bytes)
	var id int
	query := `insert into "user" (email, password, is_staff) values ($1, $2, $3) returning id;`
	err = conn.QueryRow(context.Background(), query, email, password, isStaff).Scan(&id)
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok {
			if pgErr.Code == "23505" {
				fmt.Printf("User %s already exists\n", email)
				return
			}
		}
		panic(err)
	}
	fmt.Printf("Added user: %s\n", email)

	query = `insert into settings(user_id, first_name, last_name) values ($1, $2, $3);`
	_, err = conn.Exec(context.Background(), query, id, firstName, lastName)
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok {
			if pgErr.Code == "23505" {
				fmt.Printf("Settings already exists\n")
				return
			}
		}
		panic(err)
	}
	fmt.Printf("Added settings\n")
}

func addServiceUser(conn *pgx.Conn, login string, password string, role string) {
	login = pgx.Identifier{login}.Sanitize()
	password = strings.ReplaceAll(password, "'", "''")
	query := fmt.Sprintf(`create user %s with password '%s' login;`, login, password)
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok {
			if pgErr.Code == "42710" {
				fmt.Printf("Service user %s already exists\n", login)
				return
			}
		}
		panic(err)
	}
	query = fmt.Sprintf(`grant %s to %s;`, role, login)
	_, err = conn.Exec(context.Background(), query)
	panicOnError(err)
	fmt.Printf("Added service user %s\n", login)
}

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	name := os.Getenv("POSTGRES_DB")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, host, port, name)

	conn, err := pgx.Connect(context.Background(), dbUrl)
	panicOnError(err)
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			fmt.Printf("Unable to close connection: %v\n", err)
		}
	}()

	systemUserPassword := os.Getenv("SYSTEM_USER_PASSWORD")
	if systemUserPassword == "" {
		panic("SYSTEM_USER_PASSWORD is not set")
	}
	addUser(conn, "system@example.com", systemUserPassword, true, "System", "")
	adminUserPassword := os.Getenv("ADMIN_USER_PASSWORD")
	if adminUserPassword == "" {
		panic("ADMIN_USER_PASSWORD is not set")
	}
	addUser(conn, "admin@example.com", adminUserPassword, true, "Admin", "admin")
	testUserPassword := os.Getenv("TEST_USER_PASSWORD")
	if testUserPassword == "" {
		panic("TEST_USER_PASSWORD is not set")
	}
	addUser(conn, "test@example.com", testUserPassword, false, "John", "Doe")
	test2UserPassword := os.Getenv("TEST2_USER_PASSWORD")
	if test2UserPassword == "" {
		panic("TEST2_USER_PASSWORD is not set")
	}
	addUser(conn, "test2@example.com", test2UserPassword, false, "Ivan", "Ivanov")
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
