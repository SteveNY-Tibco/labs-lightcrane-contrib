// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

// embedded flogo app descriptor file
const flogoJSON string = `{
    "imports": [
        "github.com/TIBCOSoftware/labs-air-contrib/activity/log",
        "github.com/project-flogo/flow",
        "github.com/project-flogo/contrib/activity/rest",
        "github.com/project-flogo/contrib/trigger/timer",
        "github.com/project-flogo/contrib/function/datetime"
    ],
    "name": "app",
    "description": "",
    "version": "1.0.0",
    "type": "flogo:app",
    "appModel": "1.1.1",
    "triggers": [
        {
            "ref": "#timer",
            "name": "flogo-timer",
            "description": "Simple Timer trigger",
            "settings": {},
            "id": "Timer",
            "handlers": [
                {
                    "description": "",
                    "settings": {
                        "startDelay": "",
                        "repeatInterval": "10s"
                    },
                    "action": {
                        "ref": "github.com/project-flogo/flow",
                        "settings": {
                            "flowURI": "res://flow:test"
                        },
                        "input": {
                            "Now": "=datetime.now()"
                        }
                    },
                    "name": "test"
                }
            ]
        }
    ],
    "resources": [
        {
            "id": "flow:test",
            "data": {
                "name": "test",
                "description": "",
                "links": [
                    {
                        "id": 1,
                        "from": "InputLog",
                        "to": "TargetActivity",
                        "type": "default"
                    },
                    {
                        "id": 2,
                        "from": "TargetActivity",
                        "to": "OutputLog",
                        "type": "default"
                    }
                ],
                "tasks": [
                    {
                        "id": "InputLog",
                        "name": "InputLog",
                        "description": "Logs a message",
                        "activity": {
                            "ref": "#log",
                            "input": {
                                "message": "=\"** Input **, Data = \" + $flow.Now",
                                "addDetails": false
                            }
                        }
                    },
                    {
                        "id": "TargetActivity",
                        "name": "TargetActivity",
                        "description": "Invokes a REST Service",
                        "activity": {
                            "ref": "#rest",
                            "settings": {
                                "method": "GET",
                                "uri": "http://abc",
                                "headers": "",
                                "proxy": "",
                                "timeout": 0,
                                "sslConfig": ""
                            },
                            "input": {
                                "pathParams": {
                                    "mapping": {}
                                },
                                "queryParams": {
                                    "mapping": {}
                                },
                                "headers": {
                                    "mapping": {}
                                },
                                "content": "Content"
                            }
                        }
                    },
                    {
                        "id": "OutputLog",
                        "name": "OutputLog",
                        "description": "Logs a message",
                        "activity": {
                            "ref": "#log",
                            "input": {
                                "message": "=\"** Output **, Data = \" + $activity[TargetActivity].data",
                                "addDetails": false
                            }
                        }
                    }
                ]
            }
        }
    ],
    "properties": []
}`
const engineJSON string = ``

func init () {
	cfgJson = flogoJSON
	cfgEngine = engineJSON
}
