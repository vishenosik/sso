package migrate

import (
	"context"
	"io/fs"

	"github.com/pkg/errors"

	"github.com/dgraph-io/dgo/v210"
)

const (
	gqlExt = ".gql"
)

var (
	ErrVersionFetch = errors.New("no version fetched")
)

type dgraphMigrator struct {
	dg             *dgo.Dgraph
	fsys           fs.FS
	currentVersion int64
}

func NewDgraphMigrator(client *dgo.Dgraph, fsys fs.FS) (*dgraphMigrator, error) {
	return NewDgraphMigratorContext(context.Background(), client, fsys)
}

func NewDgraphMigratorContext(
	ctx context.Context,
	client *dgo.Dgraph,
	fsys fs.FS,
) (*dgraphMigrator, error) {

	if client == nil {
		return nil, errors.New("dgraph client not initialized")
	}

	if err := applySchema(client, ctx); err != nil {
		return nil, err
	}

	version, err := fetchVersion(client, ctx)
	if err != nil && !errors.Is(err, ErrVersionFetch) {
		return nil, err
	}

	return &dgraphMigrator{
		dg:             client,
		fsys:           fsys,
		currentVersion: version.CurrentVersion,
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
