#!/usr/bin/python

# Pull data from the CWho database and generate the Web dashboard
#  Sean Caron (scaron@umich.edu)

import cgi, sys, MySQLdb

print('Content-type: text/html\n')
print('<HTML><HEAD><TITLE>CWho</TITLE></HEAD>\n')
print('<BODY BGCOLOR=white TEXT=BLACK VLINK=BLACK TEXT=BLACK>\n')
print('<STYLE>h1 { font-family: Arial, Helvetica; } p { font-family: Arial, Helvetica; } tr { font-family: Courier; }</STYLE>\n')
print('<H1>CWho</H1>\n')

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

    print('<P><B>' + host[0] + '</B><P>\n')
    print('<TABLE CELLSPACING=0 CELLPADDING=10 BORDER=0>\n')
    for row in utmps:
        if toggle == 0:
            print('<TR BGCOLOR=#CCFFCC><TD>\n')
        else:
            print('<TR><TD>\n')
    
        print(row[1])
        print('</TD><TD>')
        print(row[2])
        print('</TD><TD>')
        print(row[3])
        print('</TD><TD>')
        print(row[4])
        print('</TD></TR>\n')
  
        toggle = not toggle

    print('</TABLE>\n')

    print('<P>\n')
# We need to commit() the query on inserts and modifies after execution before they actually take effect
# db.commit()

print('</BODY></HTML>\n')

db.close()
