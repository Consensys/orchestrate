package jobs

type UseCases interface {
	CreateJob() CreateJobUseCase
	GetJob() GetJobUseCase
	StartJob() StartJobUseCase
	StartNextJob() StartNextJobUseCase
	UpdateJob() UpdateJobUseCase
	SearchJobs() SearchJobsUseCase
}
