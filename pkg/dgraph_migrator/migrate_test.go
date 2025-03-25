package migrate

import (
	"embed"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test/migrations
	test_migrations embed.FS
)

func Test_collectMigrations(t *testing.T) {

	testdir := path.Join("test", "migrations")

	mgs := migrations{
		{version: 1, filename: path.Join(testdir, "0001_create_initial_schema.gql")},
		{version: 2, filename: path.Join(testdir, "0002_add_user_table.gql")},
		{version: 3, filename: path.Join(testdir, "0003_add_post_table.gql")},
	}

	actual, err := collectMigrations(test_migrations, testdir)
	assert.NoError(t, err)
	assert.Equal(t, mgs, actual)

}
