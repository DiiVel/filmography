package entities

type Role string

var Admin Role = "admin"
var User Role = "user"

// Actor model
// @SWG.Model
type UserEntity struct {
	Role     Role
	ID       string
	Username string
}
