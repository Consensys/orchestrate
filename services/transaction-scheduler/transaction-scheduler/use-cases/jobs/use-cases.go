package jobs

type UseCases interface {
	CreateJob() CreateJobUseCase
	GetJob() GetJobUseCase
	StartJob() StartJobUseCase
	UpdateJob() UpdateJobUseCase
	SearchJobs() SearchJobsUseCase
}
