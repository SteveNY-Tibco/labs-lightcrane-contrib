/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tableupsert

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/table"
	"github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/util"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-table-upsert")

const (
	setting_Table   = "Table"
	setting_Indices = "Indices"
	input           = "Mapping"
	output_Data     = "Data"
	output_Exists   = "Exists"
)

// TableUpsertActivity is an Activity that is used to Filter a message to the console
type TableUpsertActivity struct {
	metadata           *activity.Metadata
	activityToTable    map[string]string
	activityToKeyNames map[string][]string
	mux                sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aTableActivity := &TableUpsertActivity{
		metadata:           metadata,
		activityToTable:    make(map[string]string),
		activityToKeyNames: make(map[string][]string),
	}
	return aTableActivity
}

// Metadata returns the activity's metadata
func (a *TableUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *TableUpsertActivity) Eval(ctx activity.Context) (done bool, err error) {

	log.Debug("(TableUpsertActivity.Eval) entering ..... ")
	defer log.Debug("(TableUpsertActivity.Eval) exit ..... ")

	_, myTable, err := a.getTable(ctx)

	if nil != err {
		return false, err
	}

	iData := ctx.GetInput(input).(*data.ComplexObject).Value

	log.Debug("(Eval)iData.Value = ", iData, ", type = ", reflect.TypeOf(iData))

	outputTuple, _ := myTable.Upsert(iData.(map[string]interface{}))
	exists := outputTuple != nil

	log.Debug("(Eval)output tuple = ", outputTuple, ", exist = ", exists, ", row count = ", myTable.RowCount())

	complexdata := &data.ComplexObject{Metadata: "Data", Value: outputTuple}
	ctx.SetOutput(output_Data, complexdata)
	ctx.SetOutput(output_Exists, exists)

	return true, nil
}

func (a *TableUpsertActivity) getTable(context activity.Context) ([]string, *table.Table, error) {
	myId := util.ActivityId(context)

	myTable := table.GetTableManager().GetTable(a.activityToTable[myId])
	keyName := a.activityToKeyNames[myId]
	if nil == myTable {
		a.mux.Lock()
		defer a.mux.Unlock()

		myTable = table.GetTableManager().GetTable(a.activityToTable[myId])
		keyName = a.activityToKeyNames[myId]
		if nil == myTable {

			log.Debug("(getTable) init : ", "initialize table begin ....")

			iTableInfo, exist := context.GetSetting(setting_Table)
			if !exist {
				return nil, nil, activity.NewError("(getTable)Table is not configured", "TABLE_UPSERT-4002", nil)
			}

			//Read table details
			tableInfo, _ := data.CoerceToObject(iTableInfo)
			if tableInfo == nil {
				return nil, nil, activity.NewError("(getTable)Unable extract table details", "TABLE_UPSERT-4001", nil)
			}

			var tablename string
			var schema []interface{}
			tableSettings, _ := tableInfo["settings"].([]interface{})
			if tableSettings != nil {
				for _, v := range tableSettings {
					setting, _ := data.CoerceToObject(v)

					if nil != setting {
						if setting["name"] == "schema" {
							iSchema := setting["value"]
							if nil == iSchema {
								return nil, nil, activity.NewError("(getTable)Unable to get model string", "TABLE_UPSERT-4004", nil)
							}
							err := json.Unmarshal([]byte(iSchema.(string)), &schema)
							if nil != err {
								return nil, nil, err
							}
						} else if setting["name"] == "name" {
							tablename = setting["value"].(string)
						}
					}
				}
			}

			if "" == tablename {
				return nil, nil, activity.NewError("(getTable)Unable to get table name", "TABLE_UPSERT-4003", nil)
			}

			log.Debug("-============= TABLE SCHEMA ================-")
			log.Debug(schema)
			log.Debug("-===========================================-")

			keyName = make([]string, 0)
			indexible := make([]string, 0)
			schemaArray := make([](map[string]interface{}), len(schema))
			for index, field := range schema {
				schemaArray[index] = field.(map[string]interface{})
				if "yes" == schemaArray[index]["IsKey"].(string) {
					keyName = append(keyName, schemaArray[index]["Name"].(string))
					indexible = append(indexible, schemaArray[index]["Name"].(string))
				}
				if "yes" == schemaArray[index]["Indexed"].(string) {
					indexible = append(indexible, schemaArray[index]["Name"].(string))
				}
			}

			myTable = table.GetTableManager().GetTable(tablename)
			if nil == myTable {
				tableSchema := table.CreateSchema(&schemaArray)
				myTable = table.GetTableManager().CreateTable(
					keyName,
					tablename,
					tableSchema,
				)

				for index := 0; index < len(indexible); index++ {
					myTable.GenerateKeys(indexible, make([]string, index+1), 0, len(indexible)-1, 0, index+1)
				}
			}

			log.Debug("(getTable) init : ", "initialize table done : myTable = ", myTable)
			a.activityToTable[myId] = tablename
		}
	}

	return keyName, myTable, nil
}
