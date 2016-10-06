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
    "log"
    "net"
    "fmt"
    "github.com/EricLagergren/go-gnulib/utmp"
)

func main() {
    var u = "/var/run/utmp"
    var ut []*utmp.Utmp

    host, _ := os.Hostname()
    
    if (strings.Index(host, ".") != -1) {
        host = host[0:strings.Index(host, ".")]
    }

    //
    // Read in utmp file
    //

    ut, err := utmp.ReadUtmp(u, 0x00)
    if err != nil {
        log.Fatalf("Error opening utmp file for reading")
    }

    //
    // For each line of the utmp file, do
    //

    for _, arg := range ut {
        if (arg.Type == 7) {
            fmt.Printf("%d %s %s %s %d %d\n", arg.Type, arg.User, arg.Line, arg.Host, arg.Tv.Sec, arg.Tv.Usec)
        }
    }

    conn, err := net.Dial("tcp", "localhost:5962")
    if err != nil {
        return
    }
    
    //fmt.Fprintf(conn, "", )
    
    conn.Close()
}

