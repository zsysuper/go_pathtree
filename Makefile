#This how we want to name the binary output
BINARY=pathtree
GOPATH:=$(CURDIR)/
DATETIME=`date +'%y%m%d%H%M%S'`
DIR_SRC=./src/
DIR_BIN=./bin/
DIR_CONF=./conf/
DIR_PKG=./
DIR_LOG=./log/
PKG_FILE=${BINARY}_${DATETIME}.tar.gz
DIR_INSTALL=./installer_pkg

build:
	echo "GOPATH:${GOPATH}"
	GOPATH=${GOPATH} go build -o ${BINARY} ${DIR_SRC}/*.go
	mkdir -p $(DIR_BIN) && mv -f ${BINARY} $(DIR_BIN)

build_linux64:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY} ${DIR_SRC}/*.go
	mkdir -p $(DIR_BIN) && mv -f ${BINARY} $(DIR_BIN)

pkg:
	mkdir -p ${DIR_PKG}
	tar -zcvPf ${PKG_FILE} ${DIR_BIN} ${DIR_CONF}

install:
	make clean
	make -j8
	make pkg
	mkdir -p ${DIR_INSTALL}
	mv -f ${DIR_PKG}pathtree_*.tar.gz ${DIR_INSTALL}
debug:
	make clean
	make -j8
	./${DIR_BIN}/${BINARY}

clean:
	rm -f ${DIR_PKG}/*.tar.gz ${DIR_BIN}/* ${DIR_INSTALL}/*.tar.gz
.PHONY: clean pkg debug install