#!/bin/bash
set -e

# Deployment wrapper script

which az > /dev/null || { echo "š„ Error! Azure CLI not found, please install https://aka.ms/azure-cli"; exit 1; }
az bicep version > /dev/null || { echo "š„ Error! Bicep not installed in Azure CLI, run 'az bicep install'"; exit 1; }

for varName in AZURE_DEPLOY_NAME AZURE_REGION; do
  varVal=$(eval echo "\${$varName}")
  [ -z "$varVal" ] && { echo "š„ Error! Required variable '$varName' is unset!"; varUnset=true; }
done
[ $varUnset ] && exit 1

echo -e "\n\nš Deployment started..."
echo -e "š Resource group: $AZURE_DEPLOY_NAME"
echo -e "š Region: $AZURE_REGION"
echo -e "š¦ Image: $AZURE_IMAGE\n"

az deployment sub create                 \
  --template-file deploy/main.bicep      \
  --location "$AZURE_REGION"             \
  --name food-truck                      \
  --parameters name="$AZURE_DEPLOY_NAME" \
               image="$AZURE_IMAGE"                   

echo -e "\nāØ Deployment complete!\nš Start finding food trucks here: $(az deployment sub show --name food-truck --query 'properties.outputs.appUrl.value' -o tsv)"
