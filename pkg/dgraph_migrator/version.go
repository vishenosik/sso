package migrate

import (
	"iter"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/vishenosik/sso/pkg/collections"
)

type Version struct {
	Timestamp      string   `json:"version_timestamp,omitempty"`
	CurrentVersion int64    `json:"version_current,omitempty"`
	DType          []string `json:"dgraph.type,omitempty"`
}

type migration struct {
	version  int64
	filename string
}

type migrations = []migration

func collectMigrations(filenames iter.Seq[string]) iter.Seq[migration] {
	return func() iter.Seq[migration] {
		return func(yield func(migration) bool) {
			for filename := range filenames {
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
	}()
}

func sortVersions(mgs migrations) {
	slices.SortFunc(mgs, func(a, b migration) int {
		return int(a.version - b.version)
	})
}

func parseVersion(filename string) (int64, bool) {

	base := filepath.Base(filename)
	if ext := filepath.Ext(base); ext != gqlExt {
		return 0, false
	}

	idx := strings.Index(base, "_")
	if idx < 0 {
		return 0, false
	}

	n, err := strconv.ParseInt(base[:idx], 10, 64)
	if err != nil {
		return 0, false
	}

	if n < 1 {
		return 0, false
	}

	return n, true
}

func filterVersions(current int64, _iter iter.Seq[migration]) iter.Seq[migration] {
	return collections.Filter(_iter, func(m migration) bool {
		return m.version > current
	})
}
