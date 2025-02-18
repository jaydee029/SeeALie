module github.com/jaydee029/SeeALie/notifications

go 1.21.6

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.66.2
	google.golang.org/protobuf v1.34.2
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/jaydee029/SeeALie/pubsub v0.0.0-20250129115149-15179fdf1ea8 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
replace github.com/jaydee029/SeeALie/pubsub => ../pubsub