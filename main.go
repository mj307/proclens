package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// process data 
type procInfo struct {
	pid   int32
	name  string
	cpu   float64
	ramMB float64
	proc  *process.Process
}

func main() {
	// CLI flags
	cpuThreshold := flag.Float64("cpu", 0.0, "Minimum CPU % to filter")
	memThreshold := flag.Float64("mem", 0.0, "Minimum RAM usage (MB) to filter")
	dryRun := flag.Bool("dry-run", true, "Only show results (default: true)")
	kill := flag.Bool("kill", false, "Actually kill the matching processes")
	flag.Parse()

	// get the running processes
	procs, err := process.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
		return
	}

	fmt.Printf("Filtering for: CPU > %.1f%%, RAM > %.0fMB | Dry-run: %t | Kill: %t\n\n",
		*cpuThreshold, *memThreshold, *dryRun, *kill)
	fmt.Printf("%-8s %-40s %-8s %-8s\n", "PID", "Name", "CPU%", "RAM_MB")

	// go through each process
	var results []procInfo
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}
		pid := p.Pid 

		// initial CPU read (ignored, sets baseline)
		_, _ = p.CPUPercent()
		time.Sleep(200 * time.Millisecond) // wait to sample over time
		cpu, err := p.CPUPercent()
		if err != nil {
			continue 
		}

		// memory usage
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		ramMB := float64(memInfo.RSS) / 1024 / 1024

		// threshold filters
		if cpu < *cpuThreshold || ramMB < *memThreshold {
			continue
		}

		// results list
		results = append(results, procInfo{
			pid:   pid,
			name:  name,
			cpu:   cpu,
			ramMB: ramMB,
			proc:  p,
		})

		// sleep to prevent system overload
		time.Sleep(5 * time.Millisecond)
	}

	// sort results by CPU descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].cpu > results[j].cpu
	})

	// print + kill (if kill mode)
	for _, r := range results {
		fmt.Printf("%-8d %-40s %-8.1f %-8.0f\n", r.pid, r.name, r.cpu, r.ramMB)

		// kill mode (if flag is applied)
		if *kill && !*dryRun {
			err := r.proc.Kill()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to kill PID %d (%s): %v\n", r.pid, r.name, err)
			} else {
				fmt.Printf("Killed PID %d (%s)\n", r.pid, r.name)
			}
		}
	}
}

