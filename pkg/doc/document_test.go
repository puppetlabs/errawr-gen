package doc_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/errawr-gen/pkg/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	filenames, err := filepath.Glob("testdata/docs/*.yml")
	require.NoError(t, err)

	for _, filename := range filenames {
		t.Run(filepath.Base(filename), func(t *testing.T) {
			b, err := ioutil.ReadFile(filename)
			require.NoError(t, err)

			_, err = doc.New(string(b))
			assert.NoError(t, err)
		})
	}
}
