package workers

import (
	"bytes"
	"encoding/gob"

	"github.com/ThreeDotsLabs/watermill/message"
)

func init() {
	for _, v := range RegisterWorkerStruct() {
		gob.Register(v)
	}
}

// RegisterWorkerStruct registers all worker structs for proper unmarshalling.
func RegisterWorkerStruct() []interface{} {
	return []interface{}{
		WelcomeMail{},
		WorkspaceInvitationMail{},
	}
}

// Handler interface for all worker struct
type Handler interface {
	Handle() error
}

// Process handles worker messages by decoding them into the appropriate struct and calling its Handle method.
func Process(msg *message.Message) error {
	buf := bytes.NewBuffer(msg.Payload)
	dec := gob.NewDecoder(buf)

	var result Handler
	err := dec.Decode(&result)
	if err != nil {
		return err
	}

	if err := result.Handle(); err != nil {
		return err
	}
	msg.Ack()
	return nil
}
