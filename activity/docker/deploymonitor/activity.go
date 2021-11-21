package deploymonitor

import (
	"context"
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
		ID
		Domain
		Name
		Data {}
		Status
		Reporter
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
			name := containerName[strings.Index(containerName[strings.Index(containerName, "_")+1:], "_")+len(containerName[0:strings.Index(containerName, "_")])+1:]
			currentDeploymnts = append(currentDeploymnts, map[string]interface{}{
				"ID":           containerName[1:],
				"Domain":       containerName[1:strings.Index(containerName, name)],
				"Name":         name[1:],
				"Status":       container.Status,
				"Reporter":     "Deployer",
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
