package main

import (
    "fmt"
    "time"

    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/mem"
)

func main() {
    for {
        // Fetch memory usage
        v, _ := mem.VirtualMemory()

        // Fetch CPU usage (per core)
        cpuPercents, _ := cpu.Percent(0, true)

        // Print memory usage
        fmt.Printf("Memory - Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

        // Print CPU usage per core
        for i, percent := range cpuPercents {
            fmt.Printf("CPU Core %d: %f%%\n", i, percent)
        }

        // Sleep for 0.5 seconds
        time.Sleep(500 * time.Millisecond)

        fmt.Println("-------")
    }
}

/// run:  jim-emacs-fun-go  master @ go run while_cpu_mem.go
//Memory - Total: 19327352832, Free:212729856, UsedPercent:63.236491%
//CPU Core 0: 18.367347%
//CPU Core 1: 20.408163%
//CPU Core 2: 16.000000%
//CPU Core 3: 14.000000%
//CPU Core 4: 8.163265%
//CPU Core 5: 6.000000%
//CPU Core 6: 0.000000%
//CPU Core 7: 3.921569%
//CPU Core 8: 0.000000%
//CPU Core 9: 0.000000%
//CPU Core 10: 2.040816%
//
