# internal

This directory has a special meaning in Go. The code inside it can only be imported and reused by the code in its parent directory. In other words, the packages inside `internal` can only be imported within the `snippetbox` project directory even if it is hosted on Github. All non-application-specific code that can be reused across multiple applications goes here

# PostgreSQL database

Use the SQL file in `internal` to initialize the database necessary for this application
