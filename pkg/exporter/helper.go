package exporter

import (
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func affectedZones(client *scw.Client) []scw.Zone {
	if zone, exists := client.GetDefaultZone(); exists {
		return []scw.Zone{
			zone,
		}
	}

	if region, exists := client.GetDefaultRegion(); exists {
		return region.GetZones()
	}

	return scw.AllZones
}
