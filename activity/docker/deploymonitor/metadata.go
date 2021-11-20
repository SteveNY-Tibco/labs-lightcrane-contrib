package deploymonitor

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Broker       string `md:"broker,required"` // The broker URL
	Id           string `md:"id,required"`     // The id of client
	Username     string `md:"username"`        // The user's name
	Password     string `md:"password"`        // The user's password
	Store        string `md:"store"`           // The store for message persistence
	KeepAlive    int64  `md:"keepAlive"`       // Keep Alive
	CleanSession bool   `md:"cleanSession"`    // Clean session flag

	Retain    bool                   `md:"retain"`         // Retain Messages
	Topic     string                 `md:"topic,required"` // The topic to publish to
	Qos       int                    `md:"qos"`            // The Quality of Service
	SSLConfig map[string]interface{} `md:"sslConfig"`      // SSL Configuration
}

type Input struct {
	Now int64 `md:"now"` // Current time
}

type Output struct {
	Data interface{} `md:"data"` // The data recieved
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"now": i.Now,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Now, err = coerce.ToInt64(values["now"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Data = values["data"]
	return nil
}
