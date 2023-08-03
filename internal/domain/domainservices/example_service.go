package domainservices

type ExampleService struct {
}

// Constructor of ExampleService
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// wire

var singletonExampleService *ExampleService = initSingletonExampleService()

func GetSingletonExampleService() *ExampleService {
	return singletonExampleService
}

func initSingletonExampleService() *ExampleService {
	return NewExampleService()
}

// methods

func (p *ExampleService) Hello() {
}
