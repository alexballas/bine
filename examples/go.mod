module github.com/alexballas/bine/examples

go 1.25.0

require (
	github.com/alexballas/bine v0.0.0-20260617175114-f63422fc9667
	github.com/alexballas/go-libtor v1.0.8-0.20260627152619-3478b371d708
	golang.org/x/net v0.56.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260615183401-62b3387ff324 // indirect
)

replace github.com/alexballas/bine => ../
