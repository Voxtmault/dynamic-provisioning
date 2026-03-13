#!/bin/bash

openbao="docker exec -ti shared_openbao /bao"

$openbao status > /dev/null 2>&1
exit_code=$?

if [ $exit_code -eq 1 ]; then
    echo "OpenBao is not initialized, initializing..."
    init_output=$($openbao operator init)

    # Save to file
    echo "$init_output" > ./openbao-init.txt
    echo "Init output saved to openbao-init.txt, please keep this file safe as it contains the unseal keys and root token."

    # Parse the unseal keys
    unseal_key_1=$(grep "Unseal Key 1:" ./openbao-init.txt | awk '{print $4}')
    unseal_key_2=$(grep "Unseal Key 2:" ./openbao-init.txt | awk '{print $4}')
    unseal_key_3=$(grep "Unseal Key 3:" ./openbao-init.txt | awk '{print $4}')
    root_token=$(grep "Initial Root Token:" ./openbao-init.txt | awk '{print $4}')

    rootbao="docker exec -e VAULT_TOKEN=$root_token -it shared_openbao /bao"

    echo "Unsealing OpenBao..."
    $openbao operator unseal $unseal_key_1
    $openbao operator unseal $unseal_key_2
    $openbao operator unseal $unseal_key_3

    echo "OpenBao unsealed successfully"

    echo "Enabling KV secrets engine at secret/ with version 2..."
    $openbao secrets enable -path=secret kv-v2

    echo "Copying policy file into OpenBao container..."
    docker cp ./openbao-policy.hcl shared_openbao:/openbao-policy.hcl

    echo "Applying policy..."
    $rootbao policy write dp-policy /openbao-policy.hcl

    echo "Enabling AppRole auth method..."
    $rootbao auth enable approle

    echo "Creating AppRole role for dp with dp-policy attached..."
    $rootbao write auth/approle/role/dp token_policies="dp-policy" \
    token_ttl=1h token_max_ttl=4h

    echo "Generating Role ID for dp AppRole..."
    role_id_output=$($rootbao read auth/approle/role/dp/role-id)
    role_id=$(echo "$role_id_output" | grep "role_id" | awk '{print $2}')

    echo "Generating Secret ID for dp AppRole..."
    secret_id_output=$($rootbao write -f auth/approle/role/dp/secret-id)
    secret_id=$(echo "$secret_id_output" | grep "secret_id " | awk '{print $2}')

    # Save to file
    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    echo "Role ID: $role_id" > ./openbao-approle.txt
    echo "Secret ID: $secret_id" >> ./openbao-approle.txt
    echo "Generated at $timestamp" >> ./openbao-approle.txt

    echo "AppRole credentials saved to openbao-approle.txt, please keep this file safe."
else
    echo "OpenBao is already initialized, skipping..."
fi