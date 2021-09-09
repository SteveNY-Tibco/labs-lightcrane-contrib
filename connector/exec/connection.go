package exec

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "exec-connection")
var factory = &ExecConnectionFactory{}

type Settings struct {
	Name        string `md:"name,required"`
	Description string `md:"description,required"`
}

func init() {
	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

type EXEConnectionFactory struct {
}

func (*EXEConnectionFactory) Type() string {
	return "Exec"
}

func (*EXEConnectionFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	sharedConn := &EXEConnection{}
	return sharedConn, nil
}

type EXEConnection struct {
	exeEventBrokers map[string]*EXEEventBroker
	mux             sync.Mutex
}

func (this *EXEConnection) Type() string {
	return "Exec"
}

func (this *EXEConnection) GetConnection() interface{} {
	return this
}

func (this *EXEConnection) ReleaseConnection(connection interface{}) {

}

func (this *EXEConnection) getEXEEventBroker(serverId string) *EXEEventBroker {
	return this.exeEventBrokers[serverId]
}

func (this *EXEConnection) createEXEEventBroker(
	serverId string,
	listener EXEEventListener) (*EXEEventBroker, error) {

	this.mux.Lock()
	defer this.mux.Unlock()
	broker := this.exeEventBrokers[serverId]

	broker = &EXEEventBroker{
		listener: listener,
	}
	this.exeEventBrokers[serverId] = broker

	return broker, nil
}

type EXEEventListener interface {
	ProcessEvent(event map[string]interface{}) error
}

type EXEEventBroker struct {
	listener EXEEventListener
}

func (this *EXEEventBroker) Start() {
	fmt.Println("Start broker, EXEEventBroker : ", this)
}

func (this *EXEEventBroker) Stop() {
}

func (this *EXEEventBroker) SendEvent(event map[string]interface{}) {
	this.listener.ProcessEvent(event)
}
