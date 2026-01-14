package verbosity

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetOrganizationIDs retrieves all available organization IDs.
//
// API: GET /core/org/sync
func (c *Client) GetOrganizationIDs() (*OrgSyncResponse, error) {
	req, err := c.newRequest(http.MethodGet, "/core/org/sync", nil, nil)
	if err != nil {
		return nil, err
	}

	var response OrgSyncResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrganizationsByIDs retrieves organization information by their IDs.
//
// API: GET /core/org?ids=11,12,15
func (c *Client) GetOrganizationsByIDs(ids []int64) (*OrgsResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids slice cannot be empty")
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.FormatInt(id, 10)
	}

	params := url.Values{
		"ids": {strings.Join(idStrings, ",")},
	}

	req, err := c.newRequest(http.MethodGet, "/core/org", params, nil)
	if err != nil {
		return nil, err
	}

	var response OrgsResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrganizationByID retrieves a single organization by ID.
func (c *Client) GetOrganizationByID(id int64) (*Org, error) {
	response, err := c.GetOrganizationsByIDs([]int64{id})
	if err != nil {
		return nil, err
	}

	if len(response.Orgs) == 0 {
		return nil, fmt.Errorf("organization with id %d not found", id)
	}

	return &response.Orgs[0], nil
}

// GetAllOrganizations retrieves information about all available organizations.
// This method first gets all organization IDs, then fetches their details.
func (c *Client) GetAllOrganizations() (*OrgsResponse, error) {
	syncResponse, err := c.GetOrganizationIDs()
	if err != nil {
		return nil, err
	}

	if len(syncResponse.IDs) == 0 {
		return &OrgsResponse{Orgs: []Org{}}, nil
	}

	return c.GetOrganizationsByIDs(syncResponse.IDs)
}

// GetMyOrganizations returns organizations where the bot is a member.
func (c *Client) GetMyOrganizations() (*OrgsResponse, error) {
	orgs, err := c.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	var myOrgs []Org
	for i := range orgs.Orgs {
		if orgs.Orgs[i].IsMember {
			myOrgs = append(myOrgs, orgs.Orgs[i])
		}
	}

	return &OrgsResponse{Orgs: myOrgs}, nil
}

// GetAdminOrganizations returns organizations where the bot is an admin.
func (c *Client) GetAdminOrganizations() (*OrgsResponse, error) {
	orgs, err := c.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	var adminOrgs []Org
	for i := range orgs.Orgs {
		if orgs.Orgs[i].IsAdmin {
			adminOrgs = append(adminOrgs, orgs.Orgs[i])
		}
	}

	return &OrgsResponse{Orgs: adminOrgs}, nil
}

// FindOrganizationByTitle searches for an organization by title.
func (c *Client) FindOrganizationByTitle(title string) (*Org, error) {
	orgs, err := c.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	for i := range orgs.Orgs {
		if orgs.Orgs[i].Title == title {
			return &orgs.Orgs[i], nil
		}
	}

	return nil, fmt.Errorf("organization with title %q not found", title)
}

// FindOrganizationBySlug searches for an organization by slug.
func (c *Client) FindOrganizationBySlug(slug string) (*Org, error) {
	orgs, err := c.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	for i := range orgs.Orgs {
		if orgs.Orgs[i].Slug == slug {
			return &orgs.Orgs[i], nil
		}
	}

	return nil, fmt.Errorf("organization with slug %q not found", slug)
}

// OrganizationMembers returns all member IDs of an organization.
func (c *Client) OrganizationMembers(orgID int64) ([]int64, error) {
	org, err := c.GetOrganizationByID(orgID)
	if err != nil {
		return nil, err
	}
	return org.Users, nil
}

// OrganizationAdmins returns all admin IDs of an organization.
func (c *Client) OrganizationAdmins(orgID int64) ([]int64, error) {
	org, err := c.GetOrganizationByID(orgID)
	if err != nil {
		return nil, err
	}
	return org.Admins, nil
}

// OrganizationUserCount returns the number of users in an organization.
func (c *Client) OrganizationUserCount(orgID int64) (int, error) {
	org, err := c.GetOrganizationByID(orgID)
	if err != nil {
		return 0, err
	}
	return len(org.Users), nil
}

// OrganizationGroupCount returns the number of groups in an organization.
func (c *Client) OrganizationGroupCount(orgID int64) (int, error) {
	org, err := c.GetOrganizationByID(orgID)
	if err != nil {
		return 0, err
	}
	return len(org.Groups), nil
}

// GetOrganizationStats returns statistics for an organization.
func (c *Client) GetOrganizationStats(orgID int64) (map[string]interface{}, error) {
	org, err := c.GetOrganizationByID(orgID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"id":              org.ID,
		"slug":            org.Slug,
		"title":           org.Title,
		"users_count":     len(org.Users),
		"admins_count":    len(org.Admins),
		"groups_count":    len(org.Groups),
		"guests_count":    len(org.Guests),
		"is_member":       org.IsMember,
		"is_admin":        org.IsAdmin,
		"state":           org.State,
		"email_domain":    org.EmailDomain,
		"default_chat_id": org.DefaultChatID,
	}

	return stats, nil
}

// IsOrgMember checks if a user is a member of an organization.
func (c *Client) IsOrgMember(orgID, userID int64) (bool, error) {
	members, err := c.OrganizationMembers(orgID)
	if err != nil {
		return false, err
	}

	for _, id := range members {
		if id == userID {
			return true, nil
		}
	}
	return false, nil
}

// IsOrgAdmin checks if a user is an admin of an organization.
func (c *Client) IsOrgAdmin(orgID, userID int64) (bool, error) {
	admins, err := c.OrganizationAdmins(orgID)
	if err != nil {
		return false, err
	}

	for _, id := range admins {
		if id == userID {
			return true, nil
		}
	}
	return false, nil
}

// GetOrganizationChats returns all chats in an organization.
func (c *Client) GetOrganizationChats(orgID int64) (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	var orgChats []Chat
	for i := range chats.Chats {
		if chats.Chats[i].OrganizationID != nil && *chats.Chats[i].OrganizationID == orgID {
			orgChats = append(orgChats, chats.Chats[i])
		}
	}

	return &ChatsResponse{Chats: orgChats}, nil
}

// PrettyPrintOrgs prints organization information in a readable format.
func (c *Client) PrettyPrintOrgs(orgs *OrgsResponse) string {
	if orgs == nil || len(orgs.Orgs) == 0 {
		return "No organizations found"
	}

	var result strings.Builder
	for i, org := range orgs.Orgs {
		if i > 0 {
			result.WriteString("\n---\n")
		}
		result.WriteString(fmt.Sprintf("Organization #%d:\n", i+1))
		result.WriteString(fmt.Sprintf("  ID: %d\n", org.ID))
		result.WriteString(fmt.Sprintf("  Slug: %s\n", org.Slug))
		result.WriteString(fmt.Sprintf("  Title: %s\n", org.Title))
		if org.Description != "" {
			result.WriteString(fmt.Sprintf("  Description: %s\n", org.Description))
		}
		result.WriteString(fmt.Sprintf("  Users: %d\n", len(org.Users)))
		result.WriteString(fmt.Sprintf("  Admins: %d\n", len(org.Admins)))
		result.WriteString(fmt.Sprintf("  Groups: %d\n", len(org.Groups)))
		result.WriteString(fmt.Sprintf("  Email Domain: %s\n", org.EmailDomain))
		result.WriteString(fmt.Sprintf("  Is Member: %t\n", org.IsMember))
		result.WriteString(fmt.Sprintf("  Is Admin: %t\n", org.IsAdmin))
	}
	return result.String()
}

// GetTopOrgsByUsers returns organizations sorted by number of users.
func (c *Client) GetTopOrgsByUsers(limit int) (*OrgsResponse, error) {
	orgs, err := c.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	// Sort orgs by user count
	for i := 0; i < len(orgs.Orgs)-1; i++ {
		for j := i + 1; j < len(orgs.Orgs); j++ {
			if len(orgs.Orgs[j].Users) > len(orgs.Orgs[i].Users) {
				orgs.Orgs[i], orgs.Orgs[j] = orgs.Orgs[j], orgs.Orgs[i]
			}
		}
	}

	if limit > 0 && limit < len(orgs.Orgs) {
		return &OrgsResponse{Orgs: orgs.Orgs[:limit]}, nil
	}

	return orgs, nil
}
