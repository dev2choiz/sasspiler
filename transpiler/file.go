package transpiler

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Fetch all scss files in the directory passed as parameter
func (tr *SassTranspiler) GetFiles(dir string) []afero.File {
	if !tr.isDir(dir) {
		if !tr.IsSassFile(dir) {
			return []afero.File{}
		}
		f, err := tr.fs.Open(dir)
		check(err)
		return []afero.File{f}
	}

	files, err := afero.ReadDir(tr.fs, dir)
	if err != nil {
		log.Fatal(err)
	}

	allFiles := make([]afero.File, 0)
	for _, file := range files {
		f := filepath.Join(dir, file.Name())
		if file.IsDir() {
			allFiles = append(allFiles, tr.GetFiles(f)...)
			continue
		}

		if !tr.IsSassFile(f) {
			continue
		}

		tmp, err := tr.fs.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		allFiles = append(allFiles, tmp)
	}

	return allFiles
}

// Generate css destination pathname
// `srcRoot` is the directory where .scss sources are. (example : /src/scss_root)
// `src` is the absolute path of a file source. (example : /src/scss_root/menu/menu.scss)
// `dest` is either the absolute filename of the destination
// or the root directory of destination
func (tr *SassTranspiler) GetDest(srcRoot, src, dest string) string {
	if !tr.isDir(dest) {
		return dest
	}

	// directory to directory
	rel, err := filepath.Rel(srcRoot, src)
	if err != nil {
		panic(err)
	}

	rel = strings.TrimSuffix(rel, filepath.Ext(rel)) + ".css"

	return filepath.Join(dest, rel)
}

// Test if a file is a .scss or .sass
func (tr *SassTranspiler) IsSassFile(f string) bool {
	ext := strings.ToLower(filepath.Ext(f))
	if ext != ".scss" && ext != ".sass" && ext != ".css" {
		return false
	}
	return true
}

// Test if 'p' is a directory, existing or not
func (tr *SassTranspiler) isDir(p string) bool {
	ok, err := afero.IsDir(tr.fs, p)
	if err == nil {
		return ok
	}

	if !os.IsNotExist(err) {
		log.Fatalln(err)
	}

	// 'p' does not yet exist, then we analyze the string only
	// If the base does not contains a '.' we consider it as a directory
	return !strings.Contains(filepath.Base(p), ".")
}

// Test if 'p' is a file, existing or not
func (tr *SassTranspiler) isFile(p string) bool {
	return !tr.isDir(p)
}
