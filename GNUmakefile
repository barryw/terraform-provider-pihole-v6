default: testacc

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build -o terraform-provider-pihole-v6

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/barryw/pihole-v6/0.0.1/linux_amd64
	cp terraform-provider-pihole-v6 ~/.terraform.d/plugins/registry.terraform.io/barryw/pihole-v6/0.0.1/linux_amd64/

generate:
	go generate ./...

docs:
	tfplugindocs generate --provider-name pihole
