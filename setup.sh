#!/bin/bash

USERID="$UID"

sudo chown -R 100:101 ./vol/frr
sudo chown -R $USERID:$USERID ./vol/odl
