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

Anyway, Mr Plow has also additional features, for example interacting with a database like Postgres, supporting JSON columns, we can specify JSON fields, in order to create a complex (nested) object to be created in Elastic. In the following, we show an example where in the Employee table we store two dynamic JSON fields, one containing the Payment Data and another one containing additional informations for the employee:

```yaml
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from employees
    updateDate: last_update
    JSONFields:
      - fieldName: payment_data
      - fieldName: additional_infos
    id: MyId_1
```

And additionally, we can specify the type expected for some specific fields. Please note hat field type is optional and if not specified, the field is casted as String.

Actually supported type are: String, Integer, Float and Boolean

```yaml
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from employees
    updateDate: last_update
    fields: # Optional config, casting standard sql columns to specific data type
      - name: name
        type: String
      - name: working_hours
        type: Integer
```
Merging the previous two examples, we can apply the type casting also to inner JSON fields, here is a complete example of configuration:

```yaml
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from employees
    updateDate: last_update
    fields: # Optional, casting standard sql columns to specific data type
      - name: name
        type: String
      - name: working_hours
        type: Integer
    JSONFields:
      - fieldName: payment_info
        fields: # Optional, casting json fields to specific data type
          - name: bank_account
            type: String
          - name: validated
            type: Boolean
    id: MyId_1
```

Download or build the binary (docker images will be released soon):
```bash
make 
```

Run the tool:
```bash
./bin/mr-plow -config /path/to/my/config.yml
```

To build as docker image, create a `config.yml` and put into the root folder of the project. Then run:
```bash
docker build .
```

### Mr-Plow development with Visual Studio Code

**Requirements**

- [Docker](https://docs.docker.com/)
- [Visual Studio Code](https://code.visualstudio.com/)

Linux and Mac users should ensure to have the following programs installed on their system:  

- `bash`, `id`, `getent`

Windows users must be aware that they should clone Mr-Plow repository over the `WSL` filesystem.  
This is the recommended way because mounting a `NTFS` filesystem inside a container exposes the overall user's experience to major issues.  
Windows users should also patch the `.devcontainer/devcontainer.json` as indicated in the comments inside the file.  

**Steps**

1. Clone the project to a local folder
2. VSCode -> `><` (left bottom button) -> Open Folder in Container...

**Instructions**

Users can develop and test Mr-Plow inside a docker container without being forced to install or configure anything on your own machine.  
Visual Studio Code can take care of automatically download and build the developer's docker image.  
For Linux and Mac users, an especial care has been devoted to make sure the host's user will match `UID` and `GID` with the user inside the container.  
This ensures that every modification from inside the container will be completely transparent from the host's perspective.  
Moreover, host's user `~.ssh` directory will be mounted on the container's user `~.ssh` directory. This is especially convenient if an ssh authentication type is configured to work with GitHub.  
From inside the container, users will be able to access the host's docker engine as if they were just in a regular host's shell.  
This capability allows users to launch the predefined `docker-compose` images directly from Visual Studio Code.  
Users can simply access to the task menu pressing: `ctrl + shift + p` and select `Docker: Compose Up`.  
Therefore, they can choose to spawn up the following services:

- `docker-compose-elastic.yml`: ElasticSearch
- `docker-compose-kibana.yml`: Kibana
- `docker-compose-postgres.yml`: Postgres

