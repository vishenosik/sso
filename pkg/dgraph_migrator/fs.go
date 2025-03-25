package migrate

import (
	"io/fs"
	"iter"
	"path"

	"github.com/pkg/errors"
)

func collectFilenames(fsys fs.FS, dirpath string) (iter.Seq[string], error) {
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

	_iter := func() iter.Seq[string] {
		return func(yield func(string) bool) {
			for _, filename := range filenames {
				if !yield(filename) {
					return
				}
			}
		}
	}

	return _iter(), nil
}
