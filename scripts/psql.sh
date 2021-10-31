#!/bin/bash

docker run -it --rm --network host postgres:14 psql -h localhost -U postgres -d readflow_test
