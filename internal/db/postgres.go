package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuromii5/messagio/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func PGConnectionStr(config config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)
}

func NewDB(config config.PostgresConfig) *DB {
	dbUrl := PGConnectionStr(config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("unable to parse db url")
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("unable to connect to db")
	}

	return &DB{Pool: pool}
}

func (d *DB) Save(ctx context.Context, message string) (int32, error) {
	query := `INSERT INTO messages (content, processed, created_at) VALUES ($1, $2, $3) RETURNING id;`

	var id int32
	err := d.Pool.QueryRow(ctx, query, message, false, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DB) MarkAsProcessed(ctx context.Context, messageID int32) error {
	query := `UPDATE messages SET processed = true WHERE id = $1;`

	if _, err := d.Pool.Exec(ctx, query, messageID); err != nil {
		return err
	}

	return nil
}

func (d *DB) LoadStats(ctx context.Context) (int32, error) {
	query := `SELECT COUNT(1) AS processed_count FROM messages WHERE processed = true;`

	var processedCount int32
	err := d.Pool.QueryRow(ctx, query).Scan(&processedCount)
	if err != nil {
		return 0, err
	}

	return processedCount, nil
}
