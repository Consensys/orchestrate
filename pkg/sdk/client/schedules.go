package client

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	types "github.com/consensys/orchestrate/pkg/types/api"

	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
)

func (c *HTTPClient) CreateSchedule(ctx context.Context, request *types.CreateScheduleRequest) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)
	resp := &types.ScheduleResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetSchedule(ctx context.Context, scheduleUUID string) (*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules/%v", c.config.URL, scheduleUUID)
	resp := &types.ScheduleResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetSchedules(ctx context.Context) ([]*types.ScheduleResponse, error) {
	reqURL := fmt.Sprintf("%v/schedules", c.config.URL)
	var resp []*types.ScheduleResponse

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}
