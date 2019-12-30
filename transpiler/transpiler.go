package transpiler

import (
	"github.com/spf13/afero"
	"github.com/wellington/go-libsass"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Transpiler interface {
	Verbose() bool
	SetVerbose(bool)
	FileSystem() afero.Fs
	SetFileSystem(afero.Fs)
	Run(source, dest string, files []afero.File, imports []string) error
	GetFiles(dir string) []afero.File
}

type SassTranspiler struct {
	verbose bool
	fs      afero.Fs
}

func (tr *SassTranspiler) Verbose() bool {
	return tr.verbose
}

func (tr *SassTranspiler) SetVerbose(verbose bool) {
	tr.verbose = verbose
}

func (tr *SassTranspiler) FileSystem() afero.Fs {
	return tr.fs
}

func (tr *SassTranspiler) SetFileSystem(fileSystem afero.Fs) {
	tr.fs = fileSystem
}

// Get a new instance of the SassTranspiler
func New() *SassTranspiler {
	return &SassTranspiler{}
}

// Transpile scss files in `files` to `dest` directory
func (tr *SassTranspiler) Run(source, dest string, files []afero.File, imports []string) error {
	chErr := make(chan error)
	chDone := make(chan struct{})
	nbTr := len(files)
	for _, file := range files {
		go func(file afero.File) {
			if err := tr.transpile(source, dest, file, imports); err != nil {
				chErr <- err
				return
			}
			chDone <- struct{}{}
		}(file)
	}

	// Wait for the end of all transpilations or for one error
	count := 0
	for {
		select {
		case err := <-chErr:
			return err
		case <-chDone:
			count++
			if count == nbTr {
				return nil
			}
		}
	}
}

// Transpile a unique scss file to `dest` directory
func (tr *SassTranspiler) transpile(source, dest string, file afero.File, imports []string) error {
	defer func() {
		check(file.Close())
	}()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		if tr.verbose {
			log.Println(file.Name(), "is empty")
		}
		return nil
	}

	root := source
	if tr.isFile(source) {
		root = filepath.Dir(file.Name())
	}

	d := tr.GetDest(root, file.Name(), dest)
	_ = tr.fs.MkdirAll(filepath.Dir(d), os.FileMode(0755))
	w, err := tr.fs.Create(d)
	if err != nil {
		return err
	}

	comp, err := NewCompiler(w, file)
	if err != nil {
		return err
	}

	if err := AddImportPath(comp, imports); err != nil {
		return err
	}

	if err := Compile(comp); err != nil {
		return err
	}

	if tr.verbose {
		log.Println(file.Name(), "==>", w.Name())
	}

	return nil
}

// In a variable to make it testable
var NewCompiler = func(dst io.Writer, src io.Reader) (libsass.Compiler, error) {
	return libsass.New(dst, src)
}

// In a variable to make it testable
var AddImportPath = func(c libsass.Compiler, i []string) error {
	return c.Option(libsass.IncludePaths(i))
}

// In a variable to make it testable
var Compile = func(c libsass.Compiler) error {
	return c.Run()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
