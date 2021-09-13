package services

import (
	"bufio"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	ssh2 "golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"task-runner/connection/ssh"
)

type PortBinding struct {
	HostIp   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

type InspectItem struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State struct {
		Status string `json:"Status"`
	}
	NetworkSettings struct {
		Ports map[string][]PortBinding
	}
}

func PrintInfo(client *ssh.Client, exportToFile string) {
	client.Connect()
	defer func(Connection *ssh2.Client) {
		err := Connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client.Connection)

	session, err := client.Connection.NewSession()

	if err != nil {
		log.Fatal(err)
	}

	modes := ssh2.TerminalModes{
		ssh2.ECHO:          0,
		ssh2.TTY_OP_ISPEED: 14400,
		ssh2.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatal(err)
	}

	out, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	_, err = session.Output("docker inspect --format='{{json .}}' $(sudo docker ps -aq)")
	reader := bufio.NewReader(out)
	var inspectItems []InspectItem
	for {
		readString, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		var inspectItem InspectItem
		err = json.Unmarshal([]byte(readString), &inspectItem)
		if err != nil {
			log.Fatal(err)
		}
		inspectItems = append(inspectItems, inspectItem)
	}

	outSource := os.Stdout
	if exportToFile != "" {
		file, err := os.Create(exportToFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		outSource = file
	}

	table := tablewriter.NewWriter(outSource)
	table.SetHeader([]string{"Service name", "Branch", "Internal port", "Host IP", "Host port", "Status"})
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
		info := strings.SplitN(inspectItem.Name, "_", 2)
		name := info[0]
		if strings.Contains(name, "service") {
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
				strings.Trim(name, "/"),
				info[1],
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
}
