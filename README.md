CWho
----
Sean Caron (scaron@umich.edu)

Centralized who with data inserted to a database.

Web dashboard shows aggregate who results.

Figure out who is logged into which machine over a large pool of machines.

Schema for utmp table:

```
CREATE TABLE utmp (sampletime bigint, host varchar(258), user varchar(34),
  line varchar(34), fromhost varchar(258), timestamp varchar(34));
```

Schema for hosts table:

```
CREATE TABLE hosts (host varchar(258), hostid integer NOT NULL
 AUTO_INCREMENT PRIMARY KEY, mostrecent bigint);
```
