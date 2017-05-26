/*Package js provides a JavaScript sandbox for Acid.
 */
package js

import (
	"fmt"
	"time"

	"github.com/deis/acid/pkg/js/lib"
	"github.com/deis/quokka/pkg/javascript"
	"github.com/deis/quokka/pkg/javascript/libk8s"
)

// Sandbox gives access to a particular JavaScript runtime that is configured for Acid.
//
// Do not re-use sandboxes.
type Sandbox struct {
	rt *javascript.Runtime
}

// New creates a new *Sandbox
func New() (*Sandbox, error) {
	rt := javascript.NewRuntime()
	s := &Sandbox{
		rt: rt,
	}

	// Add the "built-in" libraries here:
	if err := libk8s.Register(rt.VM); err != nil {
		return s, err
	}

	// FIXME: This should make its way into quokka.
	rt.VM.Set("sleep", func(seconds int) {
		time.Sleep(time.Duration(seconds) * time.Second)
	})
	return s, nil
}

// Preload loads scripts that have been precompiled.
//
// The script must reside in lib.Scripts.
func (s *Sandbox) Preload(script string) error {
	data, ok := lib.Scripts[script]
	if !ok {
		return fmt.Errorf("unknown library: %s", script)
	}
	return s.ExecString(data)
}

// Variable Sets a variable in the runtime.
func (s *Sandbox) Variable(name string, val interface{}) {
	s.rt.VM.Set(name, val)
}

// ExecString executes the given string as a JavaScript file.
func (s *Sandbox) ExecString(script string) error {
	_, err := s.rt.VM.Run(script)
	return err
}

// ExecAll takes a list of scripts and executes them.
func (s *Sandbox) ExecAll(scripts ...[]byte) error {
	for _, script := range scripts {
		if _, err := s.rt.VM.Run(script); err != nil {
			return err
		}
	}

	return nil
}
