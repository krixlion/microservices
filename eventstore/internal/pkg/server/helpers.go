package server

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fmtDupKeyErr(err error) error {
	badReq := &errdetails.BadRequest{}
	violation := &errdetails.BadRequest_FieldViolation{
		Field:       "event_id",
		Description: err.Error(),
	}
	badReq.FieldViolations = append(badReq.FieldViolations, violation)

	st, statusErr := status.New(codes.AlreadyExists, "Event with given ID already exists").WithDetails(badReq)
	if statusErr != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	err = st.Err()
	return err
}
