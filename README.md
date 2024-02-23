##  - WORK IN PROGRESS -

Tail a log file, resolve the geo location through Maxminds GeoLite database, then publish to a leaflet map over websockets.


```
go run . -l /var/log/nginx/access.log -g /opt/GeoLite2-City.mmdb -p 8080
```
