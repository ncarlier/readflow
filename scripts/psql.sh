#!/bin/bash

docker run -it --rm --network host  postgres:11 psql -h localhost -U postgres -d reader_test
