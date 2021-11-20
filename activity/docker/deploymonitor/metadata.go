package deploymonitor

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Id string `md:"id,required"` // The id of client
}

type Input struct {
	Now int64 `md:"now"` // Current time
}

type Output struct {
	Data []interface{} `md:"data"` // The data recieved
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

	o.Data = values["data"].([]interface{})
	return nil
}
