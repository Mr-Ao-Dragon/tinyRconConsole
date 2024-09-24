# settings
use: nutsDB
## Bucket: settings
| key            | valueType | interpretation                   |
|----------------|-----------|----------------------------------|
| storageBackend | string    | define source of data storage    |
| cryptoBackend  | string    | define source of data encryption |
| tcpKeepLive    | boolean   | enable tcp keepalive             |
| ttyMode        | boolean   | set tty mode                     |
| timezone       | integer   | set timezone                     |
| language       | string    | set language                     |
| privateKey     | string    | storage private key              |
# data
use: nutsDB
## Bucket: servers
| sets    | description         |
|---------|---------------------|
| servers | storage server list |
servers key map:

| name       | datatype | default value    | description      | can be empty |
|------------|----------|------------------|------------------|--------------|
| serverName | string   | minecraft server |                  | true         |
| address    | string   | localhost        | server address   | false        |
| port       | integer  | 25575            | serverPort       | false        |
| password   | string   |                  | is encrypted     | true         |
| nanoID     | string   |                  | only have one id | false        |
## bucket: savedCommand
| lists   | description     |
|---------|-----------------|
| history | storage history |
| micros  | storage micros  |
history key map:

| key name     | type    | description               |
|--------------|---------|---------------------------|
| command      | string  | executed command          |
| time         | integer | time of executed          |
| execServerID | string  | history executed serverID |
struct:
index sigID:
- command
- time
- execServerID
micros key map:

| key name | type   | description      |
|----------|--------|------------------|
| command  | string | recorded command |
struct:
index: microID
- command