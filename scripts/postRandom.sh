#!/bin/bash

echo '{"url": "https://en.wikipedia.org/wiki/Special:Random"}' | http \
  -a api:$API_KEY \
  :8080/articles

