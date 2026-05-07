package cli

import (
	"go.uber.org/zap"

	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/spf13/cobra"
)

// DeadLetterQ represents the structure of a dead letter queue message.
type DeadLetterQ struct {
	Handler    string `json:"handler_poisoned"`
	Reason     string `json:"reason_poisoned"`
	Subscriber string `json:"subscriber_poisoned"`
	Topic      string `json:"topic_poisoned"`
}

// GetDeadQueueCommandDef runs app
func GetDeadQueueCommandDef(cfg config.AppConfig, _ *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "dead-letter-queue",
		Short: "To start dead-letter queue",
		Long:  `This queue is used to store failed job from all worker`,
		RunE: func(_ *cobra.Command, _ []string) error {

			// Init worker
			subscriber, err := watermill.InitSubscriber(cfg, true)
			if err != nil {
				return err
			}

			// run worker with topic(queue name) and process function
			// it will run failed job until it get success
			err = subscriber.Run(cfg.MQ.DeadLetterQ, "dead_letter_queue", workers.Process)
			return err
		},
	}
	return workerCommand
}
