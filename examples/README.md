### Bine Examples

The following examples are in this directory:

* [simpleclient](simpleclient) - A simple Tor client for connecting to the web or other onion services
* [simpleserver](simpleserver) - Hosting simple "hello world" Tor onion service
* [grpc](grpc) - Example showing how to use gRPC over Tor
* [httpaltsvc](httpaltsvc) - Example showing how to use .onion address as `Alt-Svc` of regular website (in development)

These examples use [go-libtor](https://github.com/alexballas/go-libtor) as the `ProcessCreator` so Tor is embedded
directly in the binary.

To run an example, while in this directory run the following with `<example>` replaced with the desired example:

    go run ./<example>
