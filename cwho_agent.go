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
    "time"
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
    // For each line of the utmp file, parse out the information that we need.
    //

    fmt.Printf("who results for %s\n", host)

    for _, arg := range ut {
        if (arg.Type == 7) {
            tt := time.Unix(int64(arg.Tv.Sec), int64(arg.Tv.Usec))
            ts := tt.Format(time.ANSIC)
            fmt.Printf("%s\t%s\t%s\t%s\n", arg.User, arg.Line, arg.Host, ts)
        }
    }

    conn, err := net.Dial("tcp", "localhost:5962")
    if err != nil {
        return
    }
    
    //fmt.Fprintf(conn, "", )
    
    conn.Close()
}

