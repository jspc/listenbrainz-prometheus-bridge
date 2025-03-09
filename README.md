# listenbrainz-prometheus-bridge

Grab latest statistics for the artists, tracks played by a set of users and expose those statistics via the OpenMetrics line format, ahead of being scraped by prometheus.

```bash
$ export LISTENBRAINZ_TOKEN=""  # See: https://listenbrainz.org/settings/
$ export LISTENBRAINZ_USERS=""  # CSV of the users to track
$ go build
$ ./listenbrainz-prometheus-bridge
```

## Configuration

* `LISTENBRAINZ_TOKEN` - API Token used to authenticate against Listenbrainz. This is a vaguely UUID v4 looking string
* `LISTENBRAINZ_USERS` - A comma seperated string of the users to grab metrics for; for my use these are the listenbrainz users that correspond to users on my navidrome server


## Testing

Hahahahahahahaha
