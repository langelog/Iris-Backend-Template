package UserRole

type UserRole uint

type RoleList []UserRole

const (
    External      UserRole = iota // external type of user, only can see public content
    Administrator                 // administrator has super powers on system
    Normal                        // normal user that only has ownership on they content
)
