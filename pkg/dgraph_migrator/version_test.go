package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortVersions(t *testing.T) {
	t.Parallel()

	mgs := migrations{
		{version: 1, filename: "0001_create_initial_schema.gql"},
		{version: 3, filename: "0003_add_post_table.gql"},
		{version: 2, filename: "0002_add_user_table.gql"},
	}

	sortVersions(mgs)

	assert.Equal(t, int64(1), mgs[0].version)
	assert.Equal(t, int64(2), mgs[1].version)
	assert.Equal(t, int64(3), mgs[2].version)
}

func Test_parseVersion(t *testing.T) {

	testingTable := []struct {
		name            string
		filename        string
		expectedVersion int64
		expectedOk      bool
	}{
		{
			name:            "parse_version_success_1",
			filename:        "0002_add_user_table.gql",
			expectedVersion: 2,
			expectedOk:      true,
		},
		{
			name:            "parse_version_success_2",
			filename:        "0001_create_initial_schema.gql",
			expectedVersion: 1,
			expectedOk:      true,
		},
		{
			name:            "parse_version_error_1",
			filename:        "0000_invalid_version.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_filename_1",
			filename:        "invalid_filename.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_filename_2",
			filename:        "invalidFilename.gql",
			expectedVersion: 0,
			expectedOk:      false,
		},
		{
			name:            "parse_version_error_invalid_extension",
			filename:        "0002_add_user_table.sql",
			expectedVersion: 0,
			expectedOk:      false,
		},
	}

	for _, tt := range testingTable {

		t.Run(tt.name, func(t *testing.T) {
			actualVersion, ok := parseVersion(tt.filename)
			assert.Equal(t, ok, tt.expectedOk)
			assert.Equal(t, actualVersion, tt.expectedVersion)
		})

	}

}
