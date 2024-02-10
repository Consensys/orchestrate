package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"

	"github.com/consensys/orchestrate/pkg/types/entities"

	types "github.com/consensys/orchestrate/pkg/types/api"

	"github.com/consensys/orchestrate/pkg/errors"
	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"
)

func (c *HTTPClient) GetJob(ctx context.Context, jobUUID string) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting job"
			return errors.FromError(err).SetMessage(errMessage).AppendReason(err.Error()).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) GetJobs(ctx context.Context) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	var resp []*types.JobResponse

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while getting jobs"
			return errors.FromError(err).SetMessage(errMessage).AppendReason(err.Error()).ExtendComponent(component)
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) SearchJob(ctx context.Context, filters *entities.JobFilters) ([]*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	var resp []*types.JobResponse

	qParams := url.Values{}
	if len(filters.TxHashes) > 0 {
		qParams.Add("tx_hashes", strings.Join(filters.TxHashes, ","))
	}

	if filters.ChainUUID != "" {
		qParams.Add("chain_uuid", filters.ChainUUID)
	}

	if filters.Status != "" {
		qParams.Add("status", string(filters.Status))
	}

	if !filters.UpdatedAfter.IsZero() {
		qParams.Add("updated_after", filters.UpdatedAfter.Format(time.RFC3339))
	}

	if filters.OnlyParents {
		qParams.Add("only_parents", "true")
	}

	if filters.ParentJobUUID != "" {
		qParams.Add("parent_job_uuid", filters.ParentJobUUID)
	}

	if filters.WithLogs {
		qParams.Add("with_logs", "true")
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + qParams.Encode()
	}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.GetRequest(ctx, c.client, reqURL)
		if err != nil {
			errMessage := "error while searching jobs"
			return errors.FromError(err).SetMessage(errMessage).AppendReason(err.Error()).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, &resp)
	})

	return resp, err
}

func (c *HTTPClient) CreateJob(ctx context.Context, request *types.CreateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs", c.config.URL)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PostRequest(ctx, c.client, reqURL, request)
		if err != nil {
			errMessage := "error while creating job"
			return errors.FromError(err).SetMessage(errMessage).AppendReason(err.Error()).ExtendComponent(component)
		}
		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) UpdateJob(ctx context.Context, jobUUID string, request *types.UpdateJobRequest) (*types.JobResponse, error) {
	reqURL := fmt.Sprintf("%v/jobs/%s", c.config.URL, jobUUID)
	resp := &types.JobResponse{}

	err := callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PatchRequest(ctx, c.client, reqURL, request)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseResponse(ctx, response, resp)
	})

	return resp, err
}

func (c *HTTPClient) StartJob(ctx context.Context, jobUUID string) error {
	reqURL := fmt.Sprintf("%v/jobs/%s/start", c.config.URL, jobUUID)

	return callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseEmptyBodyResponse(ctx, response)
	})
}

func (c *HTTPClient) ResendJobTx(ctx context.Context, jobUUID string) error {
	reqURL := fmt.Sprintf("%v/jobs/%s/resend", c.config.URL, jobUUID)

	return callWithBackOff(ctx, c.config.backOff, func() error {
		response, err := clientutils.PutRequest(ctx, c.client, reqURL, nil)
		if err != nil {
			return err
		}

		defer clientutils.CloseResponse(response)
		return httputil.ParseEmptyBodyResponse(ctx, response)
	})
}
