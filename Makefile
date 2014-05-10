#go get github.com/jteeuwen/go-bindata/...
#go get github.com/mitchellh/gox

all: bin/exego_darwin_386 bin/exegod_darwin_386

bin/exego_darwin_386: exego/exego.go *.go
	gox -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./exego

bin/exegod_darwin_386: exegod/exegod.go *.go
	gox -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./exegod

certs.go: certs/myCA.cer
	go-bindata -pkg exego -o "./certs.go" certs
