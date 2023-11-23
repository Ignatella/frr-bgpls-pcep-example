#!/bin/bash

USERID="$UID"

sudo chown -R $USERID:$USERID .

# clear log files
sudo find . -type f -name *.log -exec truncate -s 0 {} \;

# remove sav files
sudo find . -type f -name *.sav -exec rm {} \;
