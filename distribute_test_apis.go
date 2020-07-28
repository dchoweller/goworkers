package main

/* Sample program to demonstrate sending Avant-garde test APIs
 * to IFE hosts using Go channels and select statement.
 */
import "fmt"

// parameter of a test API
type apiParam struct {
	key   string
	value string
}

// Test API
type testAPI struct {
	component  string     // CSCG component
	name       string     // name of test API
	parameters []apiParam // list of parameters
}

// stringer interface for test API
func (t testAPI) String() string {
	return fmt.Sprintf("%v.%v", t.component, t.name)
}

// Host information
type hostInfo struct {
	name string
	ip   string
}

// stringer interface for host information
func (h hostInfo) String() string {
	return fmt.Sprintf("%v(%v)", h.name, h.ip)
}

/*
 * executTestAPI waits on 'apiChannel' for APIs to be executed.
 * If any comand is received on 'commandChannel', write confirmation and
 * exit.
 */
func executeTestAPI(host hostInfo, apiChannel chan testAPI, commandChannel chan int) {
	for {
		select {
		case rcvdAPI := <-apiChannel: // read API from channel
			fmt.Printf("Sending %v to host %v\n", rcvdAPI, host)
		case <-commandChannel: // read command from channel
			fmt.Println("Ending test API executions at", host)
			commandChannel <- 0
			return
		}

	}

}

/*
 * Execute 'numAPI' test APIs in 'apis' on all 'numHosts' hosts in 'hosts'
 */
func executeTestApis(apis []testAPI, hosts []hostInfo, numAPIs int, numHosts int) {
	//var apiChannels = [numHosts]chan testAPI{}
	apiChannels := make([]chan testAPI, numHosts)
	commandChannels := make([]chan int, numHosts)

	// Initialize communication channels
	for h := 0; h < numHosts; h++ {
		// Channel to send test APIs
		apiChannels[h] = make(chan testAPI)
		// Channel to send/receive commands/results
		commandChannels[h] = make(chan int)
		// Spawn worker thread for each host to execute APIs
		go executeTestAPI(hosts[h], apiChannels[h], commandChannels[h])
	}
	// For each test API
	for a := 0; a < numAPIs; a++ {
		// Send it to all hosts to be executed
		for h := 0; h < numHosts; h++ {
			apiChannels[h] <- apis[a]
		}
	}
	// Terminate test API execution
	for h := 0; h < numHosts; h++ {
		commandChannels[h] <- 1
	}
	// Gather results
	for h := 0; h < numHosts; h++ {
		<-commandChannels[h]
	}

}

func main() {
	var apis = []testAPI{
		{"fwk.vod", "play",
			[]apiParam{
				apiParam{"name", "movie.mpg"},
				apiParam{"loop", "false"},
			}},
		{"csw.gsm", "getModedata",
			[]apiParam{
				apiParam{"parameter", "hbMode"},
			}},
	}
	var hosts = []hostInfo{
		{"dsu1", "139.182.68.1"},
		{"dsu2", "139.182.68.2"},
	}
	executeTestApis(apis, hosts, 2, 2)
	fmt.Println("Done.")
}
