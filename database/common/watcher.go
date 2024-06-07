package common

type (
	DatabaseEntityType string
	OperationType      string
	PayloadFilterFunc  func(ChangePayload) bool
)

const (
	RepositoryEntityType   DatabaseEntityType = "repository"
	OrganizationEntityType DatabaseEntityType = "organization"
	EnterpriseEntityType   DatabaseEntityType = "enterprise"
	PoolEntityType         DatabaseEntityType = "pool"
	UserEntityType         DatabaseEntityType = "user"
	InstanceEntityType     DatabaseEntityType = "instance"
	JobEntityType          DatabaseEntityType = "job"
)

const (
	CreateOperation OperationType = "create"
	UpdateOperation OperationType = "update"
	DeleteOperation OperationType = "delete"
)

type ChangePayload struct {
	EntityType DatabaseEntityType
	Operation  OperationType
	Payload    interface{}
}

type Consumer interface {
	Watch() <-chan ChangePayload
	IsClosed() bool
	Close()
	SetFilters(filters ...PayloadFilterFunc)
}

type Producer interface {
	Notify(ChangePayload) error
	IsClosed() bool
	Close()
}

type Watcher interface {
	RegisterProducer(ID string) (Producer, error)
	RegisterConsumer(ID string, filters ...PayloadFilterFunc) (Consumer, error)
}
