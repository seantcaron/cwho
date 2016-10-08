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
    "bytes"
    "github.com/EricLagergren/go-gnulib/utmp"
)

func main() {
    var u = "/var/run/utmp"
    var ut []*utmp.Utmp

    if ((len(os.Args) != 3) || (os.Args[1] != "-h")) {
        log.Fatalf("Usage: %s -h server\n", os.Args[0])
    }

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
    // Open the connection to the collection host
    //

    conn, err := net.Dial("tcp", os.Args[2]+":5963")
    if err != nil {
        log.Fatalf("Error calling net.Dial()")
    }

    //
    // For each line of the utmp file, parse out the information that we need.
    //

    fmt.Printf("who results for %s\n", host)

    for _, arg := range ut {
        if (arg.Type == 7) {
            tt := time.Unix(int64(arg.Tv.Sec), int64(arg.Tv.Usec))
            ts := tt.Format(time.RFC3339)

	    au := bytes.Trim(arg.User[:], "\x00") 
            al := bytes.Trim(arg.Line[:], "\x00")
	    ah := bytes.Trim(arg.Host[:], "\x00")

            if (len(ah) == 0) {
	        ah = []byte("local")
            }

            fmt.Printf("%s\t%s\t%s\t%s\n", au, al, ah, ts)
            fmt.Fprintf(conn, "%s %s %s %s %s\n", host, au, al, ah, ts)
        }
    }

    conn.Close()
}

