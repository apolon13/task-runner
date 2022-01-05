package restore

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"task-runner-cobra/cmd/custom"
	file "task-runner-cobra/interface"
	commandOut "task-runner-cobra/utils/command/custom/out"
	s3Connection "task-runner-cobra/utils/connection/s3"
	sshConnection "task-runner-cobra/utils/connection/ssh"
	s3Downloader "task-runner-cobra/utils/downloader/s3"
	sftpDownloader "task-runner-cobra/utils/downloader/sftp"
)

type restoreProcess struct {
	DumpFile    file.DownloadFile
	Commands    []*custom.Command
	RemoveAfter bool
}

func init() {
	Cmd.PersistentFlags().String("f", "", "dump file name")
	_ = Cmd.MarkPersistentFlagRequired("f")
	Cmd.PersistentFlags().String("db", "", "database name")
	_ = Cmd.MarkPersistentFlagRequired("db")
	Cmd.PersistentFlags().String("con", "ssh", "using connection")
}

var (
	Cmd = &cobra.Command{
		Use:   "restore-db",
		Short: "Restore you database",
		Long:  "Automatic database download and restore. Use custom restore commands in you configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			var df file.DownloadFile

			con, _ := cmd.Flags().GetString("con")
			f, _ := cmd.Flags().GetString("f")
			db, _ := cmd.Flags().GetString("db")
			restorePath := viper.GetStringMapString("restore.db.path." + con)
			switch con {
			case "ssh":
				df = sftpDownloader.NewFile(
					restorePath["local"],
					restorePath["remote"],
					f,
					sshConnection.NewClient("db"))
			case "s3":
				df = s3Downloader.NewFile(
					restorePath["local"],
					restorePath["remote"],
					f,
					s3Connection.NewClient("db"))
			}

			customCommands := custom.BuildCommandsWithViaConfig("restore.db.commands")
			for _, customCommand := range customCommands {
				customCommand.ReplaceArgs(map[string]string{
					"<-f>":  df.GetFileName(),
					"<-db>": db,
				})
			}

			restoreProcess := &restoreProcess{
				df,
				customCommands,
				viper.GetBool("restore.db.remove"),
			}

			defer func() {
				if restoreProcess.RemoveAfter == true {
					if err := restoreProcess.DumpFile.RemoveLocal(); err != nil {
						panic(fmt.Errorf("error removing downloaded file: %s", err))
					}
				}
			}()

			fmt.Println(fmt.Sprintf("Download file - %s", restoreProcess.DumpFile.GetFileName()))
			restoreProcess.DumpFile.Process()
			fmt.Println("\nDownloading complete")
			commandOut.HandleBatch(restoreProcess.Commands)
		},
	}
)
