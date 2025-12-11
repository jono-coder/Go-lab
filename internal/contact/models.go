package contact

import (
	"Go-lab/internal/utils"
	"net/mail"
	"time"
)

type Contact struct {
	Id        int          `json:"id"`
	Firstname string       `json:"first_name"`
	Surname   string       `json:"surname"`
	Email     mail.Address `json:"email"`
	CreatedAt time.Time    `json:"created_at"`
}

type Option func(*Contact)

func WithEmail(email mail.Address) Option {
	return func(contact *Contact) {
		contact.Email = email
	}
}

func NewContact(firstName, surname string) *Contact {
	return &Contact{
		Firstname: firstName,
		Surname:   surname,
	}
}

func (c *Contact) String() string {
	return utils.ToString(c)
}

// AvgSize used with caching
func AvgSize() int64 {
	return 8 + 10 + 30 + 20 + 8 // add all fields len
}
