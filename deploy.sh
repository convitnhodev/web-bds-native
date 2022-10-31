#!/bin/bash
set -e

SITE_NAME=deein
SERVER=deein
ENTRY_FILE=./cmd/deein/main.go
BUILD_PATH=./build/$SITE_NAME
CONFIG_DIR=./deploy/$SITE_NAME
THEME_PATH=./ui
NGINX_AVAIL_PATH=/etc/nginx/sites-available/$SITE_NAME
NGINX_ENABLED_PATH=/etc/nginx/sites-enabled/$SITE_NAME

SV_PATH=/etc/sv/$SITE_NAME

deploy() {
  echo "Building binary..."
  GOOS=linux go build -o "$BUILD_PATH" "$ENTRY_FILE" && chmod +x "$BUILD_PATH"
  echo "Uploading binary..."
  rsync -aP "$BUILD_PATH" root@"$SERVER":/usr/local/bin/
  echo "Uploading theme folder..."
  ssh root@"$SERVER" "mkdir -p $SV_PATH/theme"
  rsync -ahP "$THEME_PATH/" root@"$SERVER":"$SV_PATH/theme"
  echo "Uploading config files..."
  rsync -aP "$CONFIG_DIR/config.json" "$CONFIG_DIR/run" root@"$SERVER":"$SV_PATH"
  if test -f "$CONFIG_DIR/robots.txt"; then
    rsync -aP "$CONFIG_DIR/robots.txt" root@"$SERVER":"$SV_PATH"
  fi
  if test -f "$CONFIG_DIR/ads.txt"; then
    rsync -aP "$CONFIG_DIR/ads.txt" root@"$SERVER":"$SV_PATH"
  fi

  rsync -aP "./sql" root@"$SERVER":"$SV_PATH"


  rsync -aP "$CONFIG_DIR/nginx.conf" root@"$SERVER":"$NGINX_AVAIL_PATH"
  echo "Starting service on server..."
  ssh root@"$SERVER" "
    ln -sf $SV_PATH /etc/service/
    chmod +x $SV_PATH/run
    ln -sf $NGINX_AVAIL_PATH $NGINX_ENABLED_PATH
    nginx -t && systemctl reload nginx
    sv restart $SITE_NAME
  "
  echo "Done."
}

deploy