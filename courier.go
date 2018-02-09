package courier

import (
	"github.com/DankBotList/Courier/master"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// Instance an instance of a courier, locked by a file.
type Instance struct {
	melody *melody.Melody
	engine *gin.Engine
	config Config
	hub    *master.Hub

	InstanceID string
}

// New creates a new instance or throws errors.
func New(engine *gin.Engine) (i *Instance, err error) {
	i = &Instance{
		melody: melody.New(),
		engine: engine,
		config: Config{},
	}

	if err = i.config.Load(); err != nil {
		return
	}

	i.hub = master.NewHub(i.config.AuthenticationKey)

	engine.GET(i.config.WebSocketPath, func(ctx *gin.Context) {
		i.hub.ServeWs(ctx.Writer, ctx.Request)
	})

	// TODO build stuffs

	return
}

// SetInstanceID sets the instance id to be used when communicating
func (i *Instance) SetInstanceID(id string) {
	i.InstanceID = id
}
