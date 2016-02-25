#!/usr/bin/env bash

URL="https://thymio.tk/cardlogin"
for i in $(seq ${1:-10}); do
    x=$(../genid/genid -short -key "$SECRET")
    echo "$URL/$x"
done