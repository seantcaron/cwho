//
// Cwho data collection server, Sean Caron, scaron@umich.edu
//

package main

import (
    // "io"
    "net"
    // "os"
    "fmt"
    "strings"
    "bufio"
    "log"
    "strconv"
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
//  CREATE TABLE entries (host varchar(258), user varchar(34), line varchar(34), fromhost varchar(258), timestamp varchar(34));
//  CREATE TABLE hosts (host varchar(258), hostid integer NOT NULL AUTO_INCREMENT PRIMARY KEY);
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
	
        data := strings.Fields(inp)

	host := data[0]
	user := data[1]
	line := data[2]
	from := data[3]
	stamp := data[4]

        //fmt.Printf("%s\n", input.Text())

        fmt.Printf(inp+"\n")
 
        myDSN = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":3306)/" + dbName
    
        dbconn, dbConnErr := sql.Open("mysql", myDSN)
        if dbConnErr != nil {
	    log.Fatalf("Error connecting to database")
        }

	dbPingErr := dbconn.Ping()
	if dbPingErr != nil {
	    log.Fatalf("Error attempting to ping database connection")
        }

	//
	// Check to see if the host exists in the host tracking table
	//

        dbCmd := "SELECT COUNT(*) FROM hosts WHERE host = '" + host + "';"
	_, dbExecErr := dbconn.Exec(dbCmd)
	if dbExecErr != nil {
	    log.Fatalf("Failure executing SELECT for host " + host)
        }

	var hostp string;
	_ = dbconn.QueryRow(dbCmd).Scan(&hostp)
	hostpi, _ := strconv.Atoi(hostp)

	//
	// If not, add it to the hosts table. MySQL will generate an ID
	//

	if (hostpi == 0) {
            dbCmd := "INSERT INTO hosts (host) VALUES ('" + host + "');"
	    _, dbExecErr = dbconn.Exec(dbCmd)
	    if dbExecErr != nil {
	        log.Fatalf("Failure executing INSERT for host " + host)
            }
        }

	//
	// Add the most recent batch of utmp entries to the database
	//

        dbCmd = "INSERT INTO entries VALUES ('" + host + "','" + user + "','" + line + "','" + from + "','" + stamp + "');"
        _, dbExecErr = dbconn.Exec(dbCmd)
	if dbExecErr != nil {
	    log.Fatalf("Failure executing INSERT for host " + host)
        }

	dbconn.Close()
    }
    
    c.Close()
}
