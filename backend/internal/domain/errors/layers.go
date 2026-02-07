package errors

// Layer constants for use with oops.In() to identify where an error originated.
const (
	// LayerHandler is the HTTP/API handler layer.
	LayerHandler = "handler"

	// LayerService is the business logic/use case layer.
	LayerService = "service"

	// LayerRepository is the data access layer.
	LayerRepository = "repository"

	// LayerMiddleware is the middleware layer.
	LayerMiddleware = "middleware"

	// LayerInfrastructure is the infrastructure layer (database, cache, etc.).
	LayerInfrastructure = "infrastructure"
)
