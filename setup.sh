#!/bin/bash

USERID="$UID"

sudo chown -R 100:101 ./vol/*/conf
sudo find . -type f -name sysctl.conf -exec chown root:root {} \;
sudo chown -R $USERID:$USERID ./vol/odl
