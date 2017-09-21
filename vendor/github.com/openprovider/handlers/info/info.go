package info

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/takama/router"
)

// ServiceInfo defines HTTP API response giving service information
type ServiceInfo struct {
	Host    string       `json:"host"`
	Runtime *RuntimeInfo `json:"runtime"`
	Version string       `json:"version"`
	Repo    string       `json:"repo"`
	Commit  string       `json:"commit"`
}

// RuntimeInfo defines runtime part of service information
type RuntimeInfo struct {
	Compiler   string `json:"compilier"`
	CPU        int    `json:"cpu"`
	Memory     string `json:"memory"`
	Goroutines int    `json:"goroutines"`
}

// Handler provides JSON API response giving service information
func Handler(version, repo, commit string) router.Handle {
	return func(c *router.Control) {
		host, _ := os.Hostname()
		m := new(runtime.MemStats)
		runtime.ReadMemStats(m)

		rt := &RuntimeInfo{
			CPU:        runtime.NumCPU(),
			Memory:     fmt.Sprintf("%.2fMB", float64(m.Alloc)/(1<<(10*2))),
			Goroutines: runtime.NumGoroutine(),
		}

		info := ServiceInfo{
			Host:    host,
			Runtime: rt,
			Version: version,
			Repo:    repo,
			Commit:  commit,
		}

		c.Code(http.StatusOK).Body(info)
	}
}
