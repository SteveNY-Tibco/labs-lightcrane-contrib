package deploymonitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{
		settings: settings,
	}
	return act, nil
}

type Activity struct {
	settings *Settings
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Info("(fnAirDeployMonitor:Eval) entering ........ ")
	defer ctx.Logger().Info("(fnAirDeployMonitor:Eval) exit ........ ")

	input := &Input{}

	err = ctx.GetInputObject(input)

	if err != nil {
		return true, err
	}

	/*
		ProjectID
		Name
		Data {}
		Status
		LastModified
		ErrorCode
		ErrorMessage
	*/

	/* query docker container */
	dctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	containers, err := cli.ContainerList(dctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, err
	}

	currentDeploymnts := make([]interface{}, 0)
	for _, container := range containers {
		ctx.Logger().Info(container.Names[0] + "-" + container.Status)
		containerName := container.Names[0]
		if strings.HasPrefix(containerName, "/Air-") {
			currentDeploymnts = append(currentDeploymnts, map[string]interface{}{
				"ProjectID":    containerName[0:strings.Index(containerName, "_")],
				"Name":         containerName[0:strings.Index(containerName, "_")],
				"Status":       container.Status,
				"LastModified": time.Now().Unix(),
			})
		}
	}

	ctx.Logger().Info("(fnAirDeployMonitor:Eval) currentDeploymnts : ", currentDeploymnts)

	err = ctx.SetOutput("data", currentDeploymnts)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *Activity) Evalx(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Info("(fnAirDeployMonitor:Eval) entering ........ ")
	defer ctx.Logger().Info("(fnAirDeployMonitor:Eval) exit ........ ")

	input := &Input{}

	err = ctx.GetInputObject(input)

	if err != nil {
		return true, err
	}

	data := make([]interface{}, 1)
	data[0] = map[string]interface{}{
		"ProjectID": "Air-account_00001",
		"Name":      "edgex_mqtt_mqtt_fs",
	}

	deployments := make(map[string]interface{})
	deployments["Update"] = make([]interface{}, 0)
	deployments["Remove"] = make([]interface{}, 0)

	/*
		ProjectID
		Name
		Data {}
		Status
		LastModified
		ErrorCode
		ErrorMessage
	*/
	registeredDeployments := make(map[string]interface{})
	if nil != data {
		for _, registeredDeployment := range data {
			projectID := registeredDeployment.(map[string]interface{})["ProjectID"]
			name := registeredDeployment.(map[string]interface{})["Name"]
			registeredDeployments[fmt.Sprintf("/%s_%s", projectID, name)] = registeredDeployment
		}
	}

	fmt.Println("(fnAirDeployMonitor:Eval) registered deployments : ", registeredDeployments)

	/* query docker container */
	dctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	containers, err := cli.ContainerList(dctx, types.ContainerListOptions{All: true})
	if err != nil {
		return
	}

	currentDeploymnts := make(map[string]interface{})
	for _, container := range containers {
		fmt.Println(container.Names[0] + "-" + container.Status)
		currentDeploymnts[container.Names[0]] = container
		registeredDeployment := registeredDeployments[container.Names[0]]
		if nil != registeredDeployment {
			registeredDeployment.(map[string]interface{})["Status"] = container.Status
			deployments["Update"] = append(deployments["Update"].([]interface{}), registeredDeployment)
		}
	}

	ctx.Logger().Info("(fnAirDeployMonitor:Eval) currentDeploymnts : ", currentDeploymnts)
	ctx.Logger().Info("(fnAirDeployMonitor:Eval) deployments : ", deployments)

	return true, nil
}
