all: bin/exego_darwin_386 bin/exegod_darwin_386

bin/exego_darwin_386: exego/exego.go
	gox -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./exego

bin/exegod_darwin_386: exegod/exegod.go
	gox -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./exegod
