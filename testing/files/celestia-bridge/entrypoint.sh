#!/bin/bash

# Copy config files from testing/files/celestia-validator to .tmp directory
cp -r /testapp_files/* /home/celestia/

# Set permissions for keys
chmod 600 keys/OAZHALLLMV4Q 
chmod 600 keys/NJ3XILLTMVRXEZLUFZVHO5A

echo "$@"

exec $@