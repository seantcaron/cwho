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

    for _, arg := range ut {
        fmt.Printf("%s %s\n", arg.User, arg.Host)
        //fmt.Printf("%s %s %s %s %s\n", arg[0], arg[1], arg[2], arg[4], arg[5])
        //fmt.Printf("%s\n", ut[arg].User)
    }

    conn, err := net.Dial("tcp", "localhost:5962")
    if err != nil {
        return
    }
    
    //fmt.Fprintf(conn, "", )
    
    conn.Close()
}

