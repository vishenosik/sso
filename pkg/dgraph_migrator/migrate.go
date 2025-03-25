package migrate

import (
	"context"
	"io/fs"
	"iter"
	"path"

	"github.com/pkg/errors"

	"github.com/dgraph-io/dgo/v210"
)

const (
	gqlExt = ".gql"
)

type dgraphMigrator struct {
	dg   *dgo.Dgraph
	fsys fs.FS
}

func NewDgraphMigrator(
	dg *dgo.Dgraph,
	fsys fs.FS,
) (*dgraphMigrator, error) {

	if dg == nil {
		return nil, errors.New("dgraph client not initialized")
	}

	return &dgraphMigrator{
		dg:   dg,
		fsys: fsys,
	}, nil
}

func (dmr *dgraphMigrator) Up(path string) error {
	return dmr.UpContext(context.Background(), path)
}

// 1. Get current version
// 2. Collect .gql files
// 3. Sort provided versions
// 4. Alter versions greater than current
// 5. If success - update version in dgraph
// 6. If failure - Rollback
func (dmr *dgraphMigrator) UpContext(ctx context.Context, path string) error {

	txn := dmr.dg.NewTxn()
	defer txn.Discard(ctx)

	return nil
}

func (dmr *dgraphMigrator) FetchVersion(ctx context.Context) error {

	txn := dmr.dg.NewTxn()
	defer txn.Discard(ctx)

	return nil
}

func (dmr *dgraphMigrator) UpsertVersion(ctx context.Context) error {

	txn := dmr.dg.NewTxn()
	defer txn.Discard(ctx)

	return nil
}

func (dmr *dgraphMigrator) AddVersionSchema(ctx context.Context) error {

	txn := dmr.dg.NewTxn()
	defer txn.Discard(ctx)

	_ = `
		index_name: string @index(exact) .
		version_timestamp: datetime .
		current_version: int .
		type SchemaVersion {
			index_name: string
			version_timestamp: datetime
			current_version: int
		}`
	return nil
}

func collectMigrations(fsys fs.FS, dirpath string) (iter.Seq[migration], error) {

	if _, err := fs.Stat(fsys, dirpath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.Wrap(err, "migrations directory not found")
		}

		return nil, errors.Wrap(err, "migrations directory unexpected "+dirpath)
	}

	filenames, err := fs.Glob(fsys, path.Join(dirpath, "*"+gqlExt))
	if err != nil {
		return nil, errors.Wrap(err, "migrations not found")
	}

	_iter := func() iter.Seq[migration] {
		return func(yield func(migration) bool) {
			for _, filename := range filenames {
				version, ok := parseVersion(filename)
				if !ok {
					continue
				}

				_migration := migration{
					version:  version,
					filename: filename,
				}

				if !yield(_migration) {
					return
				}

			}

		}
	}
	return _iter(), nil
}
