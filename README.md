CWho
----
Sean Caron (scaron@umich.edu)

CWho centralizes the gathering of "who" data across a cluster of machines.
The data is recorded internally as a time series, which facilitates further
analytics.

A lightweight agent written in Go reads the utmp file directly on each client,
parses it and sends the data to a centralized collection server. The agent is
intended to run out of cron, at any frequency desired by the user.

A collection server also written in Go accepts connections from clients and
records each line of utmp data with a timestamp generated when the connection
is initiated.

A web dashboard written in Python shows the aggregate results recorded by the
server.

CWho permits rapid determination of whether or not a particular user is or is
not logged into any given machine at any time.

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

Schema for last table:

```
CREATE TABLE last (host varchar(258), user varchar(34), timestamp varchar(34));
```

