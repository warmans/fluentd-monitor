PREFIX=/usr/local
GOBIN=${DESTDIR}${PREFIX}/bin
PROJ=fluentd-monitor

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
