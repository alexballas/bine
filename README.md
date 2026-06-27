# <img src="logo/bine-logo.png" width="180px">

[![Go Reference](https://pkg.go.dev/badge/github.com/alexballas/bine.svg)](https://pkg.go.dev/github.com/alexballas/bine)

Bine is a Go API for using and controlling Tor. It is similar to [Stem](https://stem.torproject.org/).

This is a maintained fork of [cretz/bine](https://github.com/cretz/bine).

Features:

* Full support for the Tor controller API
* Support for `net.Conn` and `net.Listen` style APIs
* Supports statically compiled Tor to embed Tor into the binary (via [go-libtor](https://github.com/alexballas/go-libtor))
* Supports v3 onion services
* Support for embedded control socket in Tor >= 0.3.5 (non-Windows)

See info below, the [API docs](https://pkg.go.dev/github.com/alexballas/bine), and the [examples](examples). The project is
MIT licensed. The Tor docs/specs and https://github.com/yawning/bulb were great helps when building this.

## Example

It is really easy to create an onion service. For example, using
[go-libtor](https://github.com/alexballas/go-libtor) to embed Tor directly in the binary, this bit of code will show a
directory server of the current directory:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexballas/bine/tor"
	"github.com/alexballas/go-libtor"
)

func main() {
	// Start Tor with go-libtor's embedded process creator.
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(nil, &tor.StartConf{ProcessCreator: libtor.Creator})
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer t.Close()
	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()
	// Create a v3 onion service to listen on any port but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	defer onion.Close()
	fmt.Printf("Open Tor browser and navigate to http://%v.onion\n", onion.ID)
	fmt.Println("Press enter to exit")
	// Serve the current folder from HTTP
	errCh := make(chan error, 1)
	go func() { errCh <- http.Serve(onion, http.FileServer(http.Dir("."))) }()
	// End when enter is pressed
	go func() {
		fmt.Scanln()
		errCh <- nil
	}()
	if err = <-errCh; err != nil {
		log.Panicf("Failed serving: %v", err)
	}
}
```

If in `main.go` it can simply be run with `go run main.go`. The example uses
[go-libtor](https://github.com/alexballas/go-libtor), which bundles Tor and its C dependencies and exposes a
`ProcessCreator`.

In non-Windows environments, the `UseEmbeddedControlConn` field in `StartConf` can be set to `true` to use an embedded
socket that does not open a control port. With Tor statically linked the binary does not have to be distributed
separately. Of course take notice of all licenses in accompanying projects.

## Forwarding to existing services

`Listen` is best when you want Tor to back a single Go `net.Listener`. When you already run local services and want to
map several onion (virtual) ports to different local addresses, use `Forward` instead. The `PortForwards` map keys are
local addresses and the values are the remote onion ports they serve, so the following exposes onion port `80` from a
service on `127.0.0.1:5000` and onion port `90` from a service on `127.0.0.1:5001`:

```go
fwd, err := t.Forward(ctx, &tor.ForwardConf{
	PortForwards: map[string][]int{
		"127.0.0.1:5000": {80},
		"127.0.0.1:5001": {90},
	},
})
```

By default the onion service is deleted (`DEL_ONION`) when `OnionService.Close`/`OnionForward.Close` is called. Set
`Detach: true` together with `NoDeleteOnClose: true` on the `ListenConf`/`ForwardConf` to keep the service published
after the controller connection closes.

## Testing

To test, a simple `go test ./...` from the base of the repository will work (add in a `-v` in there to see the tests).
The integration tests in `tests` however will be skipped. To execute those tests, `-tor` must be passed to the test.
Also, `tor` must be on the `PATH` or `-tor.path` must be set to the path of the `tor` executable. Even with those flags,
only the integration tests that do not connect to the Tor network are run. To also include the tests that use the Tor
network, add the `-tor.network` flag. For details Tor logs during any of the integration tests, use the `-tor.verbose`
flag.
