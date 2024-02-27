# geolog
## Visualize live web server access on a browser map.

![image](https://raw.githubusercontent.com/vindolin/geolog/main/screencast.webp)


```shell

usage: run [-h|--help] -l|--log_file "<value>" -g|--geodb_file "<value>"
           [-p|--port "<value>"] [-s|--blip_size <float>] [-b|--blip_life_time
           <integer>] [-d|--dark]

           run the geolog websocket server

Arguments:

  -h  --help            Print help information
  -l  --log_file        log file to tail
  -g  --geodb_file      geolite db to use
  -p  --port            port to listen on. Default: 8080
  -s  --blip_size       Maximum size of the blib relative to the map width. Default: 0.3
  -b  --blip_life_time  life time of the map blips (Milliseconds). Default: 2000
  -d  --dark            dark mode
  ```

> [!IMPORTANT]
> Geolog needs a local *.mmdb Maxmind database to lookup the geo location.
>
> You can get a free licence to download their geoip databases at: https://dev.maxmind.com/geoip.
>
> Geolog creates a http server from which the map is served, updates happen over websockets.
>
> Websockets get automatically reconnected.
>


### Run the Go program directly from github:
```shell
go run github.com/vindolin/geolog@latest \
-l /var/log/nginx/access.log \
-g /opt/GeoLite2-City.mmdb \
-d
```


### ... or from source:
```shell
go run . \
-l /var/log/nginx/access.log \
-g /opt/GeoLite2-City.mmdb \
-p 8080
```

## Docker

### Build the container:

```shell
docker build -t vindolin/geolog \
--build-arg ACCOUNT_ID={YOUR_ACCOUNT_ID} \
--build-arg LICENSE_KEY={YOUR_ICENSE_KEY} \
.
```

### ... and run it:

#### light map
![image](https://raw.githubusercontent.com/vindolin/geolog/main/light.png)

```shell
docker run -it --rm --name geolog \
-v /var/log/nginx:/var/log/nginx:ro \
-p 8080:80 \
vindolin/geolog \
-l /var/log/nginx/access.log
```

#### dark map
![image](https://raw.githubusercontent.com/vindolin/geolog/main/dark.png)

```shell
docker run -it --rm --name geolog \
-v /var/log/nginx:/var/log/nginx:ro \
-p 8080:80 \
vindolin/geolog \
-l /var/log/nginx/access.log \
-d
```
