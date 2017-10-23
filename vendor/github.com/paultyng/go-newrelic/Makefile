NOVENDOR = $(shell glide novendor)

all: fix vet lint test

fix:
	go fix $(NOVENDOR)

vet:
	go vet $(NOVENDOR)

test:
	go test -v -cover $(NOVENDOR)

lint:
	printf "%s\n" "$(NOVENDOR)" | xargs -I {} sh -c 'golint -set_exit_status {}'
