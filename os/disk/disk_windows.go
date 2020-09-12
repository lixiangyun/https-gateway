// +build !linux,!darwin

package cpu

import (
	"fmt"
	"runtime"
)


// Stats represents disk I/O statistics for linux.
type Stats struct {
	Name            string // device name; like "hda"
	ReadsCompleted  uint64 // total number of reads completed successfully
	WritesCompleted uint64 // total number of writes completed successfully
}

// Get cpu statistics
func Get() (*Stats, error) {
	return nil, fmt.Errorf("disk statistics not implemented for: %s", runtime.GOOS)
}

type StatUsage struct {
	Name  string
	Total int
	Used  int
	Mount string
}

func GetDiskUsage() ([]StatUsage, error) {
	return nil, fmt.Errorf("disk statistics not implemented for: %s", runtime.GOOS)
}