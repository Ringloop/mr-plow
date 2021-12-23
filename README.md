![example workflow](https://github.com/Ringloop/mr-plow/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/Ringloop/mr-plow/branch/master/graph/badge.svg)](https://codecov.io/gh/Ringloop/mr-plow)

# Mr-Plow
Tiny and minimal tool to export data from relational db (postgres or mysql) to elasticsearch.

The tool does not implement all the logstash features, but its goal is to be an alternative to logstash when keeping in-sync elastic and a relational database.

### Goal
**Low memory usage**: (~15 MB when idle, great to be deployed on cloud environments).

**Stateless**: a timestamp/date column is used in order to filter inserted/update data and to avoid fetching already seen data.
During the startup Mr-Plow checks the data inserted into elasticsearch to check the last timestamp/date of the transferred data, and so it does not require a local state.

![image](https://user-images.githubusercontent.com/7256185/141697554-4e6f86d8-06e4-4c22-aea5-30145e40fc41.png )

### Usage:
Mr-Plow can execute many queries in parallel.
Specify a timestamp/date column in the queries in order to get only newly updated/inserted data.

Configuration template example:
```yaml
# example of config.yml
pollingSeconds: 5 #database polling interval
database: "postgres://user:pwd@localhost:5432/postgres?sslmode=disable" #specify here the db connection
queries: #put here one of more queries (each one will be executed in parallel):
  - query: "select * from my_table1 where last_update > $1" #please add a filter on an incrementing date/ts column using the $1 value as param
    index: "table1_index" #name of the elastic output index
    updateDate: "last_update" #name of the incrementing date column
    id: "customer_id" #optional, column to use as elasticsearch id
  - query: "select * from my_table2 where ts > $1"
    index: "table2_index"
    updateDate: "ts"
elastic:
  url: "http://localhost:9200"
  user: "elastic_user" #optional
  password: "my_secret" #optional
  numWorker: 10 #optional, number of worker for indexing each query
  caCertPath: "my/path/ca" #optional, path of custom CA file (it may be needed in some HTTPS connection..)
```

Download or build the binary (docker images will be released soon):
```bash
go build
```

Run the tool:
```bash
./mr-plow -config /path/to/my/config.yml
```


