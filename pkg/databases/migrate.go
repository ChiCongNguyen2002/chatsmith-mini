package databases

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"chatsmith-mini/pkg/failure"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool, scriptsDir string) error {
	files, err := filepath.Glob(filepath.Join(scriptsDir, "*.sql"))
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	sort.Strings(files)

	for _, file := range files {
		if err := runScript(ctx, pool, file); err != nil {
			return err
		}
	}

	return nil
}

func runScript(ctx context.Context, pool *pgxpool.Pool, file string) error {
	script, err := os.ReadFile(file)
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(script)); err != nil {
		return failure.ErrWithTrace(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}
