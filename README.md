
Project tp extract data from excle file and and show data in API and perform operations like delete update and  getinfo


## we need  this servers running 

- MySQL : have to crete new data base and use it and create new table call user insert values of all the values from excle file 
- Redis : get the data from database and and clear it after 5 min as time setup it make data easy to access
-we need this package to run program  Go packages: gin, excelize, go-sql-driver/mysql, go-redis


1. **Create MySQL Database and Table**
   Execute the SQL commands provided above to set up the database and table.

2. **Update Database Connection**
   Modify the DSN in `initDB` function to connect to your MySQL database.

## also get this packages to run the program
go get -u github.com/gin-gonic/gin
go get -u github.com/xuri/excelize/v2
go get -u github.com/go-redis/redis/v8


## for database access create new db and change it in a code accordigly

## working fine