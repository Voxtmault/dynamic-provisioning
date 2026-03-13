ui            = false
cluster_addr  = "http://openbao:8201"
api_addr      = "http://openbao:8200"

storage "raft" {
  path = "/openbao/file"
  node_id = "master"
}

listener "tcp" {
  address       = "0.0.0.0:8200"
  tls_disable   = true
}