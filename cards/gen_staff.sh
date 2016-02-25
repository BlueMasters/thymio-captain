#!/usr/bin/env bash

URL="https://thymio.tk/cardlogin"
SOURCE="staff.png"
DEST="admin_pages_all.pdf"

source lib.sh
source secret.sh
SECRET=$ADMIN_SECRET

pages=${1:-3}
doc $pages