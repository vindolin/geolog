## geolog - watches a web server log file and draw geolocation pings on a map in your browser

![image](https://github.com/vindolin/geolog/assets/834347/0bad8918-99b6-45b4-a88a-e58eca3b3c67)


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
docker build --build-arg ACCOUNT_ID={YOUR_ACCOUNT_ID} --build-arg LICENSE_KEY={YOUR_ICENSE_KEY} -t vindolin/geolog .
```

### ... and run it:

#### in light mode
```
docker run --name geolog -it --rm -e LOG_FILE=/var/log/nginx/access.log -v /var/log/nginx:/var/log/nginx:ro -p 8080:80 vindolin/geolog
```

#### in dark mode
```
docker run --name geolog -it --rm -e LOG_FILE=/var/log/nginx/access.log -e DARK_MODE=true -v /var/log/nginx:/var/log/nginx:ro -p 8080:80 vindolin/geolog
```
