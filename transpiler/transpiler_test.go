package transpiler_test

import (
	"bytes"
	"errors"
	"github.com/dev2choiz/sasspiler/transpiler"
	tmocks "github.com/dev2choiz/sasspiler/transpiler/mocks"
	"github.com/spf13/afero"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"github.com/wellington/go-libsass"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Mocks generations
//go:generate mockery -dir . -name Transpiler
//go:generate mockery -dir ./../vendor/github.com/spf13/afero -name File

var (
	// Fake file system
	appFs = new(afero.MemMapFs)

	// Where sources are stocked in this tests
	src = "/src/scss"

	// Destination of generated css
	dest = "/public/css"

	// Directories where are scss imported with '@import'
	importDirs = make([]string, 0)
)

func TestMain(m *testing.M) {
	//setup before to run
	generateFixtures(appFs, src, data)

	os.Exit(m.Run())
}

func TestRun(t *testing.T) {
	tr := getTranspiler(true)

	err := tr.Run(src, dest, tr.GetFiles(src), importDirs)
	check(err, t)

	// check generated css
	for _, datum := range data {
		tarRelPath := strings.TrimSuffix(datum["relPath"], filepath.Ext(datum["relPath"]))
		tarPath := filepath.Join(dest, tarRelPath+".css")

		if datum["content"] == "" || datum["relPath"] == "not_scss.txt" {
			exist, err := afero.Exists(tr.FileSystem(), tarPath)
			check(err, t)
			assert.False(t, exist, "%s should not exist", tarPath)
			continue
		}

		f, err := appFs.OpenFile(tarPath, os.O_RDONLY, os.FileMode(0755))
		assert.NoError(t, err)

		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, f)
		check(err, t)
		check(f.Close(), t)

		assert.Equal(t, datum["expected"], string(buf.Bytes()))
	}
}

func TestRunWithFileError(t *testing.T) {
	msgErr := "error when Stat()"
	tr := getTranspiler(false)

	f := &tmocks.File{}
	f.On("Stat").Return(nil, errors.New(msgErr)).
		On("Close").Return(nil)

	err := tr.Run("dummy", "dummy", []afero.File{f}, importDirs)

	assert.NotNil(t, err)
	assert.EqualError(t, err, msgErr)
}

func TestRunWithNewCompilerError(t *testing.T) {
	msgErr := "cannot create compiler"
	tr := getTranspiler(false)

	file, err := tr.FileSystem().Open(filepath.Join(src, "style.scss"))
	check(err, t)

	sav := transpiler.NewCompiler
	defer func() {
		transpiler.NewCompiler = sav
	}()
	transpiler.NewCompiler = func(dst io.Writer, src io.Reader) (compiler libsass.Compiler, err error) {
		return nil, errors.New(msgErr)
	}

	err = tr.Run(src, dest, []afero.File{file}, importDirs)
	assert.EqualError(t, err, msgErr)
}
func TestRunWithAddImportPathError(t *testing.T) {
	msgErr := "cannot add option"
	tr := getTranspiler(false)

	file, err := tr.FileSystem().Open(filepath.Join(src, "style.scss"))
	check(err, t)

	sav := transpiler.AddImportPath
	defer func() {
		transpiler.AddImportPath = sav
	}()
	transpiler.AddImportPath = func(c libsass.Compiler, i []string) error {
		return errors.New(msgErr)
	}

	err = tr.Run(src, dest, []afero.File{file}, importDirs)
	assert.EqualError(t, err, msgErr)
}

func TestRunWithErrorWhileCompile(t *testing.T) {
	msgErr := "error while compile"
	tr := getTranspiler(false)

	file, err := tr.FileSystem().Open(filepath.Join(src, "style.scss"))
	check(err, t)

	sav := transpiler.Compile
	defer func() {
		transpiler.Compile = sav
	}()
	transpiler.Compile = func(comp libsass.Compiler) error {
		return errors.New(msgErr)
	}

	err = tr.Run(src, dest, []afero.File{file}, importDirs)
	assert.EqualError(t, err, msgErr)
}

func generateFixtures(fs *afero.MemMapFs, dir string, data []map[string]string) {
	for _, d := range data {
		p := filepath.Join(dir, d["relPath"])
		err := fs.MkdirAll(filepath.Dir(p), os.FileMode(0755))
		check(err, nil)
		f, err := fs.Create(p)
		check(err, nil)
		_, err = f.WriteString(d["content"])
		check(err, nil)
		err = f.Close()
		check(err, nil)
	}
}

func check(e error, t *testing.T) {
	if e != nil {
		if t != nil {
			t.Error(e)
		} else {
			panic(e)
		}
	}
}

func getTranspiler(v bool) *transpiler.SassTranspiler {
	t := transpiler.New()
	t.SetFileSystem(appFs)
	t.SetVerbose(v)
	return t
}

func fixPart(relPath, content, expected string) map[string]string {
	return map[string]string{
		"relPath":  relPath,
		"content":  content,
		"expected": expected,
	}
}

// Fixtures
var data = []map[string]string{
	fixPart("style.scss", `body {color: #fff; h1 {font-size: 17px;}}`, `body {
  color: #fff; }
  body h1 {
    font-size: 17px; }
`),
	fixPart("partials/menu.scss", `ul.menu {diplay: flex; li {color: red;}}`, `ul.menu {
  diplay: flex; }
  ul.menu li {
    color: red; }
`),
	fixPart("empty.scss", "", ""),
	fixPart("not_scss.txt", "dummy", ""),
}
