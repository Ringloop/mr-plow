![example workflow](https://github.com/Ringloop/mr-plow/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/Ringloop/mr-plow/branch/main/graph/badge.svg?token=PE53PJ8HHR)](https://codecov.io/gh/Ringloop/mr-plow)

# Mr-Plow
Tiny and minimal tool to export data from relational db (postgres or mysql) to elasticsearch.

The tool does not implement all the logstash features, but its goal is to be an alternative to logstash when keeping in-sync elastic and a relational database.

### Goal
**Low memory usage**: (~15 MB when idle, great to be deployed on cloud environments).

**Stateless**: a timestamp/date column is used in order to filter inserted/update data and to avoid fetching already seen data.
During the startup Mr-Plow checks the data inserted into elasticsearch to check the last timestamp/date of the transferred data, and so it does not require a local state.

![image](https://user-images.githubusercontent.com/7256185/141697554-4e6f86d8-06e4-4c22-aea5-30145e40fc41.png )

### Usage:
Mr-Plow essentially executes queries on a relational database and writes these data to ElasticSearch.
The configured queries are run in parallel, and data are written incrementally, it's only sufficient to specify a timestamp/date column in the queries in order to get only newly updated/inserted data.

This is a basic configuration template example, where we only specify two queries and the endpoint configuration (one Postgres database and one ElasticSearch cluster:
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

There are some additional features Mr Plow can manage, such as specifying the expected data type to be inserted in ElasticSearch JSON (not specified fields are parsed as String):

```yaml
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from employees
    updateDate: last_update
    fields:
      - name: name
        type: String
      - name: working_hours
        type: Integer
```

And additionally (for database like Postgres) you can also specify JSON columns fields:
```yaml
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from employees
    updateDate: last_update
    fields:
      - name: name
        type: String
      - name: working_hours
        type: Integer
    JSONFields:
      - fieldName: dataField_1
        fields:
          - name: attribute_1_Name
            type: attribute_1_Type
      - fieldName: dataField_2
        fields:
          - name: attribute_2_Name
            type: attribute_2_Type
    id: MyId_1
```


Download or build the binary (docker images will be released soon):
```bash
go build
```

Run the tool:
```bash
./mr-plow -config /path/to/my/config.yml
```

To build as docker image, create a `config.yml` and put into the root folder of the project. Then run:
```bash
docker build .
```


