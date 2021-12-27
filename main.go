package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	// Pass as a command line argument:
	//user:password@tcp(127.0.0.1:3306)/db_name
	runMigration(os.Args[1])
}

func runMigration(datasource string) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	migrate(tx,
		`CREATE TABLE A (
		id char(36) NOT NULL,
		PRIMARY KEY(id),
		name VARCHAR (32) NOT NULL
	);`,
	)

	migrate(tx,
		`CREATE TABLE B (
		id char(36) NOT NULL,
		PRIMARY KEY(id),
		A_id char(36),
		FOREIGN KEY (A_id) REFERENCES A(id)
	);`,
	)

	//migrate(tx,
	//	`SET FOREIGN_KEY_CHECKS = 0;`,
	//)
	migrate(tx,
		`ALTER TABLE B MODIFY A_id char(36) NOT NULL;`,
	)
	//migrate(tx,
	//	`SET FOREIGN_KEY_CHECKS = 1;`,
	//)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func migrate(tx *sql.Tx, statement string) {
	_, err := tx.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
}
