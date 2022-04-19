// ============================================================================
// Deploy food truck app into Azure Container Apps with Azure Maps account
// ============================================================================

targetScope = 'subscription'

@description('Name used for resource group, and base name for most resources')
param name string = 'temp-food-truck'
@description('Azure region for all resources')
param location string = deployment().location

param image string = 'ghcr.io/benc-uk/food-truck:latest'

// ===== Modules & Resources ==================================================

resource resGroup 'Microsoft.Resources/resourceGroups@2021-04-01' = {
  name: name
  location: location
}

module logAnalytics './modules/log-analytics.bicep' = {
  scope: resGroup
  name: 'monitoring'
  params: {
    name: 'food-truck-logs'
  }
}

module mapAccount './modules/maps.bicep' = {
  scope: resGroup
  name: 'mapAccount'
  params: {
    name: 'food-truck-maps'
    sku: 'G2'
  }
}

module containerAppEnv './modules/app-env.bicep' = {
  scope: resGroup
  name: 'containerAppEnv'
  params: {
    name: 'food-truck-env'
    logAnalyticsName: logAnalytics.outputs.name
    logAnalyticsResGroup: resGroup.name
  }
}

module containerApp './modules/app.bicep' = {
  scope: resGroup
  name: 'containerApp'
  params: {
    name: 'food-truck'
    environmentId: containerAppEnv.outputs.id

    image: image
    revisionMode: 'single'

    ingressPort: 8080
    ingressExternal: true
    replicasMax: 1
    replicasMin: 1

    probePath: '/health'
    probePort: 8080

    envs: [
      {
        name: 'AZURE_MAPS_KEY'
        value: mapAccount.outputs.key
      }
    ]
  }
}

// ===== Outputs ==================================================

output appUrl string = 'https://${containerApp.outputs.fqdn}/app/'
