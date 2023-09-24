.PHONY: test clean all

NAME=installer
OS_ARCH ?= linux_amd64_v1

build:
	goreleaser build --snapshot --clean

install: build clean
	mkdir -p /tmp/tfproviders/
	mv dist/terraform-provider-${NAME}_${OS_ARCH}/* /tmp/tfproviders/

test:
	go1.19.12 test $(TESTARGS) -race -parallel=4 ./...

testacc:
	TF_ACC=1 go1.19.12 test $(TESTARGS) -race -parallel=4 ./...

clean:
	rm -rf /tmp/tfproviders/