PREFIX=/usr/local
DESTDIR=/usr/local
BINDIR=${PREFIX}/bin

BUILDDIR=build

APPS = logster
all: $(APPS)


$(BUILDDIR)/logster: $(wildcard apps/logster/*.go logster/*.go)

$(BUILDDIR)/%:
	@mkdir -p $(dir $@)
	GO111MODULE=on go build ${GOFLAGS} -o $@ ./apps/$*
	@bash ./build_plugins.sh

$(APPS): %: $(BUILDDIR)/%

clean:
	rm -fr $(BUILDDIR)

.PHONY: install clean all
.PHONY: $(APPS)

install: $(APPS)
	install -m 755 -d ${DESTDIR}${BINDIR}
	for APP in $^ ; do install -m 755 ${BUILDDIR}/$$APP ${DESTDIR}${BINDIR}/$$APP${EXT} ; done
