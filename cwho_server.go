//
// Cwho data collection server, Sean Caron, scaron@umich.edu
//

package main

import (
    // "io"
    "net"
    // "os"
    "fmt"
    //"strings"
    "bufio"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

    listener, err := net.Listen("tcp", "localhost:5963")
    if err != nil {
        return
    }
    
    for {
        conn, err := listener.Accept()
	if err != nil {
            continue
	}
	
	go handle_connection(conn)
    }
}

//
// Database schema:
//  CREATE TABLE ENTRIES (host varchar(255), user varchar(255), line varchar(255), fromhost varchar(255), timestamp varchar(255));
//

func handle_connection(c net.Conn) {

    var dbUser string = "cwho"
    var dbPass string = "xyzzy123"
    var dbName string = "cwho"
    var dbHost string = "localhost"

    var myDSN string;
    
    input := bufio.NewScanner(c)
    
    fmt.Printf("%s\n", input.Text())
    
    for input.Scan() {
    
        inp := input.Text()
	
        //data := strings.Fields(inp)

        //fmt.Printf("%s\n", input.Text())

        fmt.Printf(inp)
 
        myDSN = dbUser + ":" + dbPass + "@" + dbHost + "/" + dbName
    
        dbconn, _ := sql.Open("mysql", myDSN)

	dbconn.Close()
    }
    
    c.Close()
}
