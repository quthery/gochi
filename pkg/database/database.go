package main

import (
	"context"
	"fmt"
	"gochi/pkg/structs"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

//storage, err := ConnectDB("postgres://postgres:root@localhost:5432/gochi")

type db struct {
	db   *pgxpool.Pool
	conn *pgxpool.Conn
}

func ConnectDB(dbURL string) (*db, error) {
	ctx := context.Background()

	pgxPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось создать пул подключений: %v\n", err)
		return nil, err
	}
	fmt.Println("Подключение к базе данных успешно!")

	conn, err := pgxPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Unable to acquire a database connection: " + err.Error())
	}
	defer conn.Release()

	return &db{db: pgxPool, conn: conn}, nil
}

func (db *db) createTable() error {
	const tableQuery = `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(50) NOT NULL,
				email VARCHAR(100) NOT NULL
			)
		`
	_, err := db.db.Exec(context.Background(), tableQuery)
	if err != nil {
		return err
	}
	fmt.Println("Таблица 'users' успешно создана.")
	return nil
}

func (db *db) DropTable() error {
	const DropTableQuery = `DROP TABLE users`
	_, err := db.db.Exec(context.Background(), DropTableQuery)
	if err != nil {
		log.Fatal("database error drop table error", err)
		return err
	}
	return nil
}

func (db *db) GetAllColumns() ([]structs.User, error) {
	// Define your query here
	const query = "SELECT * FROM users"

	rows, err := db.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var user structs.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *db) InsertUser() {
	_, err := db.db.Exec(context.Background(), "INSERT INTO users (username, email) VALUES ($1, $2)",
		"john_doe", "john@example.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при выполнении запроса: %v\n", err)
		return
	}
	fmt.Print("succes")
}

func main() {
	//todo сделать все через пул в структуре т.к. все работает через пул и конекшин
	//делать не обязательно
	storage, err := ConnectDB("postgres://postgres:root@localhost:5432/gochi")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к базе данных: %v\n", err)
		return
	}

	err = storage.createTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при создании таблицы: %v\n", err)
		return
	}

	storage.InsertUser()
	stack, _ := storage.GetAllColumns()
	fmt.Println(stack)
}
