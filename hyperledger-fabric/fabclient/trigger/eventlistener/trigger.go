package eventlistener

import (
	"context"
	"encoding/json"

	client "github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/common"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	conName          = "name"
	conConfig        = "config"
	conEntityMatcher = "entityMatcher"
	conChannel       = "channelID"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})
var logger log.Logger

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Trigger eventlistener trigger struct
type Trigger struct {
	listeners []*Listener
}

// Factory to create a trigger
type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (f *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{}, nil
}

// Initialize implements trigger.Trigger.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	logger = ctx.Logger()
	t.listeners = []*Listener{}
	for _, handler := range ctx.GetHandlers() {
		s := &HandlerSettings{}
		if err := metadata.MapToStruct(handler.Settings(), s, true); err != nil {
			return err
		}
		configs, err := client.GetSettings(s.FabricConnector)
		if err != nil {
			return errors.Wrapf(err, "failed to get settings of fabric connector")
		}
		networkConfig, err := client.ExtractFileContent(configs[conConfig])
		if err != nil {
			return errors.Wrapf(err, "invalid network config")
		}
		entityMatcher, err := client.ExtractFileContent(configs[conEntityMatcher])
		if err != nil {
			return errors.Wrapf(err, "invalid entity-matchers-override")
		}
		spec := Spec{
			Name:           configs[conName].(string),
			NetworkConfig:  networkConfig,
			EntityMatchers: entityMatcher,
			UserName:       s.User,
			OrgName:        s.Org,
			ChannelID:      configs[conChannel].(string),
			EventType:      s.EventType,
			ChaincodeID:    s.ChaincodeID,
			EventFilter:    s.EventFilter,
		}
		// logger.Debugf("initialize event listener spec: %+v", spec)
		listener, err := NewListener(&spec, flowEventHandler(handler))
		if err != nil {
			return errors.Wrapf(err, "failed to crete event listener")
		}
		t.listeners = append(t.listeners, listener)
	}
	return nil
}

// processes event data by passing it to flogo flow
func flowEventHandler(handler trigger.Handler) EventHandler {
	return func(data interface{}) {
		output := &Output{}
		output.Data = data
		if logger.DebugEnabled() {
			if jsonbytes, err := json.MarshalIndent(data, "", "  "); err == nil {
				logger.Debug("Got event data: ", string(jsonbytes))
			}
		}
		if _, err := handler.Handle(context.Background(), output); err != nil {
			logger.Errorf("error processing event: %s", err.Error())
		}
	}
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	for _, c := range t.listeners {
		logger.Debug("start listener of type:", c.eventType)
		if err := c.Start(); err != nil {
			return errors.Wrapf(err, "failed to start event listener")
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Stop
func (t *Trigger) Stop() error {
	for _, c := range t.listeners {
		c.Stop()
	}
	return nil
}
