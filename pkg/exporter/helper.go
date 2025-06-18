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

func retrieveProject(client *scw.Client, projectId string) (string, error) {

	if projectName, ok := projectCache[projectId]; ok {
		log.Info("Hit in cache", "project", projectId, "project_name", projectName)
		return projectName, nil
	}

	projectClient := account.NewProjectAPI(client)

	project, err := projectClient.GetProject(&account.ProjectAPIGetProjectRequest{ProjectID: projectId})

	if err != nil {
		log.Error("Got error retrieving project", "project", projectId, "err", err)
		return "", err
	}

	log.Info("Got project from the api", "project", project)
	projectCache[projectId] = project.Name
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
