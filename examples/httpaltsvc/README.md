## HTTP Alt-Svc Example

The Tor browser now supports `Alt-Svc` headers to be onion services as
[HTTP alternate services](https://tools.ietf.org/html/rfc7838) to traditional sites.
[This Cloudflare post](https://blog.cloudflare.com/cloudflare-onion-service/) explains how they use it. This example
shows how to run the same pattern locally.

Specifically, this example listens on all IPs on 80 for insecure HTTP requests. It also listens on an onion service for
insecure HTTP requests. Both services return `Alt-Svc` addresses to an onion service that is run securely.

**NOTE: Only do the steps in this example if you are comfortable with and aware of the consequences.**

### Setup

We're going to self-sign a certificate for use in the Tor browser. First, download `mkcert` and run:

    mkcert -install

This will install a fake CA on your local machine that certs can be generated from. Remember the path it says the local
CA is at or run `mkcert -CAROOT` to get it back. We will use this "CA path" later.

Now the CA must be added to the Tor browser. By default the Tor browser doesn't use the cert database so it cannot store
the overrides. To change this, go to `about:config` click through warning and change `security.nocertdb` to `false`. Now
that CAs can be added, go to `Options` (i.e. `about:preferences`) > `Privacy & Security` > `Certificates` section at the
bottom > `View Certificates...` > `Authorities` tab > `Import...` > choose `rootCA.pem` from the "CA path" from
earlier > check "Trust this CA to identify websites" and click `OK`. Restart the Tor browser.

### Running

Make sure `tor` and `mkcert` are both on your `PATH`. The example binds to port 80, so run it with the permissions
needed to bind a privileged port on your system.

Point your DNS to this machine's IP, or start `ngrok http 80` via [ngrok.com](https://ngrok.com/) and use the generated
hostname as the domain argument. Then run:

    go run ./httpaltsvc [-verbose] example.com

Open `http://example.com` in Tor Browser, replacing `example.com` with the domain you passed to the command. The first
request reaches the regular HTTP listener, which advertises the secure onion service with an `Alt-Svc` header. If the
browser accepts that alternative service, later requests to the same visible URL should be served by the secure onion
listener while the address bar still shows the original domain.

The command also prints the two generated onion addresses. You can open the first onion over HTTP to see it advertise the
second onion as its alternative service, or open the second onion over HTTPS directly to check the certificate and server
output. Press enter in the terminal to stop the example. Generated Tor state and certificates are written under
`tor-data/`.

Tor Browser's handling of onion `Alt-Svc` has varied across versions and configuration, so verify the response body and
printed request logs when testing. Requests served over the onion listener should not show a Tor exit node as the remote
address.
