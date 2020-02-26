package business

import "errors"

// ErrNoSuchSensor is returned when a sensor is expected but couldn't be found
var ErrNoSuchSensor error = errors.New("No such sensor")
