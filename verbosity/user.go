package verbosity

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetUsersByIDs retrieves user information by their IDs.
//
// API: GET /core/user?ids=11,12,15
func (c *Client) GetUsersByIDs(ids []int64) (*UsersResponse, error) {
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

	req, err := c.newRequest(http.MethodGet, "/core/user", params, nil)
	if err != nil {
		return nil, err
	}

	var response UsersResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUsersByUniqueNames retrieves user information by their unique names.
//
// API: GET /core/user?unames=user0,user1,user2
func (c *Client) GetUsersByUniqueNames(names []string) (*UsersResponse, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("unique names slice cannot be empty")
	}

	params := url.Values{
		"unames": {strings.Join(names, ",")},
	}

	req, err := c.newRequest(http.MethodGet, "/core/user", params, nil)
	if err != nil {
		return nil, err
	}

	var response UsersResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUserByID retrieves a single user by ID.
func (c *Client) GetUserByID(id int64) (*User, error) {
	response, err := c.GetUsersByIDs([]int64{id})
	if err != nil {
		return nil, err
	}

	if len(response.Users) == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &response.Users[0], nil
}

// GetUserByUniqueName retrieves a single user by unique name.
func (c *Client) GetUserByUniqueName(name string) (*User, error) {
	response, err := c.GetUsersByUniqueNames([]string{name})
	if err != nil {
		return nil, err
	}

	if len(response.Users) == 0 {
		return nil, fmt.Errorf("user with unique_name %s not found", name)
	}

	return &response.Users[0], nil
}
