#!/bin/bash

# Make sure the container name is the same as the one defined in the docker-compose.yaml file
garage="docker exec -i shared_garage /garage"

# Create a bucket
$garage bucket create dynamic-provisioning-bucket

# Create an API Key
key_output=$($garage key create dp-app-key)

key_name=$(echo "$key_output" | grep "Key name:" | awk '{print $3}')
key_id=$(echo "$key_output" | grep "Key ID:" | awk '{print $3}')
secret_key=$(echo "$key_output" | grep "Secret key:" | awk '{print $3}')

# Save to file
timestamp=$(date +"%Y-%m-%d %H:%M:%S")
echo "Key Name: $key_name" > ../garage.keys
echo "Key ID: $key_id" >> ../garage.keys
echo "Secret Key: $secret_key" >> ../garage.keys
echo "Generated at $timestamp" >> ../garage.keys

# Assign the API Key to the bucket with read-write-owner permissions
$garage bucket allow --read --write --owner dynamic-provisioning-bucket --key $key_name