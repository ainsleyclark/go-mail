#!/usr/local/bin/bash

# Shell script for executing tests based on input.
# Usage:
# ./tests.sh for all drivers
# ./tests.sh sparkpost for a particular driver
# Author - Ainsley Clark

DRIVER=$1

declare -A tests=(
	["mailgun"]="Test_MailGun"
	["postal"]="Test_Postal"
	["postmark"]="Test_Postmark"
	["sendgrid"]="Test_SendGrid"
	["smtp"]="Test_SMTP"
	["sparkpost"]="Test_SparkPost"
)

if [ -z "$DRIVER" ]
then
	for name in "${!tests[@]}";
		do go test -v ./tests/ -run "${tests[$name]}";
	done
else
	go test -v ./tests/ -run "${tests["$DRIVER"]}";
fi
