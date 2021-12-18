#!/bin/bash
docker run -itd -p 3306:3306 \
       -v $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/mysql-data:/var/lib/mysql \
       --name stree-mysql \
       -e MYSQL_ROOT_PASSWORD=mypassword \
       mysql:5.7