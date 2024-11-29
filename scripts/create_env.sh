#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_ROOT_PATH=$SCRIPT_DIR"/.."
ENV_FILE_PATH=$PROJECT_ROOT_PATH"/.env"



CONFIG_PATH=./config/config.yaml

MAIL_EMAIL=golang.auth.service@gmail.com
MAIL_PASSWORD="jybh ayjb qosq kykn"
MAIL_ADDRESS=smtp.gmail.com
MAIL_PORT=587


echo "CONFIG_PATH=$CONFIG_PATH" > $ENV_FILE_PATH

echo "MAIL_EMAIL=$MAIL_EMAIL" >> $ENV_FILE_PATH
echo "MAIL_PASSWORD=$MAIL_PASSWORD" >> $ENV_FILE_PATH
echo "MAIL_ADDRESS=$MAIL_ADDRESS" >> $ENV_FILE_PATH
echo "MAIL_PORT=$MAIL_PORT" >> $ENV_FILE_PATH
