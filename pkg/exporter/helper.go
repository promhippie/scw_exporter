package exporter

import (
	"log/slog"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	projectCache *ttlcache.Cache[string, string]
)

func init() {
	projectCache = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](5 * time.Minute),
	)

	go projectCache.Start()
}

func retrieveProject(logger *slog.Logger, client *scw.Client, projectID string) (string, error) {
	if projectCache.Has(projectID) {
		projectName := projectCache.Get(projectID).Value()

		logger.Debug("Hit project cache",
			"project", projectID,
			"name", projectName,
		)

		return projectName, nil
	}

	project, err := account.NewProjectAPI(
		client,
	).GetProject(
		&account.ProjectAPIGetProjectRequest{
			ProjectID: projectID,
		},
	)

	if err != nil {
		logger.Error("Failed to fetch project",
			"project", projectID,
			"err", err,
		)

		return "", err
	}

	logger.Debug("Set project cache",
		"project", projectID,
		"name", project.Name,
	)

	projectCache.Set(
		projectID,
		project.Name,
		ttlcache.DefaultTTL,
	)

	return project.Name, nil
}

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
