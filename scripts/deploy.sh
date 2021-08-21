#!/bin/bash
echo "DB_URL"
read url
[[ -z "$url" ]] && exit 1

url=${url#"mysql://"}
IFS=':' read -r -a array <<< "$url"
username="${array[0]}"
url="${array[1]}"
portname="${array[2]}"
# echo $username

IFS='@' read -r -a array <<< "$url"
password="${array[0]}"
host="${array[1]}"
# echo $password
# echo $host

IFS='/' read -r -a array <<< "$portname"
port="${array[0]}"
dbname="${array[1]}"
# echo $port
# echo $dbname

file=`ls queries | grep ".sql"`
for i in $file
do
    echo $i
    `mysql --host=$host --user=$username -p$password --database=$dbname < "./queries/$i"`
done