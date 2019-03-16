#!/bin/bash

cat $1 | http \
  -a api:89b5700d-e4da-407e-94a0-7303417189c5 \
  :8081/articles
