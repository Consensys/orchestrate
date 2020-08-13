backend "file" {
    path            = "/vault/file"
}

listener "tcp" {
    address         = "vault:8200"
    tls_disable     = true
}

default_lease_ttl   = "15m"
max_lease_ttl       = "1h"

log_level = "Debug"

ui = true
