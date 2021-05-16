#!/bin/bash

source .env

authority=https://login.readflow.app/auth/realms/readflow
endpoint=${endpoint:-https://api.readflow.app}

payload="grant_type=client_credentials&client_id=${client_id}&client_secret=${client_secret}"

res=`curl -s -k --data "${payload}" ${authority}/protocol/openid-connect/token`
access_token=`echo $res | jq -r .access_token`

curl -s -H "Authorization: Bearer $access_token" ${endpoint}$1

