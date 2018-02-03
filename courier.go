package courier

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// Instance an instance of a courier, locked by a file.
type Instance struct {
	melody *melody.Melody
	engine *gin.Engine

	InstanceID string
}

// New creates a new instance or throws errors.
func New(engine *gin.Engine) (i *Instance, err error) {
	i = &Instance{
		melody: melody.New(),
		engine: engine,
	}

	return
}

// SetInstanceID sets the instance id to be used when communicating
func (i *Instance) SetInstanceID(id string) {
	i.InstanceID = id
}

func (i *Instance) handler(ctx *gin.Context) {

}
