PREFIX=/usr/local
DESTDIR=
GOFLAGS=
BINDIR=${PREFIX}/bin

BLDDIR = build
EXT=
ifeq (${GOOS},windows)
    EXT=.exe
endif

APPS = agent monitor dashboard
all: $(APPS)

$(BLDDIR)/%:
	@mkdir -p $(dir $@)
	go build ${GOFLAGS} -o $@ ./cmd/$*
	@cp ./cmd/agent/agent.toml ./build/agent.toml
	@cp ./cmd/monitor/monitor.toml ./build/monitor.toml
	@cp ./cmd/dashboard/dashboard.toml ./build/dashboard.toml

$(APPS): %: $(BLDDIR)/%

clean:
	rm -fr $(BLDDIR)

install: $(APPS)
	install -m 755 -d ${DESTDIR}${BINDIR}
	for APP in $^ ; do install -m 755 ${BLDDIR}/$$APP ${DESTDIR}${BINDIR}/$$APP${EXT} ; done
