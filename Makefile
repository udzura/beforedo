VERSION := $(shell cat ./VERSION)
.PHONY: install test setup clean-zip all compress release

beforedo: test
	go build -ldflags "-X main.Version=$(VERSION)"

test:
	go test ./...

setup:
	go get ./...
	which gox || go get github.com/mitchellh/gox
	which ghr || go get github.com/tcnksm/ghr

install: beforedo
	install beforedo /usr/local/bin

clean-zip:
	find pkg -name '*.zip' | xargs rm

all: setup test
	gox \
	    -os="linux,darwin" \
	    -arch="amd64" \
	    -output "pkg/{{.Dir}}_$(VERSION)-{{.OS}}-{{.Arch}}" \
	    $(CMDDIR)

compress: all clean-zip
	cd pkg && ( find . -perm -u+x -type f -name 'beforedo*' | gxargs -i zip -m {}.zip {} )

release: compress
	git push origin master
	ghr $(VERSION) pkg
	git fetch origin --tags

