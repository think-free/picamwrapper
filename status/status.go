package status

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/think-free/other/picamwrapper/config"
)

// Status object state checker
type Status struct {
	last     bool
	auto     bool
	conf     *config.Config
	internal *config.Internal
}

// New Create a new state checker
func New(conf *config.Config, internal *config.Internal) *Status {

	internal.Lock()
	internal.Auto = false
	internal.State = false
	internal.Unlock()
	return &Status{auto: false, conf: conf, internal: internal}
}

// Run the state checker
func (s *Status) Run() {

	s.last = !s.getState()

	go func() {
		for {
			select {
			case st := <-s.internal.Chwritest:

				fmt.Println("Writting manual status", st)
				s.setState(st)

			case st := <-s.internal.Chwriteauto:

				if s.auto {
					fmt.Println("Writting auto status", st)
					s.setState(st)
				}

			case auto := <-s.internal.AutoMode:
				fmt.Println("Automode status changed :", auto)
				s.internal.Lock()
				s.auto = auto
				s.internal.Auto = auto
				s.internal.Unlock()
				s.internal.Chgetauto <- auto
			}
		}
	}()

	for {
		current := s.getState()

		if s.last != current {
			s.internal.Lock()
			s.last = current
			s.internal.State = current
			s.internal.Unlock()
			s.internal.Chgetst <- current
		}

		time.Sleep(time.Second)
	}
}

func (s *Status) setState(st bool) {
	if st {
		cmd := exec.Command("touch", s.conf.PiCamPath+"/hooks/start_record")
		cmd.Run()
	} else {
		cmd := exec.Command("touch", s.conf.PiCamPath+"/hooks/stop_record")
		cmd.Run()
	}
}

func (s *Status) getState() bool {

	content, _ := ioutil.ReadFile(s.conf.PiCamPath + "/state/record")
	return strings.Contains(string(content), "true")
}
