package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
	"github.com/X-AROK/urlcut/migrations"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
	"time"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(db *sql.DB) (*DBStore, error) {
	dbs := &DBStore{db: db}
	if err := dbs.migrate(); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	return dbs, nil
}

func (dbs *DBStore) migrate() error {
	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect error: %w", err)
	}

	if err := goose.Up(dbs.db, "."); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}

	return nil
}

func (dbs *DBStore) Add(ctx context.Context, u *url.URL) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	id := util.GenerateID(8)

	_, err := dbs.db.ExecContext(ctxTimeout, "INSERT INTO urls (short, original) VALUES ($1, $2)", id, u.OriginalURL)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		id, err := dbs.GetByOriginal(ctx, u)
		if err != nil {
			return "", fmt.Errorf("get by original url error: %w", err)
		}
		return "", url.NewAlreadyExistsError(id)
	}
	if err != nil {
		return "", fmt.Errorf("insert query error: %w", err)
	}
	u.ShortURL = id
	return id, nil
}

func (dbs *DBStore) Get(ctx context.Context, id string) (*url.URL, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(ctxTimeout, "SELECT short, original FROM urls WHERE short=$1", id)
	if row == nil {
		return nil, url.ErrNotFound
	}

	u := &url.URL{}
	err := row.Scan(&u.ShortURL, &u.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("row scan error: %w", err)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return u, nil
}

func (dbs *DBStore) GetByOriginal(ctx context.Context, u *url.URL) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(ctxTimeout, "SELECT short FROM urls WHERE original=$1", u.OriginalURL)
	if row == nil {
		return "", url.ErrNotFound
	}

	err := row.Scan(&u.ShortURL)
	if err != nil {
		return "", fmt.Errorf("row scan error: %w", err)
	}

	err = row.Err()
	if err != nil {
		return "", fmt.Errorf("row error: %w", err)
	}

	return u.ShortURL, nil
}

func (dbs *DBStore) AddBatch(ctx context.Context, urls *url.URLsBatch) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (short, original) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("stmt prepare error: %w", err)
	}

	for _, v := range *urls {
		id := util.GenerateID(8)
		_, err := stmt.ExecContext(ctx, id, v.OriginalURL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert query error: %w", err)
		}
		v.ShortURL = id
	}
	tx.Commit()

	return nil
}
