package transpiler_test

import (
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetFilesWithADir(t *testing.T) {
	tr := getTranspiler(true)

	files := tr.GetFiles(src)

	assert.Equal(t, 3, len(files))
	assert.Equal(t, "/src/scss/empty.scss", files[0].Name())
	assert.Equal(t, "/src/scss/partials/menu.scss", files[1].Name())
	assert.Equal(t, "/src/scss/style.scss", files[2].Name())
	// /src/scss/not_scss.txt is ignored as expected.
}

func TestGetFilesWithAFile(t *testing.T) {
	tr := getTranspiler(true)

	files := tr.GetFiles(filepath.Join(src, "style.scss"))

	assert.Equal(t, 1, len(files))
	assert.Equal(t, "/src/scss/style.scss", files[0].Name())
}

func TestGetFilesWithANoScssFile(t *testing.T) {
	tr := getTranspiler(true)

	files := tr.GetFiles(filepath.Join(src, "not_scss.txt"))

	assert.Equal(t, 0, len(files))
}

func TestGetDest(t *testing.T) {
	tr := getTranspiler(false)

	srcRoot := "/src"
	for _, d := range []struct {
		src  string
		dst  string
		want string
	}{
		// Existing file
		{"/src/style.css", "/dst", "/dst/style.css"},
		// Existing file in a sub folder of the 'srcRoot'
		{"/src/partials/menu.css", "/dst", "/dst/partials/menu.css"},

		// No Existing file in a sub folder of the 'srcRoot'
		{"/src/partials/no_exist.css", "/dst", "/dst/partials/no_exist.css"},
		// No Existing file in a sub sub folder of the 'srcRoot'
		{"/src/partials/partials/no_exist.css", "/dst", "/dst/partials/partials/no_exist.css"},

		// Existing file with 'dst' being a full filepath
		{"/src/style.css", "/dst/custom_name.css", "/dst/custom_name.css"},
		// No existing file with 'dst' being a full filepath
		{"/src/no_exist.css", "/dst/custom_name.css", "/dst/custom_name.css"},
	} {
		got := tr.GetDest(srcRoot, d.src, d.dst)
		assert.Equal(t, d.want, got)
	}
}

func TestIsSassFile(t *testing.T) {
	tr := getTranspiler(false)

	for _, d := range []struct {
		filepath string
		want     bool
	}{
		{"/src/file.scss", true},
		{"/src/file.SCSS", true},
		{"/src/file.sass", true},
		{"/src/file.SASS", true},
		{"/src/file.css", true},
		{"/src/file.CSS", true},
		{"/src/file.js", false},
		{"/src/file", false},
	} {
		got := tr.IsSassFile(d.filepath)
		assert.Equal(t, d.want, got)
	}
}
