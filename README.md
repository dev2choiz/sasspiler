# Sasspiler
![](https://img.shields.io/github/tag/dev2choiz/sasspiler.svg) ![](https://img.shields.io/github/release/dev2choiz/sasspiler.svg) ![](https://img.shields.io/github/issues/dev2choiz/sasspiler.svg)

Sasspiler is written in golang to transpile easily `.scss` files to `.css`.  
It uses [wellington/go-libsass](https://github.com/wellington/go-libsass).  

## Installation with `go get`
This method require golang.
```sh
go get github.com/dev2choiz/sasspiler
```
The command above will generate `sasspiler` binary in `$GOPATH\bin`

## Installation with `docker-compose`

```sh
git clone git@github.com:dev2choiz/sasspiler.git
cd sasspiler
docker-compose up --build
```
The `sasspiler` binary will be generated in `./bin`
Then :

```sh
cp ./bin/sasspiler $GOPATH/bin/
```

## Usage
```sh
sasspiler --source=/absolute/scss/dir --dest=/absolute/css/dir --importDir=/absolute/scss/dir1,/absolute/scss/dir2 --verbose
```
or  
```sh
sasspiler -s=/absolute/scss/dir -d=/absolute/css/dir -i=/absolute/scss/dir1,/absolute/scss/dir2 -v
```
`--importDir` argument can be ignored if there is not `@import` in your scss to transpile.


You can also transpile a unique file, in order for example to configure the file watcher of your IDE.
```sh
sasspiler --source=/path/to/file.scss --dest=/path/to/destination/file.css --importDir=/absolute/scss/dir1,/absolute/scss/dir2 --verbose
```