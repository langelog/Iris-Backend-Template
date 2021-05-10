package types

import "backend/models/user/UserRole"

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type ActivationRequest struct {
    ActivationTicket string `json:"ticket"`
    Password         string `json:"password"`
}

type CreationRequest struct {
    Username string            `json:"username"`
    Password string            `json:"password"`
    Role     UserRole.UserRole `json:"role"`
}
