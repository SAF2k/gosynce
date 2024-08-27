package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

var (
	sourceDir      string
	backupDir      string
	schedulePeriod time.Duration
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Start the backup process",
	Long: `Start the rsync backup process with dynamic source and destination directories.
You can set the interval at which the backup is performed, ensuring that your data is securely stored at regular intervals.`,
	Run: func(cmd *cobra.Command, args []string) {
		schedule := gocron.NewScheduler(time.Local)

		// Schedule the backup job
		schedule.Every(schedulePeriod).Do(runBackup)

		// Start the scheduler
		schedule.StartBlocking()
	},
	Example: `  gorsync backup --source=/home/user/data --destination=/mnt/backup --interval=1h
  gorsync backup -s /path/to/source -d /path/to/destination -i 30m`,
}

func init() {
	// Adding backup command to the root command
	rootCmd.AddCommand(backupCmd)

	// Defining flags for the backup command
	backupCmd.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory (required)")
	backupCmd.Flags().StringVarP(&backupDir, "destination", "d", "", "Backup directory (required)")
	backupCmd.Flags().DurationVarP(&schedulePeriod, "interval", "i", 1*time.Hour, "Schedule interval (e.g., 1h, 30m)")

	// Marking flags as required
	backupCmd.MarkFlagRequired("source")
	backupCmd.MarkFlagRequired("destination")
}

func runBackup() {
	// Ensure source and destination directories are specified
	if sourceDir == "" || backupDir == "" {
		fmt.Println("Source or backup directory is not specified.")
		return
	}

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		fmt.Printf("Source directory '%s' does not exist.\n", sourceDir)
		return
	}

	// Create backup directory if it does not exist
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.MkdirAll(backupDir, os.ModePerm)
	}

	// Command to sync 2024 data with progress
	rsyncCmd := exec.Command("rsync", "-avh", "--progress", "--include=2024*/", "--include=2024*/**", "--exclude=*", fmt.Sprintf("%s/", sourceDir), fmt.Sprintf("%s/", backupDir))

	// Execute the rsync command
	fmt.Printf("Starting backup at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	rsyncCmd.Stdout = os.Stdout
	rsyncCmd.Stderr = os.Stderr
	err := rsyncCmd.Run()
	if err != nil {
		fmt.Printf("Error during backup: %v\n", err)
	} else {
		fmt.Printf("Backup completed at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}
