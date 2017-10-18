go ?= go
GO_LD_X = -X $(1)=$(2)
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_ID), $(call GO_LD_X,github.com/mistsys/accord.ClientID,'$(GOOGLE_CLIENT_ID)'))
GO_LDFLAGS += $(if $(GOOGLE_CLIENT_SECRET), $(call GO_LD_X,github.com/mistsys/accord.ClientSecret,'$(GOOGLE_CLIENT_SECRET)'))
set_vars += $(if $(GO_LDFLAGS), -ldflags="$(GO_LDFLAGS)")


install:
	$(go) install $(set_vars) ./...

# clean
clean :
	-$(go) clean -i ./...

test:
	$(go) test ./...

all: install

# fetch depedencies
getdeps :
	$(go) get -u -v -d ./...

.PHONY: all clean test getdeps
