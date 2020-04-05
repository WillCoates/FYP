package business

import "errors"

// ErrNoSuchSensor is returned when a sensor is expected but couldn't be found
var ErrNoSuchSensor error = errors.New("No such sensor")

// ErrNoSuchField is returned when a field is expected but couldn't be found
var ErrNoSuchField error = errors.New("No such sensor")
