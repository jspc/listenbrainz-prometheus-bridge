go_build_cmd ?= CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath

listenbrainz-prometheus-bridge: *.go go.*
	$(go_build_cmd) -o $@
