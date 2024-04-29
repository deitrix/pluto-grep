package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"text/tabwriter"
)

func main() {
	servers, err := getServers()
	if err != nil {
		panic(err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "IP\tGame\tHostname\tMap\tGametype\tPlayers\tMaxPlayers\tHardcore\tPassword\tBots\tVoice\tDescription\tRevision")
	slices.SortFunc(servers, func(a, b Server) int {
		return cmp.Compare(len(b.Players), len(a.Players))
	})

	for _, server := range servers {
		if len(server.Players) == 0 {
			continue
		}
		if server.Game != "t6mp" {
			continue
		}
		if server.Revision != 4035 {
			continue
		}
		fmt.Fprintf(tw,
			"%s\t%s\t%s\t%s\t%s\t%d\t%d\t%t\t%t\t%d\t%d\t%s\t%d\n",
			server.IP+":"+strconv.Itoa(server.Port),
			server.Game,
			server.Hostname,
			server.Map,
			server.GameType,
			len(server.Players),
			server.MaxPlayers,
			server.Hardcore,
			server.Password,
			server.Bots,
			server.Voice,
			server.Description,
			server.Revision,
		)
	}
}

func getServers() ([]Server, error) {
	req, err := http.NewRequest(http.MethodGet, "https://plutonium.pw/api/servers", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var servers []Server
	if err := json.NewDecoder(res.Body).Decode(&servers); err != nil {
		return nil, err
	}

	return servers, nil
}

type Server struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Game     string `json:"game"`
	Hostname string `json:"hostname"`
	Map      string `json:"map"`
	GameType string `json:"gametype"`
	Players  []struct {
		Username string `json:"username"`
		ID       int    `json:"id"`
		Ping     int    `json:"ping"`
	} `json:"players"`
	MaxPlayers  int    `json:"maxplayers"`
	Hardcore    bool   `json:"hardcore"`
	Password    bool   `json:"password"`
	Bots        int    `json:"bots"`
	Voice       int    `json:"voice"`
	Description string `json:"description"`
	CodInfo     string `json:"codInfo"`
	Revision    int    `json:"revision"`
}
