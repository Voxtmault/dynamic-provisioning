#!/bin/bash

openbao="docker exec -i shared_openbao /bao"

# Parse the unseal keys
unseal_key_1=$(grep "Unseal Key 1:" ./openbao-init.txt | awk '{print $4}')
unseal_key_2=$(grep "Unseal Key 2:" ./openbao-init.txt | awk '{print $4}')
unseal_key_3=$(grep "Unseal Key 3:" ./openbao-init.txt | awk '{print $4}')

echo "Unsealing OpenBao..."

$openbao operator unseal $unseal_key_1
$openbao operator unseal $unseal_key_2
$openbao operator unseal $unseal_key_3

echo "OpenBao unsealed successfully"