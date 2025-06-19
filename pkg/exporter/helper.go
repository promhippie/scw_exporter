package exporter

import (
	log "log/slog"

	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var projectCache map[string]string

func init() {
	projectCache = make(map[string]string, 0)
}

func retrieveProject(client *scw.Client, projectID string) (string, error) {

	if projectName, ok := projectCache[projectID]; ok {
		log.Info("Hit in cache", "project", projectID, "project_name", projectName)
		return projectName, nil
	}

	projectClient := account.NewProjectAPI(client)

	project, err := projectClient.GetProject(&account.ProjectAPIGetProjectRequest{ProjectID: projectID})

	if err != nil {
		log.Error("Got error retrieving project", "project", projectID, "err", err)
		return "", err
	}

	log.Info("Got project from the api", "project", project)
	projectCache[projectID] = project.Name
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
