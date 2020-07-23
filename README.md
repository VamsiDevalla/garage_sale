# Garage Sales Web Service

This web will be used to experiment, practice and develop my own design philosophy for **go lang**
This application has two binaries:

1. **Sales Admin**
This binary can be used to create/migrate tables in the database and to seed the database
1. **Sales API**
This binary can be used to run curd operations against the garage sale items

## Table of contents

- [Application Checklist](#application-checklist)
- [Commands](#commands)
- [KEY POINTS / GUIDE LINES](#key-points--guide-lines)

## Application Checklist

- [x] Start and shutdown the serve gracefully
- [x] Encode/decode json
- [x] Connect to a database for permanent storage
- [x] Package oriented design
- [x] Reading config from external source
- [ ] Logging
- [ ] Error Handling
- [ ] Testing and benchmarking
- [ ] Tracing and profiling  

## Commands

### Project Commands

1. Run Postgres database

   ```shell
   docker-compose up
   ```

1. Run Sales admin
   - Migrate

     ```shell
     go run ./cmd/sales-admin migrate
     ```

   - Seed

     ```shell
     go run ./cmd/sales-admin seed
     ```

1. Run Sales Api

   ```shell
   go run ./cmd/sales-api
   ```

### General Commands / Tools

## KEY POINTS / GUIDE LINES

- Always shutdown the server gracefully
- Init functions:
  - They run during the package initialization.
  - All the init functions from the imported package are ran before the main function.
  - If you want to initialize a package that is not required in the current package then use '_' to ignore the import but run the init functions.
  
  ```go
  import (
      _ "github.com/lib/pq" // to register postgres driver
  )
  ```

- Package Oriented design:
  - If the application has more than one binary place them in the into separate main files in their own folders in cmd folder.
  - Make use of internal folder/directory so that only packages at its parent level can import them. This is help full when you want created a library and don't want the clients to use the under lying packages.
  - **Read about package oriented design [here](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)**

- When choosing a config package consider:
  - where you wanted to read your config from?
    - JSON / YML
    - ENV Variables
    - Command line Arguments / flags
  - do you want the ability to update the config with out restarting the server.

- Caution when using `log.Fatal`  
  - log.Fatal logs the error and calls the `os.Exit(1)` function
  - Because of this `defer functions` wouldn't be triggered
  - To main the integrity we always let the defer functions to run
  - Instead we could just return the error from the function and let the calling function take care of the exit.

  ```go
  package main
  
  import (
    "log"
    "os"
    "github.com/pkg/errors"
  )
    // in the following defer functions won't be executed
  /*
  func main() {
    defer close()
    log.Fatal("error: shutting down", err)
  }
  */

   // defer functions inside run will be executed even then there is a fatal error
  func main() {
    // defer close(some file/ http response)
    if err := run(); err != nil {
      log.Print("error: shutting down", err)
      os.Exit(1)
    }
  }
  
  func run() error {
    if err := error.New("Test error"); err!= nil {
      return errors.Wrap(err, "error while running running")
    }
    return nil
  }
  ```
