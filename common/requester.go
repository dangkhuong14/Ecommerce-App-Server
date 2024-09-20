package common

// import "github.com/google/uuid"

type Requester interface {
	UserId() UUID
	TokenId() UUID
	FirstName() string
	LastName() string
	Role() string
	Status() string
}

type requesterData struct {
	userId    UUID
	tid       UUID
	firstName string
	lastName  string
	role      string
	status    string
}

func (r *requesterData) UserId() UUID {
	return r.userId
}
func (r *requesterData) TokenId() UUID {
	return r.tid
}
func (r *requesterData) FirstName() string { return r.firstName }
func (r *requesterData) LastName() string  { return r.lastName }
func (r *requesterData) Role() string      { return r.role }
func (r *requesterData) Status() string    { return r.status }

func NewRequester(sub, tid UUID, firstName, lastName, role, status string) Requester {
	return &requesterData{
		userId:    sub,
		tid:       tid,
		firstName: firstName,
		lastName:  lastName,
		role:      role,
		status:    status,
	}
}