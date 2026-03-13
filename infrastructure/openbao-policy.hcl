# Access secrets under any prefix, e.g. secret/data/DEBUG/dp/*
path "secret/data/+/dp/*" {
capabilities = ["create", "read", "update", "delete", "patch", "list"]
}

# Metadata access

path "secret/metadata/+/dp/*" {
capabilities = ["create", "read", "update", "delete", "patch", "list"]
}