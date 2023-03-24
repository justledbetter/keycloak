package keycloak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// User representation.
//
// https://github.com/keycloak/keycloak/blob/master/core/src/main/java/org/keycloak/representations/idm/UserRepresentation.java
type User struct {
	ID                         *string              `json:"id,omitempty"`
	CreatedTimestamp           *int64               `json:"createdTimestamp,omitempty"`
	Username                   *string              `json:"username,omitempty"`
	Enabled                    *bool                `json:"enabled,omitempty"`
	Totp                       *bool                `json:"totp,omitempty"`
	EmailVerified              *bool                `json:"emailVerified,omitempty"`
	FirstName                  *string              `json:"firstName,omitempty"`
	LastName                   *string              `json:"lastName,omitempty"`
	Email                      *string              `json:"email,omitempty"`
	DisableableCredentialTypes []string             `json:"disableableCredentialTypes,omitempty"`
	RequiredActions            []string             `json:"requiredActions,omitempty"`
	NotBefore                  *int                 `json:"notBefore,omitempty"`
	Access                     *map[string]bool     `json:"access,omitempty"`
	Attributes                 *map[string][]string `json:"attributes,omitempty"`
}

// Credential representation.
//
// https://github.com/keycloak/keycloak/blob/master/core/src/main/java/org/keycloak/representations/idm/CredentialRepresentation.java
type Credential struct {
	Type      *string `json:"type,omitempty"`
	Value     *string `json:"value,omitempty"`
	Temporary *bool   `json:"temporary,omitempty"`
}

// UsersService ...
type UsersService service

// Create a new user.
func (s *UsersService) Create(ctx context.Context, realm string, user *User) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users", realm)
	req, err := s.keycloak.NewRequest(http.MethodPost, u, user)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// List users.
func (s *UsersService) List(ctx context.Context, realm string) ([]*User, *http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users", realm)
	req, err := s.keycloak.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	res, err := s.keycloak.Do(ctx, req, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, res, nil
}

// GetByID get a single user by ID.
func (s *UsersService) GetByID(ctx context.Context, realm, id string) (*User, *http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s", realm, id)
	req, err := s.keycloak.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var user *User
	res, err := s.keycloak.Do(ctx, req, &user)
	if err != nil {
		return nil, nil, err
	}

	return user, res, nil
}

// GetByUsername get a single user by username.
func (s *UsersService) GetByUsername(ctx context.Context, realm, username string) ([]*User, *http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users?username=%s", realm, username)
	req, err := s.keycloak.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	res, err := s.keycloak.Do(ctx, req, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, res, nil
}

// GetByUsername get a single user by attribute.
func (s *UsersService) GetByAttribute(ctx context.Context, realm, attributeName string, value string) ([]*User, *http.Response, error) {
	// Assume we are on a modern release first.
	var ver string = "22"
	var queryUrl string

	if si, e := s.keycloak.GetServerInfo(); e == nil {
		if si != nil {
			ver = si.SystemInfo.Version
		}
	}

	// If we are on a version that doesn't support q=attr:val syntax:
	//
	if ver < "20" {
		queryUrl = fmt.Sprintf("admin/realms/%s/users?filter=%s=%s", realm, url.PathEscape(attributeName), url.PathEscape(value))
	} else {
		queryUrl = fmt.Sprintf("admin/realms/%s/users?q=%s:%s", realm, url.PathEscape(attributeName), url.PathEscape("\""+value+"\""))
	}
	req, err := s.keycloak.NewRequest(http.MethodGet, queryUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	res, err := s.keycloak.Do(ctx, req, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, res, nil
}

// Update update a single user.
func (s *UsersService) Update(ctx context.Context, realm string, user *User) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s", realm, *user.ID)
	req, err := s.keycloak.NewRequest(http.MethodPut, u, user)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// Delete user.
func (s *UsersService) Delete(ctx context.Context, realm, userID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s", realm, userID)
	req, err := s.keycloak.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// ResetPassword sets or resets the user's password.
func (s *UsersService) ResetPassword(ctx context.Context, realm, userID string, credential *Credential) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/reset-password", realm, userID)
	req, err := s.keycloak.NewRequest(http.MethodPut, u, credential)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// Update user.

// JoinGroup adds user to a group.
func (s *UsersService) JoinGroup(ctx context.Context, realm, userID, groupID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/groups/%s", realm, userID, groupID)
	req, err := s.keycloak.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// LeaveGroup removes a user from a group.
func (s *UsersService) LeaveGroup(ctx context.Context, realm, userID, groupID string) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/groups/%s", realm, userID, groupID)
	req, err := s.keycloak.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// AddRealmRoles adds realm roles to user.
func (s *UsersService) AddRealmRoles(ctx context.Context, realm, userID string, roles []*Role) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/role-mappings/realm", realm, userID)
	req, err := s.keycloak.NewRequest(http.MethodPost, u, roles)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// RemoveRealmRoles removes assigned realm roles from user.
func (s *UsersService) RemoveRealmRoles(ctx context.Context, realm, userID string, roles []*Role) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/role-mappings/realm", realm, userID)
	req, err := s.keycloak.NewRequest(http.MethodDelete, u, roles)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// ListRealmRoles returns a list of realm roles assigned to user.
func (s *UsersService) ListRealmRoles(ctx context.Context, realm, userID string) ([]*Role, *http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/role-mappings/realm", realm, userID)
	req, err := s.keycloak.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var roles []*Role
	res, err := s.keycloak.Do(ctx, req, &roles)
	if err != nil {
		return nil, nil, err
	}

	return roles, res, nil
}

// AddClientRoles adds client roles to user.
func (s *UsersService) AddClientRoles(ctx context.Context, realm, userID, clientID string, roles []*Role) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/role-mappings/clients/%s", realm, userID, clientID)
	req, err := s.keycloak.NewRequest(http.MethodPost, u, roles)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// RemoveClientRoles removes assigned client roles from user.
func (s *UsersService) RemoveClientRoles(ctx context.Context, realm, userID, clientID string, roles []*Role) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/role-mappings/clients/%s", realm, userID, clientID)
	req, err := s.keycloak.NewRequest(http.MethodDelete, u, roles)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// VerifyEmailOptions ...
type VerifyEmailOptions struct {
	ClientID    string `url:"client_id,omitempty"`
	RedirectUri string `url:"redirect_uri,omitempty"`
}

// Send an email-verification email to the user.
// An email contains a link the user can click to verify their email address.
func (s *UsersService) SendVerifyEmail(ctx context.Context, realm, userID string, opts *VerifyEmailOptions) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/send-verify-email", realm, userID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.keycloak.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}

// ExecuteActionsEmailOptions ...
type ExecuteActionsEmailOptions struct {
	ClientID    string `url:"client_id,omitempty"`
	Lifespan    int    `url:"lifespan,omitempty"`
	RedirectUri string `url:"redirect_uri,omitempty"`
}

// ExecuteActionsEmail sends an update account email to the user.
// An email contains a link the user can click to perform a set of required actions.
func (s *UsersService) ExecuteActionsEmail(ctx context.Context, realm, userID string, opts *ExecuteActionsEmailOptions, actions []string) (*http.Response, error) {
	u := fmt.Sprintf("admin/realms/%s/users/%s/execute-actions-email", realm, userID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.keycloak.NewRequest(http.MethodPut, u, actions)
	if err != nil {
		return nil, err
	}

	return s.keycloak.Do(ctx, req, nil)
}
