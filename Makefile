LDFLAGS += -X "github.com/amyangfei/go-logster/logster.ReleaseVersion=$(shell git describe --tags --dirty="-dev")"
LDFLAGS += -X "github.com/amyangfei/go-logster/logster.BuildTS=$(shell date -u '+%Y-%m-%d %H:%M:%S')"
LDFLAGS += -X "github.com/amyangfei/go-logster/logster.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "github.com/amyangfei/go-logster/logster.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "github.com/amyangfei/go-logster/logster.GoVersion=$(shell go version)"

PREFIX		:= /usr/local
DESTDIR		:= /usr/local
BINDIR		:= ${PREFIX}/bin
GO 			:= GO111MODULE=on go
GOBUILD 	:= $(GO) build
GOTEST		:= $(GO) test
PACKAGES	:= $$(go list ./... | grep -vE 'vendor')

BUILDDIR=build

APPS = logster
all: $(APPS)


$(BUILDDIR)/logster: $(wildcard apps/logster/*.go logster/*.go)

$(BUILDDIR)/%:
	@mkdir -p $(dir $@)
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $@ ./apps/$*
	@bash ./build_plugins.sh

$(APPS): %: $(BUILDDIR)/%

clean:
	rm -fr $(BUILDDIR)

.PHONY: install clean all test
.PHONY: $(APPS)

install: $(APPS)
	install -m 755 -d ${DESTDIR}${BINDIR}
	for APP in $^ ; do install -m 755 ${BUILDDIR}/$$APP ${DESTDIR}${BINDIR}/$$APP${EXT} ; done

test:
	$(GOTEST) -cover -race $(PACKAGES)

check:
	./test.sh
