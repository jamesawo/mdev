package doctor

import (
	"fmt"
	"sync"
	"time"

	"github.com/jamesawo/mdev/internal/ui/messages"
	"github.com/jamesawo/mdev/internal/ui/printer"
)

type noopReporter struct{}

func (noopReporter) StartSection(string)    {}
func (noopReporter) StartCheck(string)      {}
func (noopReporter) SystemCheck(Check)      {}
func (noopReporter) EnvironmentCheck(Check) {}
func (noopReporter) ToolCheck(ToolCheck)    {}

type progressReporter struct {
	mu      sync.Mutex
	stop    chan struct{}
	stopped chan struct{}
}

func (r *progressReporter) StartSection(title string) {
	r.stopCheck()
	printer.Section(title)
}

func (r *progressReporter) StartCheck(name string) {
	r.stopCheck()

	stop := make(chan struct{})
	stopped := make(chan struct{})

	r.mu.Lock()
	r.stop = stop
	r.stopped = stopped
	r.mu.Unlock()

	go spinCheck(name, stop, stopped)
}

func (r *progressReporter) SystemCheck(result Check) {
	r.stopCheck()
	if result.Status {
		printer.Success(result.Name)
		return
	}

	printer.Fail(fmt.Sprintf("%s %s", messages.Missing, result.Name))
}

func (r *progressReporter) EnvironmentCheck(result Check) {
	r.stopCheck()
	if result.Status {
		if result.Detail != "" {
			printer.Success(result.Name + ": " + result.Detail)
		} else {
			printer.Success(result.Name)
		}
		return
	}

	printer.Fail(messages.DoctorNotConfigured(result.Name))
	printer.Indent(2, messages.Run+" mdev install")
}

func (r *progressReporter) ToolCheck(result ToolCheck) {
	r.stopCheck()
	if result.Installed {
		printer.Success(result.Name)
		return
	}

	printer.Fail(result.Name)
}

func (r *progressReporter) stopCheck() {
	r.mu.Lock()
	stop := r.stop
	stopped := r.stopped
	r.stop = nil
	r.stopped = nil
	r.mu.Unlock()

	if stop == nil {
		return
	}

	close(stop)
	<-stopped
}

func spinCheck(name string, stop <-chan struct{}, stopped chan<- struct{}) {
	frames := []string{"|", "/", "-", "\\"}
	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	index := 0
	printer.OverwriteLine(printer.FormatIndent(1, fmt.Sprintf("%s checking %s...", frames[index], name)))

	for {
		select {
		case <-stop:
			printer.ClearLine()
			close(stopped)
			return
		case <-ticker.C:
			index = (index + 1) % len(frames)
			printer.OverwriteLine(printer.FormatIndent(1, fmt.Sprintf("%s checking %s...", frames[index], name)))
		}
	}
}

// renderSummary renders the final doctor summary after streaming progress.
func renderSummary(report *Report) {

	printer.Section(messages.DoctorNextSteps)

	printer.Indent(1, messages.DoctorInstallIndividual)
	for _, t := range report.Tools {
		if !t.Installed {
			printer.Indent(2, messages.DoctorInstallTool(t.Name))
		}
	}

	printer.Blank()
	printer.Info(messages.DoctorInstallEverything)
	printer.Indent(2, "mdev install --all")

	printer.Blank()
	printer.Info(messages.DoctorFixHint)
	printer.Indent(2, "mdev doctor --fix")
	printer.Blank()
}
