package service

import (
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
	"github.com/WillCoates/FYP/sensors/business"
)

type SensorsService struct {
	proto.UnimplementedSensorsServiceServer
	logic *business.Logic
}

func NewSensorsService(logic *business.Logic) *SensorsService {
	service := new(SensorsService)
	service.logic = logic
	return service
}
