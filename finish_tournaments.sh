#!/bin/bash

# API endpoint
API_URL="http://localhost:8080/tournaments/finish-all"

# Log file path
LOG_FILE="$HOME/Desktop/good-api/tournament_cron.log"

# Ensure the log file exists
touch "$LOG_FILE"

# Log the execution time
echo "Running tournament finish script at $(date)" >> $LOG_FILE

# Make the API request using curl
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL")

# Log the response
echo "Response code: $response" >> $LOG_FILE

# Check if the request was successful
if [ "$response" -eq 200 ]; then
    echo "Tournaments finished successfully." >> $LOG_FILE
else
    echo "Failed to finish tournaments. Check API logs." >> $LOG_FILE
fi
