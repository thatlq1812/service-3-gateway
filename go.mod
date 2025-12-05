module service-3-gateway

go 1.25.4

replace github.com/thatlq1812/service-1-user => ../service-1-user

replace github.com/thatlq1812/service-2-article => ../service-2-article

replace agrios => ..

require (
	github.com/gorilla/mux v1.8.1
	github.com/thatlq1812/service-1-user v1.2.3
	github.com/thatlq1812/service-2-article v1.2.3
	google.golang.org/grpc v1.77.0
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
