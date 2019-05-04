#!/bin/bash

source .env

authority=https://login.nunux.org/auth/realms/readflow
endpoint=https://api.readflow.app
#endpoint=http://localhost:8080

payload="grant_type=client_credentials&client_id=${client_id}&client_secret=${client_secret}"


res=`curl -s -k --data "${payload}" ${authority}/protocol/openid-connect/token`
access_token=`echo $res | jq -r .access_token`

curl -s -H "Authorization: Bearer $access_token" ${endpoint}$1
