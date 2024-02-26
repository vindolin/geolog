# geolog
## Visualize web server access on a browser map.

![image](https://raw.githubusercontent.com/vindolin/geolog/main/screencast.webp)



> [!IMPORTANT]
> Geolog needs a local *.mmdb Maxmind database to lookup the geo location.
>
> You can get a free licence to download their geoip databases at: https://dev.maxmind.com/geoip.


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
