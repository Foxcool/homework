package storage

type UserStorage interface {
	StoreUser(in User) (User, error)
	UpdateUser(id string, in User) (User, error)
	GetUser(in User) (User, error)
	GetUsers() ([]User, error)
	DeleteUser(id string) error
}

type User struct {
	ID *string `json:"id" bson:"_id"`

	FirstName  *string `json:"firstName" bson:"firstName"`
	LastName   *string `json:"lastName" bson:"lastName"`
	MiddleName *string `json:"middleName" bson:"middleName"`

	Mobile *string `json:"mobile" bson:"mobile"`
	Email  *string `json:"email" bson:"email"`
}
