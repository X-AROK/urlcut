package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(db *sql.DB) (*DBStore, error) {
	dbs := &DBStore{db: db}
	err := dbs.createTables()
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (dbs *DBStore) createTables() error {
	_, err := dbs.db.Exec("CREATE TABLE IF NOT EXISTS urls (short VARCHAR(8) PRIMARY KEY, original TEXT UNIQUE)")
	return err
}

func (dbs *DBStore) Add(ctx context.Context, u *url.URL) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	id := util.GenerateID(8)

	_, err := dbs.db.ExecContext(ctxTimeout, "INSERT INTO urls (short, original) VALUES ($1, $2)", id, u.OriginalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			id, err := dbs.GetByOriginal(ctx, u)
			if err != nil {
				return "", err
			}
			return id, url.ErrorAlreadyExists
		} else {
			return "", err
		}
	}
	u.ShortURL = id
	return id, nil
}

func (dbs *DBStore) Get(ctx context.Context, id string) (*url.URL, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(ctxTimeout, "SELECT short, original FROM urls WHERE short=$1", id)
	if row == nil {
		return nil, url.ErrorNotFound
	}

	u := &url.URL{}
	err := row.Scan(&u.ShortURL, &u.OriginalURL)
	if err != nil {
		return nil, err
	}

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (dbs *DBStore) GetByOriginal(ctx context.Context, u *url.URL) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := dbs.db.QueryRowContext(ctxTimeout, "SELECT short FROM urls WHERE original=$1", u.OriginalURL)
	if row == nil {
		return "", url.ErrorNotFound
	}

	err := row.Scan(&u.ShortURL)
	if err != nil {
		return "", err
	}

	err = row.Err()
	if err != nil {
		return "", err
	}

	return u.ShortURL, nil
}

func (dbs *DBStore) AddBatch(ctx context.Context, urls *url.URLsBatch) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO urls (short, original) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	for _, v := range *urls {
		id := util.GenerateID(8)
		_, err := stmt.ExecContext(ctx, id, v.OriginalURL)
		if err != nil {
			tx.Rollback()
			return err
		}
		v.ShortURL = id
	}
	tx.Commit()

	return nil
}
