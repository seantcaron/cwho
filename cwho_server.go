//
// Cwho data collection server, Sean Caron, scaron@umich.edu
//

package main

import (
    "net"
    "os"
    "fmt"
    "strings"
    "bufio"
    "log"
    "strconv"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var dbUser, dbPass, dbName, dbHost string

func main() {

    var bindaddr, conffile string

    if (len(os.Args) != 5) {
        log.Fatalf("Usage: %s -b bindaddr -f configfile\n", os.Args[0])
    }

    for i := 1; i < len(os.Args); i++ {
        switch os.Args[i] {
	    case "-b":
	        bindaddr = os.Args[i+1]
            case "-f":
	        conffile = os.Args[i+1]
        }
    }

    conf, err := os.Open(conffile)
    if err != nil {
        log.Fatalf("Failed opening configuration file for reading")
    }

    inp := bufio.NewScanner(conf)

    for inp.Scan() {
        line := inp.Text()

        if (len(line) > 0) {
	    fields := strings.Fields(line)
	    key := strings.ToLower(fields[0])

	    switch key {
                case "dbuser":
	            dbUser = fields[1]
                case "dbpass":
	            dbPass = fields[1]
                case "dbname":
	            dbName = fields[1]
                case "dbhost":
	            dbHost = fields[1]
                default:
	            log.Print("Ignoring nonsense configuration %s\n", fields[1])
            } 
        }
    }

    conf.Close()

    listener, err := net.Listen("tcp", bindaddr+":5963")
    if err != nil {
        log.Fatalf("Failure calling net.Listen()\n")
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
//  CREATE TABLE utmp ([sampletime bigint], host varchar(258), user varchar(34), line varchar(34), fromhost varchar(258), timestamp varchar(34), [XXXlatest booleanXXX]);
//  CREATE TABLE hosts (host varchar(258), hostid integer NOT NULL AUTO_INCREMENT PRIMARY KEY, [mostrecent bigint]);
//

func handle_connection(c net.Conn) {

    var myDSN string

    // var flip bool = true

    input := bufio.NewScanner(c)
    
    fmt.Printf("%s\n", input.Text())

    //
    // Generate a timestamp for these samples
    //

    t := time.Now().Unix()
    tt := strconv.FormatInt(t, 10)

    for input.Scan() {
    
        inp := input.Text()
	
        data := strings.Fields(inp)

	host := data[0]
	user := data[1]
	line := data[2]
	from := data[3]
	secs, _ := strconv.ParseInt(data[4], 10, 64)
        usecs, _ := strconv.ParseInt(data[5], 10, 64)

        fmt.Printf(inp+"\n")
 
        myDSN = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":3306)/" + dbName
    
        dbconn, dbConnErr := sql.Open("mysql", myDSN)
        if dbConnErr != nil {
	    log.Fatalf("Failed connecting to database")
        }

	dbPingErr := dbconn.Ping()
	if dbPingErr != nil {
	    log.Fatalf("Failed pinging database connection")
        }

	//
	// Check to see if the host exists in the host tracking table
	//

        dbCmd := "SELECT COUNT(*) FROM hosts WHERE host = '" + host + "';"
	_, dbExecErr := dbconn.Exec(dbCmd)
	if dbExecErr != nil {
	    log.Fatalf("Failed executing SELECT for host " + host)
        }

	var hostp string
	_ = dbconn.QueryRow(dbCmd).Scan(&hostp)
	hostpi, _ := strconv.Atoi(hostp)

	//
	// If not, add it to the hosts table. MySQL will generate an ID
	// If so, we need to update the sample time in the mostrecent
	//  field
	//

	if (hostpi == 0) {
            dbCmd := "INSERT INTO hosts (host, mostrecent) VALUES ('" + host + "'," + tt + ");"
	    _, dbExecErr = dbconn.Exec(dbCmd)
	    if dbExecErr != nil {
	        log.Fatalf("Failed executing host table INSERT for host " + host)
            }
        } else {
            dbCmd := "UPDATE hosts SET mostrecent = " + tt + " WHERE host = '" + host + "';"
	    _, dbExecErr = dbconn.Exec(dbCmd)
            if dbExecErr != nil {
	        log.Fatalf("Failed executing host table mostrecent field UPDATE for host " + host)
            }
        }

        //
	// Mark any existing entries as not being the latest, for the
	//  display layer to use
	//

        //if (flip == true) {
        //    dbCmd = "UPDATE utmp SET latest = false WHERE host = '" + host +"';"
	//    _, dbExecErr = dbconn.Exec(dbCmd)
	//    if dbExecErr != nil {
        //        log.Fatalf("Failed executing UPDATE on latest for host " + host)
	//    }
	//    flip = false
        //}

	//
	// Add the most recent batch of utmp entries to the database
	//

        stamp := time.Unix(secs,usecs).Format(time.Stamp)

        dbCmd = "INSERT INTO utmp VALUES (" + tt + ",'" + host + "','" + user + "','" + line + "','" + from + "','" + stamp + "');"
	fmt.Printf("%s\n", dbCmd)
        _, dbExecErr = dbconn.Exec(dbCmd)
	if dbExecErr != nil {
	    log.Fatalf("Failed executing utmp table INSERT for host " + host)
        }

	dbconn.Close()
    }
    
    c.Close()
}
