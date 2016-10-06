//
// CWho agent, Sean Caron, scaron@umich.edu
//

package main

// go get github.com/EricLagergren/go-gnulib/utmp

import (
    //"bufio"
    //"fmt"
    "os"
    //"os/exec"
    "strings"
    //"strconv"
    //"time"
    "net"
    "github.com/EricLagergren/go-gnulib/utmp"
)

func main() {
    var u = "/var/run/utmp"
    var ut *utmp.Utmp
    var uf *utmp.File

    host, _ := os.Hostname()
    
    if (strings.Index(host, ".") != -1) {
        host = host[0:strings.Index(host, ".")]
    }

    uf, err := utmp.Open(u, Reading)
    if err != nil {
        os.Exit(1)
    }

    ut = utmp.GetUtEnt(u)
        
    conn, err := net.Dial("tcp", "localhost:5962")
    if err != nil {
        return
    }
    
    //fmt.Fprintf(conn, "%d,%s,%d,%d,%f,%f,%f,%f,%s\n", timestamp, host, nc, mt, om, fivm, fifm, swap_used_pct, diskReport)
    
    conn.Close()
}

