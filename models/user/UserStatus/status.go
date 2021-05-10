package UserStatus

type UserStatus uint

const (
    Inactive UserStatus = iota
    Transition
    Active
)
