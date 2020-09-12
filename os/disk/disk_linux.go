// +build linux

package disk

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/lixiangyun/autoproxy/proc"
	"github.com/lixiangyun/autoproxy/util"
	"io"
	"os"
	"strconv"
	"strings"
)

// Get disk I/O statistics.
func Get() ([]Stats, error) {
	// Reference: Documentation/iostats.txt in the source of Linux
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return collectDiskStats(file)
}

// Stats represents disk I/O statistics for linux.
type Stats struct {
	Name            string // device name; like "hda"
	ReadsCompleted  uint64 // total number of reads completed successfully
	WritesCompleted uint64 // total number of writes completed successfully
}

func collectDiskStats(out io.Reader) ([]Stats, error) {
	scanner := bufio.NewScanner(out)
	var diskStats []Stats
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		name := fields[2]
		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reads completed of %s", name)
		}
		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse writes completed of %s", name)
		}
		diskStats = append(diskStats, Stats{
			Name:            name,
			ReadsCompleted:  readsCompleted,
			WritesCompleted: writesCompleted,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for /proc/diskstats: %s", err)
	}
	return diskStats, nil
}

type StatUsage struct {
	Name  string
	Total int
	Used  int
	Mount string
}

/*
Filesystem     Type     1K-blocks    Used Available Use% Mounted on
udev           devtmpfs   1976444       0   1976444   0% /dev
tmpfs          tmpfs       401604    5924    395680   2% /run
/dev/sda1      ext4     205373416 6239964 188631412   4% /
tmpfs          tmpfs      2008012       0   2008012   0% /dev/shm
tmpfs          tmpfs         5120       0      5120   0% /run/lock
tmpfs          tmpfs      2008012       0   2008012   0% /sys/fs/cgroup
tmpfs          tmpfs       401600       0    401600   0% /run/user/0
*/

func lineToList(line string) []string {
	var output []string
	list := strings.Split(line," ")
	for _, v := range list {
		if len(v) == 0 {
			continue
		}
		output = append(output, v)
	}
	return output
}

func GetDiskUsage() ([]StatUsage, error) {
	var stat []StatUsage
	cmd := proc.NewCmd("df","-T")
	code := cmd.Run()
	if code != 0 {
		return nil, errors.New("run cmd [/bin/bash df] fail" + cmd.Stderr())
	}
	output := cmd.Stdout()
	lines := strings.Split(output, "\n")
	for _, v := range lines {
		if -1 != strings.Index(v, "tmpfs") {
			continue
		}
		if -1 != strings.Index(v, "Filesystem") {
			continue
		}
		list := lineToList(v)
		if len(list) == 0 {
			continue
		}
		stat = append(stat, StatUsage{
			Name: list[0],
			Total: util.StringToInt(list[2]),
			Used: util.StringToInt(list[3]),
			Mount: list[6],
		})
	}
	return stat, nil
}