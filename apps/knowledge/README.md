# Knoweldge
This provides data access and storage for the Hub UI. It provides a number of REST endpoints that can be called.

/ - HealthCheck on the service <br/>
/record - Allows CRUD operations on records<br/>
/field - Allows CRUD operations for fields that specify what data can be seen and entered.<br/>

## Build
Dep ensure wants to have its path as $GOPATH/src/{project} so need to create a symbolic link from there to where the code actually is to get it to run.

## Run
Set the following environment variables:<br/>
PORT - To run server on (defaults to 8000, if using it built by the Dockerfile it's always 8000)<br/>
MONGODB_URI - in the format mongo://<\user>:<pass/token>@<\server><br/>
(Optional)<br/>
LOG_LEVEL - Level to log at, debug or info (defaults to info)<br/>