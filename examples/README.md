# Sample App

The `sample-app` directory is a sample splunk app that exposes two REST endpoints:

1. One using the python persistent server framework, and is exposed at `servicesNS/nobody/sample-app/py_interface`
2. One using the `splunk-persistentconn` framework written in `go` and is exposed at `servicesNS/nobody/sample-app/go_interface`

The `go-app` directory contains all the code of your actual application.

## How to build and deploy

1. To build the actual application in `go-app`

```bash
cd go-app
go build -o ./server main.go
```

2. Move the compiled executable into the sample Splunk app

```bash
mv ./go-app/server ./sample-app/bin
```

3. Deploy the sample Splunk app

```bash
cp -r ./sample-app $SPLUNK_HOME/etc/apps
$SPLUNK_HOME/bin/splunk restart
```