#!/bin/sh

# set the dark mode flag if the DARK_MODE environment variable is set to true
d=""
if [ "$DARK_MODE" = true ]; then
    d="-d"
fi

/geolog -g /GeoLite2-City.mmdb -l $LOG_FILE -p 80 $d
