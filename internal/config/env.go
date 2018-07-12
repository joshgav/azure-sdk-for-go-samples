package config

import (
	"log"
	"strconv"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/gobuffalo/envy"
)

// ParseEnvironment loads a sibling `.env` file then looks through all environment
// variables to set global configuration.
func ParseEnvironment() error {
	envy.Load()
	azureEnv, _ := azure.EnvironmentFromName("AzurePublicCloud") // shouldn't fail

	// these must be provided by environment
	var err1, err2, err3, err4 error
	clientID, err1 = envy.MustGet("AZURE_CLIENT_ID")
	clientSecret, err2 = envy.MustGet("AZURE_CLIENT_SECRET")
	tenantID, err3 = envy.MustGet("AZURE_TENANT_ID")
	subscriptionID, err4 = envy.MustGet("AZURE_SUBSCRIPTION_ID")

	for _, err := range []error{err1, err2, err3, err4} {
		if err != nil {
			return err
		}
	}

	// we can choose defaults for these
	groupName = envy.Get("AZURE_GROUP_NAME", "azure-go-samples")
	baseGroupName = envy.Get("AZURE_BASE_GROUP_NAME", groupName)
	resourceURL = envy.Get("AZURE_RESOURCE_URL", azureEnv.ResourceManagerEndpoint)
	location = envy.Get("AZURE_LOCATION_DEFAULT", "westus2")

	var err error
	useDeviceFlow, err = strconv.ParseBool(envy.Get("AZURE_USE_DEVICEFLOW", "0"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_USE_DEVICEFLOW, disabling\n")
		useDeviceFlow = false
	}
	keepResources, err = strconv.ParseBool(envy.Get("AZURE_SAMPLES_KEEP_RESOURCES", "0"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_SAMPLES_KEEP_RESOURCES, discarding\n")
		keepResources = false
	}

	authorizationServerURL = azureEnv.ActiveDirectoryEndpoint
	return nil
}
