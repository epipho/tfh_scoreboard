#! /bin/bash

set -eu

if [ $# -ne 1 ]; then
    echo "Usage: install.sh api_key"
    exit 1
fi

API_KEY=$1

systemctl stop tfh_scoreboard || true

mkdir -p /opt/tfh_scoreboard
cp ./bin/tfh_scoreboard /opt/tfh_scoreboard
cp -r ./ui /opt/tfh_scoreboard

mkdir -p /var/lib/tfh_scoreboard

cp tfh_scoreboard.service /etc/systemd/system/tfh_scoreboard.service
sed -i "s/{{API_KEY}}/$API_KEY/" /etc/systemd/system/tfh_scoreboard.service

systemctl daemon-reload
systemctl start tfh_scoreboard
