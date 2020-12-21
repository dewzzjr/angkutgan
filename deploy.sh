#!/bin/bash

file=`ls deploy | grep ".sql"`

# Read Username
echo -n "Host: "
read host

# Read Username
echo -n "DB Name: "
read dbname

# Read Username
echo -n "Username: "
read username

# Read Password
echo -n "Password: "
read -s password


for i in $file
do
    `mysql --host=$host --user=$username --password=$password --database=$dbname --execute=./deploy/$i`
done

