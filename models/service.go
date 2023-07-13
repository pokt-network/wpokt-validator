package models

type Service interface {
	Start()
	Stop()
}

type EmptyService struct{}

func (e *EmptyService) Start() {}

func (e *EmptyService) Stop() {}

func NewEmptyService() *EmptyService {
	return &EmptyService{}
}
