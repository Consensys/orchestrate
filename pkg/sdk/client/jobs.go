package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	clientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/client-utils"
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

	var qParams []string
	if len(filters.TxHashes) > 0 {
		qParams = append(qParams, "tx_hashes="+strings.Join(filters.TxHashes, ","))
	}

	if filters.ChainUUID != "" {
		qParams = append(qParams, "chain_uuid="+filters.ChainUUID)
	}

	if filters.Status != "" {
		qParams = append(qParams, "status="+string(filters.Status))
	}

	if !filters.UpdatedAfter.IsZero() {
		qParams = append(qParams, "updated_after="+filters.UpdatedAfter.Format(time.RFC3339))
	}

	if filters.OnlyParents {
		qParams = append(qParams, "only_parents=true")
	}

	if filters.ParentJobUUID != "" {
		qParams = append(qParams, "parent_job_uuid="+filters.ParentJobUUID)
	}

	if len(qParams) > 0 {
		reqURL = reqURL + "?" + strings.Join(qParams, "&")
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
