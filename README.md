Fluentd Monitor UI
========================

UI for fluentd plugin monitoring API (http://docs.fluentd.org/articles/monitoring#monitoring-agent).

### Running

The deployed config file will be located in `/etc/fluentd-monitor/`. This should
be updated with the target fluentd hosts.

Once the config is updated `service fluentd-monitor start` and check localhost:8080 on the host.


### Deploying

Distributable packages (RPM etc.) can be built using the `package.sh` script. Installing from
source is possible using `make install`.

### Developing

The front end files are by default embedded in static.go. They can be updated and rebuilt
by doing the following:

1. Checkout npm/bower sources `cd ui && npm install && bower install`
2. Rebuild embedded files (from project root) `${GOPATH}/bin/esc -prefix="ui/static" -o static.go ui/static`
