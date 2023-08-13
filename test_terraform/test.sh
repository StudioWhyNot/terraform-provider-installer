#!/bin/bash
#export TF_REATTACH_PROVIDERS='{"registry.terraform.io/shihanng/installer":{"Protocol":"grpc","ProtocolVersion":5,"Pid":2793,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin867798207"}}}'
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