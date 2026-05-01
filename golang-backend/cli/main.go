package cli

import (
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/constants"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetAPICommandDef(&cfg, logger)
	workerCmd := GetWorkerCommandDef(cfg, logger)
	workerCmd.PersistentFlags().Int("retry-delay", 100, "time intertval for two retry in ms")
	workerCmd.PersistentFlags().Int("retry-count", 3, "number of retry")
	workerCmd.PersistentFlags().String("topic", constants.TopicWorkspaceInvites, "topic name(queue name)")

	deadQueueCmd := GetDeadQueueCommandDef(cfg, logger)
	rootCmd := &cobra.Command{Use: "nexus-tasks"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd, &workerCmd, &deadQueueCmd)
	return rootCmd.Execute()
}
