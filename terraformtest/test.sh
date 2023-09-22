#!/bin/bash
#export TF_REATTACH_PROVIDERS='{"registry.terraform.io/shihanng/installer":{"Protocol":"grpc","ProtocolVersion":6,"Pid":25228,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin2075873966"}}}'
export TF_CLI_CONFIG_FILE=$(pwd)/../.terraformrc
export TF_LOG_PROVIDER=""
export TF_LOG="WARN"
# Set default value for input
input=${1:-"up"}

# Check input value and run appropriate command
if [ "$input" = "up" ]; then
    terraform apply --auto-approve
elif [ "$input" = "down" ]; then
    terraform destroy --auto-approve
else
    echo "Invalid input. Please enter 'up' or 'down'."
fi

#terraform-provider-installer_v0.6.1-SNAPSHOT-d18e997
#dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient attach $(pgrep -f terraform-provider-installer)
#dlv exec --accept-multiclient --continue --headless /tmp/tfproviders/* -- -debug --check-go-version=false