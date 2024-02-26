# geolog
## Visualize web server access on a browser map.

![image](https://raw.githubusercontent.com/vindolin/geolog/main/screencast.webp)



> [!IMPORTANT]
> You have to sign up to get a Maxmind licence key for their free geoip service at: https://dev.maxmind.com/geoip before you can use this tool.

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

#### light mode

```shell
docker run  -it --rm --name geolog \
-v /var/log/nginx:/var/log/nginx:ro \
-p 8080:80 \
vindolin/geolog \
-l /var/log/nginx/access.log
```

#### dark mode

```shell
docker run  -it --rm --name geolog \
-v /var/log/nginx:/var/log/nginx:ro \
-p 8080:80 \
vindolin/geolog \
-l /var/log/nginx/access.log \
-d
```
