#!/bin/bash

set -e

go build -o bookings cmd/web/*.go \
  && ./bookings