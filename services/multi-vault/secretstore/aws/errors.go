package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/error"
)

// FromAWSError transform an AWS error into an internal error
func FromAWSError(err error) *ierror.Error {
	if err == nil {
		return nil
	}

	if awsErr, ok := err.(awserr.Error); ok {
		// You can refer to https://github.com/aws/aws-sdk-go/blob/master/service/secretsmanager/errors.go
		// for error signification
		switch awsErr.Code() {
		case secretsmanager.ErrCodeDecryptionFailure:
			return errors.CryptoOperationError(awsErr.Message())
		case secretsmanager.ErrCodeEncryptionFailure:
			return errors.CryptoOperationError(awsErr.Message())
		case secretsmanager.ErrCodeInternalServiceError:
			return errors.InternalError(awsErr.Message())
		case secretsmanager.ErrCodeInvalidNextTokenException:
			return errors.DataError(awsErr.Message())
		case secretsmanager.ErrCodeInvalidParameterException:
			return errors.DataError(awsErr.Message())
		case secretsmanager.ErrCodeInvalidRequestException:
			return errors.ConflictedError(awsErr.Message())
		case secretsmanager.ErrCodeLimitExceededException:
			return errors.InsuficientResourcesError(awsErr.Message())
		case secretsmanager.ErrCodeMalformedPolicyDocumentException:
			return errors.InvalidFormatError(awsErr.Message())
		case secretsmanager.ErrCodePreconditionNotMetException:
			return errors.FailedPreconditionError(awsErr.Message())
		case secretsmanager.ErrCodeResourceExistsException:
			return errors.ConstraintViolatedError(awsErr.Message())
		case secretsmanager.ErrCodeResourceNotFoundException:
			return errors.NotFoundError(awsErr.Message())
		default:
			return errors.InternalError(awsErr.Message())
		}
	}

	return errors.FromError(err)
}
