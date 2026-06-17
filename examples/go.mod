module github.com/alexballas/bine/examples

go 1.25.0

require (
	github.com/alexballas/bine v0.0.0-00010101000000-000000000000
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
