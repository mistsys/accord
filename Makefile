go ?= go
GO_LD_X = -X $(1)=$(2)
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_ID), $(call GO_LD_X,github.com/mistsys/accord.ClientID,'$(GOOGLE_CLIENT_ID)'))
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_SECRET), $(call GO_LD_X,github.com/mistsys/accord.ClientSecret,'$(GOOGLE_CLIENT_SECRET)'))
set_vars += $(if $(GO_LDFLAGS), -ldflags="$(GO_LDFLAGS)")
TOP ?= .
DEPLOYMENT_VAULT ?= depvault

install:
	$(go) install $(set_vars) ./...

_builds:
	mkdir -p _builds/{linux,osx}

build-osx: _builds
	GOOS=darwin $(go) build $(set_vars) -o _builds/osx/accord_client ./cmd/accord_client

build-linux: _builds
	GOOS=linux $(go) build $(set_vars) -o _builds/linux/accord_client ./cmd/accord_client

# clean
clean :
	-$(go) clean -i ./...

test:
	$(go) test ./...

add-deployment:
	$(go) run ./cmd/accord/accord.go -task add-deployment -path.psk (TOP)/terraform/playbooks/files/deployments.json $(DEPLOYMENT_ID)

release-server:
	GOOS=linux go build -o $(TOP)/terraform/playbooks/files/accord $(TOP)/cmd/accord_server
	cd $(TOP)/terraform && make upload

release-client: build-osx build-linux
	aws-vault exec $(DEPLOYMENT_VAULT) -- aws s3 cp  _builds/linux/accord_client $(DEPLOYMENT_S3_URL)/accord_client/accord_client-linux
	aws-vault exec $(DEPLOYMENT_VAULT) -- aws s3 cp  _builds/osx/accord_client $(DEPLOYMENT_S3_URL)/accord_client/accord_client-osx

integration-test:
	./integration.sh

all: install

# fetch depedencies
getdeps :
	$(go) get -u -v -d ./...

.PHONY: all clean test getdeps
