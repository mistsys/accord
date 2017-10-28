go ?= go
GO_LD_X = -X $(1)=$(2)
# shrink the binaries
GO_LDFLAGS += -s -w
GO_LDFLAGS += $(if $(DEFAULT_SERVER), $(call GO_LD_X,github.com/mistsys/accord.DefaultServer,'$(DEFAULT_SERVER)'))
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_ID), $(call GO_LD_X,github.com/mistsys/accord.ClientID,'$(GOOGLE_CLIENT_ID)'))
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_SECRET), $(call GO_LD_X,github.com/mistsys/accord.ClientSecret,'$(GOOGLE_CLIENT_SECRET)'))
set_vars += $(if $(GO_LDFLAGS), -ldflags="$(GO_LDFLAGS)")
TOP ?= .
DEPLOYMENT_VAULT ?= depvault
COVERAGE_DIR = _coverage
coveragesubdirs =$(sort $(shell find . \( -path ./git -o -path ./vendor \) -prune -o -name '*.go' -a -printf '%h\n' | sed 's/^.\///g' | xargs))
RACE ?= "-race"

install:
	$(go) install $(set_vars) ./...

_builds:
	mkdir -p _builds/{linux,osx}

build-osx: _builds
	GOOS=darwin $(go) build $(set_vars) -o _builds/osx/accord_client ./cmd/accord_client
	upx --brute _builds/osx/accord_client

build-linux: _builds
	GOOS=linux $(go) build $(set_vars) -o _builds/linux/accord_client ./cmd/accord_client
	upx --brute _builds/linux/accord_client

# clean
clean :
	$(go) clean -i ./...
	rm $(COVERAGE_DIR)/* || true


test:
	$(go) test ./...

dev-server-insecure:
	$(go) run $(TOP)/cmd/accord_server/server.go -rootca $(TOP)/test_assets/root_ca_20170927 -rootcapassword="staple horse apple newton" -userca $(TOP)/test_assets/user_ca_20170927 -usercapassword "staple horse apple thatcher" -insecure

dev-client-insecure-hostauth:
	$(go) run $(TOP)/cmd/accord_client/client.go -task=hostcert -insecure -deploymentId=test -psk=JpUtbRukLuIFyjeKpA4fIpjgs6MTV8eH -hostkeys=$(TOP)/test_assets/test_pubkeys/

add-deployment:
	$(go) run $(TOP)/cmd/accord/accord.go -task add-deployment -path.psk $(TOP)/terraform/playbooks/files/deployments.json $(DEPLOYMENT_ID)

release-server:
	GOOS=linux go build $(set_vars) -o $(TOP)/terraform/playbooks/files/accord $(TOP)/cmd/accord_server
	cd $(TOP)/terraform && make upload

release-client: build-osx build-linux
	aws-vault exec $(DEPLOYMENT_VAULT) -- aws s3 cp  _builds/linux/accord_client $(DEPLOYMENT_S3_URL)/accord_client/accord_client-linux
	aws-vault exec $(DEPLOYMENT_VAULT) -- aws s3 cp  _builds/osx/accord_client $(DEPLOYMENT_S3_URL)/accord_client/accord_client-osx

integration-test:
	./integration.sh

$(COVERAGE_DIR):
	mkdir -p $(COVERAGE_DIR)


${COVERAGE_DIR}/coverage.txt:
	echo "mode: atomic" > ${COVERAGE_DIR}/coverage.txt
	-@for dir in $(shell $(go) list ./... | grep -v vendor |grep -v _vendor| grep -v vendor); do\
		$(go) test $(RACE) -coverprofile=${COVERAGE_DIR}/profile.out -covermode=atomic $$dir;\
		test ${COVERAGE_DIR}/$(dir).out && (cat ${COVERAGE_DIR}/profile.out | grep -v 'mode:' >> ${COVERAGE_DIR}/coverage.txt && rm ${COVERAGE_DIR}/profile.out)\
	done

test-coverage: $(COVERAGE_DIR) $(COVERAGE_DIR)/coverage.txt

refresh-coverage: clean test-coverage
	$(go) tool cover -html=$(COVERAGE_DIR)/coverage.txt -o $(COVERAGE_DIR)/coverage.html

all: install

# fetch depedencies
getdeps :
	$(go) get -u -v -d ./...

.PHONY: all clean test getdeps test-coverage
