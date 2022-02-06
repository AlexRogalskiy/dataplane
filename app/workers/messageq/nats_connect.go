package messageq

import (
	"dataplane/workers/logging"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var NATS *nats.Conn
var NATSencoded *nats.EncodedConn
var NATSProtobuf *nats.EncodedConn
var NATSStream nats.JetStreamContext

func NATSConnect() {
	var err error
	NATS, err = nats.Connect(os.Getenv("dataplane_nats"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(3*time.Second))

	if err != nil {
		logging.PrintSecretsRedact(err.Error())
		log.Fatal("Failed to connect to NATS")
	}

	NATSencoded, _ = nats.NewEncodedConn(NATS, nats.JSON_ENCODER)
	if err != nil {
		logging.PrintSecretsRedact(err.Error())
		log.Fatal("Failed to connect to encoded NATS")
	}

	// NATSProtobuf, _ = nats.NewEncodedConn(NATS, protobuf.PROTOBUF_ENCODER)
	// if err != nil {
	// 	logging.PrintSecretsRedact(err.Error())
	// 	log.Fatal("Failed to connect to encoded NATS")
	// }

	// Create JetStream Context
	// NATSStream, _ = NATS.JetStream(nats.PublishAsyncMaxPending(256))

	log.Println("📧 NATS connected")
	//DBConn.Config.PrepareStmt = true
}
