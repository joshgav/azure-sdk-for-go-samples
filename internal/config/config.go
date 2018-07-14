// package config manages loading configuration from environment and command-line params
// Some of these should be considered base names and defaults rather than exact
// settings.
package config

import (
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/marstr/randname"
)

var (
	clientID               string
	clientSecret           string
	tenantID               string
	subscriptionID         string
	location               string
	resourceURL            string
	authorizationServerURL string
	cloudName              string
	useDeviceFlow          bool
	keepResources          bool
	groupName              string
	baseGroupName          string
)

// ClientID is the OAuth client ID
func ClientID() string {
	return clientID
}

// ClientSecret is the OAuth client secret
func ClientSecret() string {
	return clientSecret
}

// TenantID is the AAD tenant to which this client belongs
func TenantID() string {
	return tenantID
}

// SubscriptionID is a target subscription for resource management
func SubscriptionID() string {
	return subscriptionID
}

// ResourceURL is the URL of a resource for use with OAuth requests
func ResourceURL() string {
	return resourceURL
}

// deprecated: use DefaultLocation() instead
// Location returns the Azure location to be utilized.
func Location() string {
	return location
}

// DefaultLocation() returns the default location wherein to create new resources.
// Some resource types are not available in all locations so another location might need
// to be chosen.
func DefaultLocation() string {
	return location
}

// AuthorizationServerURL is the OAuth authorization server URL.
func AuthorizationServerURL() string {
	return authorizationServerURL
}

// UseDeviceFlow() specifies if interactive auth should be used.
func UseDeviceFlow() bool {
	return useDeviceFlow
}

// deprecated: do not use global group names
// utilize `BaseGroupName()` for a shared prefix
func GroupName() string {
	return groupName
}

// deprecated: we have to set this because we use a global for group names
// once that's fixed this should be removed
func SetGroupName(name string) {
	groupName = name
}

// BaseGroupName() returns a prefix for new groups.
func BaseGroupName() string {
	return baseGroupName
}

func KeepResources() bool {
	return keepResources
}

func Environment() azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		return azure.PublicCloud
	}
	return env
}

// GenerateGroupName leverages BaseGroupName() to return a more detailed name,
// helping to avoid collisions.  It appends each of the `affixes` to
// BaseGroupName() separated by dashes, and adds a 5-character random string.
func GenerateGroupName(affixes ...string) string {
	var b strings.Builder
	b.WriteString(BaseGroupName())
	b.WriteRune('-')
	for _, affix := range affixes {
		b.WriteString(affix)
		b.WriteRune('-')
	}
	return randname.GenerateWithPrefix(b.String(), 5)
}
