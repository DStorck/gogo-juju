package gogo

import "sync"

// Juju defines the cluster name, which bundle to use, and the manifest for credentials and clouda
type Juju struct {
	Name     string
	Bundle   string
	Manifest string
	p        Parallel
}

// Parallel sets the waitgroup if user wishes to bring up several clusters at once
type Parallel struct {
	wg sync.WaitGroup
}

// the following structs are for json parsing used with GetJujuStatus()

type jujuStatus struct {
	ApplicationResults map[string]applications `json:"applications"`
	Machines           map[string]machines     `json:"machines"`
}

type machines struct {
	MachStatus map[string]string `json:"juju-status"`
}

type applications struct {
	AppStatus map[string]string `json:"application-status"`
}
