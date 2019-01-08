

## Running
Set the following environment variables:<br/>
  APP_CONFIG_VERSION - the current version (defaults to 0.0.1)<br/>
  APP_CONFIG_VERSION_COLOUR - the colour for the version bar (defaults to #58af58 which is green)<br/>
  APP_CONFIG_MAP_PROVIDER - where to get the map tiles from e.g. http://{s}.tile.osm.org/{z}/{x}/{y}.png (defaults to this OSM URL)<br/>
  APP_CONFIG_API_URL - URL of the knowledge API service (defaults to what's in .env files)<br/>
  APP_SERVER_PORT - the port for the web app to listen on (defaults to 3000)


Then run:<br/>
```  
  npm run start
```

## Build
### Local
For a production build run: <br/>
```
  npm run build
```

### Docker
Alternatively it can be built into a docker image:
```
  docker build -t <tag> .
```

It exposes the app on port 3000