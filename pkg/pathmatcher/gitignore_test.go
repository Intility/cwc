package pathmatcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempLoc() string {
	return filepath.Join(os.TempDir(), "cwc", "pathmatcher", "gitignore_test")
}

func createFiles(t *testing.T, names []string) {
	for _, name := range names {
		f, err := os.Create(filepath.Join(tempLoc(), name))
		require.NoError(t, err, "Error creating temp file %s", err)
		f.Close()
	}
}

func deleteFiles(t *testing.T, names []string) {
	for _, name := range names {
		err := os.Remove(filepath.Join(tempLoc(), name))
		require.NoError(t, err, "Error removing temp file %s", err)
	}
}

func TestMatch(t *testing.T) {
	v := tempLoc()
	require.NoError(t, os.MkdirAll(v, 0777))
	defer func() {
		_ = os.RemoveAll(filepath.Join(os.TempDir(), "cwc"))
	}()

	// match with multiple casing
	assert.True(t, match("[Bb]uild", "Build"))
	assert.True(t, match("[Bb]uild", "build"))

	// match with dir
	assert.True(t, match("[Bb]in/", "bin/"))
	assert.True(t, match("[Bb]in/", "bin/foo.txt"))
	assert.False(t, match("[Bb]in/", "Fin/foo.txt"))

	// blank line
	assert.False(t, match("", v))

	// a comment
	assert.False(t, match("#a comment", v))

	// regular match no slash
	assert.True(t, match("gitglob.go", "gitglob.go"))

	// negation no slash
	assert.False(t, match("!gitglob.go", "gitglob.go"))

	// match with slash
	tmpFiles := []string{"foo.txt"}
	createFiles(t, tmpFiles)
	assert.True(t, match(tempLoc()+"/foo.txt", v+"/foo.txt"))
	deleteFiles(t, tmpFiles)

	// negate match with slash
	tmpFiles = []string{"foo.txt"}
	createFiles(t, tmpFiles)
	assert.False(t, match("!"+tempLoc()+"/foo.txt", v+"/foo.txt"))
	deleteFiles(t, tmpFiles)

	// directory
	assert.True(t, match(tempLoc(), v))

	// directory with trailing slash
	//assert.True(t, match(tempLoc()+"/", v)) // this is wrong, no? pattern `foo/` should not match `foo`

	// star matching
	tmpFiles = []string{"foo.txt"}
	createFiles(t, tmpFiles)
	assert.True(t, match(tempLoc()+"/*.txt", v+"/foo.txt"))
	assert.False(t, match(tempLoc()+"/*.txt", v+"/somedir/foo.txt"))
	deleteFiles(t, tmpFiles)

	// double star prefix
	assert.True(t, match("**/foo.txt", v+"/hello/foo.txt"))
	assert.True(t, match("**/foo.txt", v+"/some/dirs/foo.txt"))

	// double star suffix
	assert.True(t, match(tempLoc()+"/hello/**", v+"/hello/foo.txt"))
	assert.False(t, match(tempLoc()+"/hello/**", v+"/some/dirs/foo.txt"))

	// double star in path
	assert.True(t, match(tempLoc()+"/hello/**/world.txt", v+"/hello/world.txt"))
	assert.True(t, match(tempLoc()+"/hello/**/world.txt", v+"/hello/stuff/world.txt"))
	assert.False(t, match(tempLoc()+"/hello/**/world.txt", v+"/some/dirs/foo.txt"))

	// negate doubl start patterns
	assert.False(t, match("!**/foo.txt", v+"/hello/foo.txt"))
	assert.False(t, match("!"+tempLoc()+"/hello/**", v+"/hello/foo.txt"))
	assert.False(t, match("!"+tempLoc()+"/hello/**/world.txt", v+"/hello/world.txt"))
}
