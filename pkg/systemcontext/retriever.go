package systemcontext

type ContextRetriever interface {
	RetrieveContext() (string, error)
}
