# Knoweldge
This provides data access and storage for the Hub UI. It provides a number of REST endpoints that can be called.

/ - HealthCheck on the service <br/>
/record - Allows CRUD operations on records<br/>
/field - Allows CRUD operations for fields that specify what data can be seen and entered.<br/>

## Build
Dep ensure wants to have its path as $GOPATH/src/{project} so need to create a symbolic link from there to where the code actually is to get it to run.

### Build the server
The following command will build the server, it produces an executable called `main`
Note: It requires go v1.11+ to make use of modules.
```
    go build ./cmd/server/main.go
```

### Docker build
Alternatively the app can be build into a docker image and the deployed.
```
    docker build -t knowledge:<version> .
```

## Test
There are a number of unit tests that can be run, use the following command to run all tests:
```
    go test ./...
```

## Run
### Local
Set the following environment variables:<br/>
PORT - To run server on (defaults to 8000, if using it built by the Dockerfile it's always 8000)<br/>
MONGODB_URI - in the format mongo://<\user>:<pass/token>@<\server><br/>
(Optional)<br/>
LOG_LEVEL - Level to log at, debug or info (defaults to info)<br/>

And then run 
```
./main
```

### Docker
Set the following environment variables:<br/>
        MONGODB_URI - URL of the mongo DB to connect to<br/>
        (Optional) <br/>
        LOG_LEVEL - info or debug<br/>
Then start the container. It runs listening on port 8000 within the container. It must be able to connect to the Mongo one<br/>
    