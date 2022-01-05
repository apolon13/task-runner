package info

import (
	"bufio"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"sort"
	"strings"
	sshConnection "task-runner-cobra/utils/connection/ssh"
)

type portBinding struct {
	HostIp   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

type inspectItem struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State struct {
		Status string `json:"Status"`
	}
	NetworkSettings struct {
		Ports map[string][]portBinding
	}
}

func init() {
	Cmd.PersistentFlags().String("export", "", "export to file")
}

var (
	Cmd = &cobra.Command{
		Use:   "grpc-info",
		Short: "grpc info",
		Long:  "Show a table with available grps services.",
		Run: func(cmd *cobra.Command, args []string) {
			sshClient := sshConnection.NewClient("grpc")
			exportToFile, _ := cmd.Flags().GetString("export")
			sshClient.Connect()
			defer func(Connection *ssh.Client) {
				err := Connection.Close()
				if err != nil {
					panic(err)
				}
			}(sshClient.Connection)

			session, err := sshClient.Connection.NewSession()

			if err != nil {
				panic(err)
			}

			modes := ssh.TerminalModes{
				ssh.ECHO:          0,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}

			err = session.RequestPty("xterm", 80, 40, modes)
			if err != nil {
				panic(err)
			}

			out, err := session.StdoutPipe()
			if err != nil {
				panic(err)
			}

			_, err = session.Output("docker inspect --format='{{json .}}' $(sudo docker ps -aq)")
			reader := bufio.NewReader(out)
			var inspectItems []inspectItem
			for {
				readString, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				var inspectItem inspectItem
				err = json.Unmarshal([]byte(readString), &inspectItem)
				if err != nil {
					panic(err)
				}
				inspectItems = append(inspectItems, inspectItem)
			}

			outSource := os.Stdout
			if exportToFile != "" {
				file, err := os.Create(exportToFile)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				outSource = file
			}

			table := tablewriter.NewWriter(outSource)
			table.SetHeader([]string{"Service name", "Internal port", "Host IP", "Host port", "Status"})
			var items [][]string
			sort.Slice(inspectItems, func(i, j int) bool {
				if strings.Compare(inspectItems[i].Name, inspectItems[j].Name) == 1 {
					return true
				}
				return false
			})
			sort.Slice(inspectItems, func(i, j int) bool {
				return inspectItems[i].Name[1] < inspectItems[j].Name[1]
			})
			for _, inspectItem := range inspectItems {
				if strings.Contains(inspectItem.Name, "service") {
					var internalPorts []string
					var hostIps []string
					var hostPorts []string
					for internalPort, bindings := range inspectItem.NetworkSettings.Ports {
						internalPorts = append(internalPorts, internalPort)
						for _, binding := range bindings {
							hostIps = append(hostIps, binding.HostIp)
							hostPorts = append(hostPorts, binding.HostPort)
						}
					}
					items = append(items, []string{
						inspectItem.Name,
						strings.Join(internalPorts, "\n"),
						strings.Join(hostIps, "\n"),
						strings.Join(hostPorts, "\n"),
						inspectItem.State.Status,
					})
				}
			}
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.AppendBulk(items)
			table.Render()
		},
	}
)
