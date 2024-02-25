## geolog - watches a web server log file and draw geolocation pings on a map in your browser

> [!IMPORTANT]
> You have to sign up to get a Maxmind licence key for their free geoip service at: https://dev.maxmind.com/geoip before you can use this tool.

### Run the Go program from it's github repository:
```
go run github.com/vindolin/geolog@latest -l /var/log/nginx/access.log -g /opt/GeoLite2-City.mmdb -d
```


### Run the Go program from source:
```
go run . -l /var/log/nginx/access.log -g /opt/GeoLite2-City.mmdb -p 8080
```

### Build the docker container:

```
docker build --build-arg ACCOUNT_ID=976666 --build-arg LICENSE_KEY=r872TfXFxdEJjBGgvkCVwU6zDY4Au3WMK5RqmNts -t vindolin/geolog .
```

### ... and run it:

#### in light mode
```
docker run --name geolog -it --rm -e LOG_FILE=/var/log/nginx/brummellock_access.log -v /var/log/nginx:/var/log/nginx:ro -p 8080:80 vindolin/geolog
```

#### in dark mode
```
docker run --name geolog -it --rm -e LOG_FILE=/var/log/nginx/brummellock_access.log -e DARK_MODE=true -v /var/log/nginx:/var/log/nginx:ro -p 8080:80 vindolin/geolog
```
