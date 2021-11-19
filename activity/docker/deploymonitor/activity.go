package deploymonitor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/ssl"
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
	client   mqtt.Client
	topic    Topic
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Debugf("(fnAirDeployMonitor:Eval) entering ........ ")
	defer ctx.Logger().Debugf("(fnAirDeployMonitor:Eval) exit ........ ")

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

	fmt.Println("(fnAirDeployMonitor:Eval) entering ........ ")
	defer fmt.Println("(fnAirDeployMonitor:Eval) exit ........ ")

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
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
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

	fmt.Println("(fnAirDeployMonitor:Eval) currentDeploymnts : ", currentDeploymnts)
	fmt.Println("(fnAirDeployMonitor:Eval) deployments : ", deployments)

	ctx.Logger().Debugf("Published Message: %v", input.Message)

	return true, nil
}
