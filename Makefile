VERSION ?= "${shell git tag --sort=-version:refname | head -n 1}-${shell git log -n 1 | head -n 1 | cut -d' ' -f2 | cut -b 1-16}"

GO ?= "go"
BIN ?= "lrcsnc"
DESTDIR ?= 
PREFIX ?= "/usr/local"
PACKAGEDIR ?= "${BIN}_${VERSION}"
PACKAGENAME := "${PACKAGEDIR}.tar.gz"

LDFLAGS_VERSION := -X lrcsnc/internal/setup.version=${VERSION}
LDFLAGS := \
	${LDFLAGS_VERSION}

default: build
all: build install clean

build:
	${GO} build -ldflags="${LDFLAGS}" -o ${BIN} -v
build-dev:
	${GO} build -ldflags="-X lrcsnc/internal/setup.version=dev" -o lrcsnc-dev -v
install:
	install -Dm644 LICENSE "${DESTDIR}/usr/share/licenses/${BIN}/LICENSE"
	install -Dm755 ${BIN} "${DESTDIR}${PREFIX}/bin/${BIN}"
package:
	strip ${BIN}
	cp -t ${PACKAGEDIR} ${BIN}
	tar -czvf ${PACKAGE} ${PACKAGEDIR}
clean:
	rm -f lrcsnc