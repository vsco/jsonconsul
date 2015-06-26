// Adapted from
// https://github.com/hashicorp/envconsul/blob/master/runner.go

package jsonconsul

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	dep "github.com/hashicorp/consul-template/dependency"
	"github.com/hashicorp/consul-template/watch"
	"github.com/hashicorp/consul/api"
)

func (c *JsonExport) RunWatcher() {
	runner, err := NewRunner(c, false)
	if err != nil {
		log.Fatal(err)
	}

	runner.Start()
}

type Runner struct {
	sync.RWMutex

	// Prefix is the KeyPrefixDependency associated with this Runner.
	Prefix *dep.StoreKeyPrefix

	// ErrCh and DoneCh are channels where errors and finish notifications occur.
	ErrCh  chan error
	DoneCh chan struct{}

	// ExitCh is a channel for parent processes to read exit status values from
	// the child processes.
	ExitCh chan int

	// config is the Config that created this Runner. It is used internally to
	// construct other objects and pass data.
	config *JsonExport

	// client is the consul/api client.
	client *api.Client

	// once indicates the runner should get data exactly one time and then stop.
	once bool

	// minTimer and maxTimer are used for quiescence.
	minTimer, maxTimer <-chan time.Time

	// outStream and errStream are the io.Writer streams where the runner will
	// write information.
	outStream, errStream io.Writer

	// watcher is the watcher this runner is using.
	watcher *watch.Watcher

	// data is the latest representation of the data from Consul.
	data map[string][]*dep.KeyPair

	// killSignal is the signal to send to kill the process.
	killSignal os.Signal
}

// NewRunner accepts a JsonExport, and boolean value for once mode.
func NewRunner(config *JsonExport, once bool) (*Runner, error) {
	var err error

	log.Printf("[INFO] (runner) creating new runner (once: %v)", once)

	runner := &Runner{
		config: config,
		once:   once,
	}

	s := strings.TrimPrefix(config.Prefix, "/")
	runner.Prefix, err = dep.ParseStoreKeyPrefix(s)
	if err != nil {
		return nil, err
	}

	if err := runner.init(); err != nil {
		return nil, err
	}

	return runner, nil
}

// Start creates a new runner and begins watching dependencies and quiescence
// timers. This is the main event loop and will block until finished.
func (r *Runner) Start() {
	var (
		exitCh <-chan int
	)

	log.Printf("[INFO] (runner) starting")

	// Add the dependencies to the watcher
	r.watcher.Add(r.Prefix)

	for {
		select {
		case data := <-r.watcher.DataCh:
			r.Receive(data.Dependency, data.Data)

			// Drain all views that have data
		OUTER:
			for {
				select {
				case data = <-r.watcher.DataCh:
					r.Receive(data.Dependency, data.Data)
				default:
					break OUTER
				}
			}
		case <-r.minTimer:
			log.Printf("[INFO] (runner) quiescence minTimer fired")
			r.minTimer, r.maxTimer = nil, nil
		case <-r.maxTimer:
			log.Printf("[INFO] (runner) quiescence maxTimer fired")
			r.minTimer, r.maxTimer = nil, nil
		case err := <-r.watcher.ErrCh:
			// Intentionally do not send the error back up to the runner. Eventually,
			// once Consul API implements errwrap and multierror, we can check the
			// "type" of error and conditionally alert back.
			//
			// if err.Contains(Something) {
			//   errCh <- err
			// }
			log.Printf("[ERR] (runner) watcher reported error: %s", err)
		case <-r.watcher.FinishCh:
			log.Printf("[INFO] (runner) watcher reported finish")
			return
		case code := <-exitCh:
			r.ExitCh <- code
		case <-r.DoneCh:
			log.Printf("[INFO] (runner) received finish")
			return
		}

		// If we got this far, that means we got new data or one of the timers
		// fired, so attempt to re-process the environment.
		r.Run()
	}
}

// Stop halts the execution of this runner and its subprocesses.
func (r *Runner) Stop() {
	log.Printf("[INFO] (runner) stopping")
	r.watcher.Stop()

	close(r.DoneCh)
}

// Receive accepts data from Consul and maps that data to the prefix.
func (r *Runner) Receive(d dep.Dependency, data interface{}) {
	r.Lock()
	defer r.Unlock()
	r.data[d.HashCode()] = data.([]*dep.KeyPair)
}

// Run executes and manages the child process with the correct environment. The
// current enviornment is also copied into the child process environment.
func (r *Runner) Run() {
	log.Printf("[INFO] (runner) running")

	// TODO: Just call the app to consul again. This should
	// probably be updated to actually receive the values that we
	// got but MVP.
	r.config.Run()
}

// init creates the Runner's underlying data structures and returns an error if
// any problems occur.
func (r *Runner) init() error {
	// Print the final config for debugging
	result, err := json.MarshalIndent(r.config, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] (runner) final config (tokens suppressed):\n\n%s\n\n",
		result)

	r.client = client

	// Create the watcher
	watcher, err := newWatcher(r.config, client, r.once)
	if err != nil {
		return fmt.Errorf("runner: %s", err)
	}
	r.watcher = watcher

	r.data = make(map[string][]*dep.KeyPair)

	r.outStream = os.Stdout
	r.errStream = os.Stderr

	r.ErrCh = make(chan error)
	r.DoneCh = make(chan struct{})
	r.ExitCh = make(chan int, 1)

	return nil
}

// newWatcher creates a new watcher.
func newWatcher(config *JsonExport, client *api.Client, once bool) (*watch.Watcher, error) {
	log.Printf("[INFO] (runner) creating Watcher")

	clientSet := dep.NewClientSet()
	if err := clientSet.Add(client); err != nil {
		return nil, err
	}

	watcher, err := watch.NewWatcher(&watch.WatcherConfig{
		Clients: clientSet,
		Once:    once,
	})
	if err != nil {
		return nil, err
	}

	return watcher, err
}
