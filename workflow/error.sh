#!/usr/bin/env bash

cat << EOB
{"items": [
    {
        "title": "Error!",
        "subtitle": "$error_msg",
    },
    {
        "title": "Clear credentials to reauthenticate",
    },
]}
EOB
