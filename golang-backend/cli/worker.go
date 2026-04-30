package cli

import (
	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/pkg/watermill"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetWorkerCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	workerCmd := cobra.Command{
		Use:   "worker",
		Short: "Start the background worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get details from flag
			topic, err := cmd.Flags().GetString("topic")
			if err != nil {
				return err
			}

			retryCount, err := cmd.Flags().GetInt("retry-count")
			if err != nil {
				return err
			}

			delay, err := cmd.Flags().GetInt("retry-delay")
			if err != nil {
				return err
			}

			// Init subscriber
			subscriber, err := watermill.InitSubscriber(cfg, false)
			if err != nil {
				return err
			}

			// init router for add middleware,retry count,etc
			router, err := subscriber.InitRouter(cfg, delay, retryCount)
			if err != nil {
				return err
			}

			// run worker with topic(queue name) and process function
			err = router.Run(topic, cfg.MQ.HandlerName, workers.Process)
			return err

		},
	}
	return workerCmd
}
