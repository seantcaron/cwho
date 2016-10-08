#!/usr/bin/python

# Pull data from the CWho database and generate the Web dashboard
#  Sean Caron (scaron@umich.edu)

import cgi, time, sys, MySQLdb

print('Content-type: text/html\n')
print('<html>')
print('<head>')
print('<title>CWho</title>')
print('<style type="text/css">h1 { font-family: Arial, Helvetica; } p { font-size: medium; font-weight: bold; font-family: Arial, Helvetica; width: 80%; margin: 10px auto; } table { height: 15%; margin: 10px auto; width: 80%; } td { font-family: Courier; }</style>')
print('</head>')
print('<body bgcolor=White text=Black vlink=Black text=Black>')
print('<h1>CWho: ' + time.strftime("%A %b %d %H:%m:%S %Z", time.localtime()) + '</h1>')

db = MySQLdb.connect(user="cwho",passwd="xyzzy123",db="cwho")

curs = db.cursor()

query = 'SELECT host from hosts;'
curs.execute(query)
hosts = curs.fetchall()

for host in hosts:
    query = 'SELECT * FROM utmp WHERE host = \'' + host[0] + '\' and latest = true;'

    curs.execute(query)

    utmps = curs.fetchall()

    toggle = 0

    # user port fromhost time

    print('<p>' + host[0] + '</p>')
    print('<table>')
    for row in utmps:
        if toggle == 0:
            print('<tr bgcolor=#ccffcc><td>')
        else:
            print('<tr><td>')
    
        print(row[1])
        print('</td><td>')
        print(row[2])
        print('</td><td>')
        print(row[4])
        print('</td><td>')
        print(row[3])
        print('</td></tr>')
 
        toggle = not toggle

    print('</table>')

# We need to commit() the query on inserts and modifies after execution before they actually take effect
# db.commit()

print('</body>')
print('</html>')

db.close()
