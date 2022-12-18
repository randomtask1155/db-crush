

## SETUP

build a databse with a table to get things started

```
CREATE DATABASE crushdb;
```

```
use crushdb;
```

```
CREATE TABLE IF NOT EXISTS crushit (code varchar(255), 
    authentication BLOB, 
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    expiresat BIGINT DEFAULT 0 NOT NULL, 
    user_id VARCHAR(36) NULL, 
    client_id VARCHAR(36) NULL, 
    `id` int(11) unsigned PRIMARY KEY AUTO_INCREMENT,
    identity_zone_id VARCHAR(36) NULL);
```



## USAGE

```
export MYSQL_HOST=IP-ADDRESS MYSQL_PWD=PASSWORD MYSQL_USER=USER
nohup /root/db-crush > /dev/null &
```