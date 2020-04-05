package service

import (
	proto "github.com/WillCoates/FYP/common/protocol/scripting"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*ScriptingService) GetScriptErrors(req *proto.GetScriptErrorsRequest, srv proto.ScriptingService_GetScriptErrorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetScriptErrors not implemented")
}
