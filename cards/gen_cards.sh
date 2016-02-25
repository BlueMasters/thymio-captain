#!/usr/bin/env bash

URL="https://thymio.tk/start"
SOURCE="thymio-captain_s.png"
DEST="user_pages_all.pdf"

source lib.sh
source secret.sh
SECRET=$USER_SECRET

pages=${1:-3}
doc $pages