#!/bin/bash
docker build -f Dockerfile -t pind-build:latest .
id=$(docker create pind-build:latest)