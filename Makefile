PREFIX=/usr
GOBIN=${DESTDIR}${PREFIX}/bin
PROJ=fluentd-monitor
PACKAGE_TYPE=rpm
PACKAGE_BUILD_DIR=pkg
PACKAGE_DIR=dist

static:

	${GOPATH}/bin/esc -prefix="ui/static" -o static.go ui/static

build:

	go get
	go build

install: build

	#install binary
	GOBIN=${GOBIN} go install -v

	#install config file
	install -Dm 644 config/config.yaml ${DESTDIR}/etc/${PROJ}/config.yaml

	#install init script
	install -Dm 755 init.d/${PROJ} ${DESTDIR}/etc/init.d/${PROJ}

package:

	#
	# export PACKAGE_TYPE to vary package type (e.g. deb, tar, rpm)
	#

	@if [ -z "$(shell which fpm 2>/dev/null)" ]; then \
		echo "error:\nPackaging requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"; \
		exit 1; \
	fi

	#run make install against the packaging dir
	$(MAKE) install DESTDIR=${PACKAGE_BUILD_DIR}

	#clean
	mkdir -p dist && rm -f dist/*.${PACKAGE_TYPE}

	#build package
	fpm --rpm-os linux \
		-s dir \
		-p dist \
		--config-files /etc/fluentd-monitor/config.yaml \
		-t ${PACKAGE_TYPE} \
		-n fluentd-monitor \
		-v $(shell cat version) \
		-C ${PACKAGE_BUILD_DIR} .
