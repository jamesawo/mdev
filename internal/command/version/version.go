// Package version renders deterministic metadata for the running mdev binary.
package version

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jamesawo/mdev/internal/ui/messages"
)

const (
	defaultVersion = "0.4.2"
	unknownValue   = "unknown"
)

// Metadata identifies an mdev binary.
type Metadata struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Built   string `json:"built"`
}

// Options controls how version metadata is rendered.
type Options struct {
	JSON bool
}

// DefaultMetadata returns the deterministic metadata used by local builds.
func DefaultMetadata() Metadata {
	return Metadata{
		Version: defaultVersion,
		Commit:  unknownValue,
		Built:   unknownValue,
	}
}

// Run writes the supplied metadata in human-readable or JSON form.
func Run(out io.Writer, metadata Metadata, options Options) error {
	metadata = withDefaults(metadata)

	var output []byte
	if options.JSON {
		encoded, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		output = append(encoded, '\n')
	} else {
		output = []byte(fmt.Sprintf(
			messages.VersionProductLineFormat+messages.VersionCommitLineFormat+messages.VersionBuiltLineFormat,
			metadata.Version,
			metadata.Commit,
			metadata.Built,
		))
	}

	n, err := out.Write(output)
	if err != nil {
		return fmt.Errorf(messages.VersionWriteFailed, err)
	}
	if n != len(output) {
		return fmt.Errorf(messages.VersionWriteFailed, io.ErrShortWrite)
	}
	return nil
}

func withDefaults(metadata Metadata) Metadata {
	defaults := DefaultMetadata()
	if metadata.Version == "" {
		metadata.Version = defaults.Version
	}
	if metadata.Commit == "" {
		metadata.Commit = defaults.Commit
	}
	if metadata.Built == "" {
		metadata.Built = defaults.Built
	}
	return metadata
}
