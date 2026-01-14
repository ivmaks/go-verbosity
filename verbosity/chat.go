package verbosity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetChatIDs retrieves all available chat IDs.
//
// API: GET /core/chat/sync
func (c *Client) GetChatIDs() (*ChatSyncResponse, error) {
	req, err := c.newRequest(http.MethodGet, "/core/chat/sync", nil, nil)
	if err != nil {
		return nil, err
	}

	var response ChatSyncResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetChatsByIDs retrieves chat information by their IDs.
//
// API: GET /core/chat?ids=11,12,15
func (c *Client) GetChatsByIDs(ids []int64) (*ChatsResponse, error) {
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

	req, err := c.newRequest(http.MethodGet, "/core/chat", params, nil)
	if err != nil {
		return nil, err
	}

	var response ChatsResponse
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetChatByID retrieves a single chat by ID.
func (c *Client) GetChatByID(id int64) (*Chat, error) {
	response, err := c.GetChatsByIDs([]int64{id})
	if err != nil {
		return nil, err
	}

	if len(response.Chats) == 0 {
		return nil, fmt.Errorf("chat with id %d not found", id)
	}

	return &response.Chats[0], nil
}

// GetAllChats retrieves information about all available chats.
// This method first gets all chat IDs, then fetches their details.
func (c *Client) GetAllChats() (*ChatsResponse, error) {
	syncResponse, err := c.GetChatIDs()
	if err != nil {
		return nil, err
	}

	if len(syncResponse.Chats) == 0 {
		return &ChatsResponse{Chats: []Chat{}}, nil
	}

	return c.GetChatsByIDs(syncResponse.Chats)
}

// GetOrCreatePrivateChat retrieves or creates a private chat with a user.
//
// API: POST /core/chat/pm/{user_id}
func (c *Client) GetOrCreatePrivateChat(userID int64) (*Chat, error) {
	path := fmt.Sprintf("/core/chat/pm/%d", userID)
	req, err := c.newRequest(http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var response Chat
	if err := c.do(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ChatMemberIDs returns all member IDs of a chat.
func (c *Client) ChatMemberIDs(chatID int64) ([]int64, error) {
	chat, err := c.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	return chat.MemberIDs, nil
}

// ChatAdminIDs returns all admin IDs of a chat.
func (c *Client) ChatAdminIDs(chatID int64) ([]int64, error) {
	chat, err := c.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	return chat.AdminIDs, nil
}

// FindChatByTitle searches for a chat by title.
func (c *Client) FindChatByTitle(title string) (*Chat, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	for i := range chats.Chats {
		if chats.Chats[i].Title == title {
			return &chats.Chats[i], nil
		}
	}

	return nil, fmt.Errorf("chat with title %q not found", title)
}

// IsChatMember checks if a user is a member of a chat.
func (c *Client) IsChatMember(chatID, userID int64) (bool, error) {
	members, err := c.ChatMemberIDs(chatID)
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

// IsChatAdmin checks if a user is an admin of a chat.
func (c *Client) IsChatAdmin(chatID, userID int64) (bool, error) {
	admins, err := c.ChatAdminIDs(chatID)
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

// GetMyChats filters chats where the current bot is a member.
func (c *Client) GetMyChats() (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	// Get bot's user ID from its info
	botUser, err := c.GetBotInfo()
	if err != nil {
		return nil, err
	}

	var myChats []Chat
	for i := range chats.Chats {
		for _, memberID := range chats.Chats[i].MemberIDs {
			if memberID == botUser.ID {
				myChats = append(myChats, chats.Chats[i])
				break
			}
		}
	}

	return &ChatsResponse{Chats: myChats}, nil
}

// GetFavoriteChats returns chats marked as favorite.
func (c *Client) GetFavoriteChats() (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	var favoriteChats []Chat
	for i := range chats.Chats {
		if chats.Chats[i].IsFavorite {
			favoriteChats = append(favoriteChats, chats.Chats[i])
		}
	}

	return &ChatsResponse{Chats: favoriteChats}, nil
}

// GetPublicChats returns non-private (public) chats.
func (c *Client) GetPublicChats() (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	var publicChats []Chat
	for i := range chats.Chats {
		if !chats.Chats[i].PM {
			publicChats = append(publicChats, chats.Chats[i])
		}
	}

	return &ChatsResponse{Chats: publicChats}, nil
}

// GetPrivateChats returns private chats.
func (c *Client) GetPrivateChats() (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	var privateChats []Chat
	for i := range chats.Chats {
		if chats.Chats[i].PM {
			privateChats = append(privateChats, chats.Chats[i])
		}
	}

	return &ChatsResponse{Chats: privateChats}, nil
}

// PrettyPrintChats prints chat information in a readable format.
func (c *Client) PrettyPrintChats(chats *ChatsResponse) string {
	if chats == nil || len(chats.Chats) == 0 {
		return "No chats found"
	}

	var result strings.Builder
	for i, chat := range chats.Chats {
		if i > 0 {
			result.WriteString("\n---\n")
		}
		result.WriteString(fmt.Sprintf("Chat #%d:\n", i+1))
		result.WriteString(fmt.Sprintf("  ID: %d\n", chat.ID))
		result.WriteString(fmt.Sprintf("  Title: %s\n", chat.Title))
		if chat.Description != "" {
			result.WriteString(fmt.Sprintf("  Description: %s\n", chat.Description))
		}
		result.WriteString(fmt.Sprintf("  Members: %d\n", len(chat.MemberIDs)))
		result.WriteString(fmt.Sprintf("  Admins: %d\n", len(chat.AdminIDs)))
		result.WriteString(fmt.Sprintf("  Posts: %d\n", chat.PostsCount))
		if chat.OrganizationID != nil {
			result.WriteString(fmt.Sprintf("  Organization ID: %d\n", *chat.OrganizationID))
		}
		result.WriteString(fmt.Sprintf("  Private: %t\n", chat.PM))
	}
	return result.String()
}

// UserChats returns all chats where the user is a member.
func (c *Client) UserChats(userID int64) (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	var userChatList []Chat
	for i := range chats.Chats {
		for _, memberID := range chats.Chats[i].MemberIDs {
			if memberID == userID {
				userChatList = append(userChatList, chats.Chats[i])
				break
			}
		}
	}

	return &ChatsResponse{Chats: userChatList}, nil
}

// GetChatStats returns statistics for a chat.
func (c *Client) GetChatStats(chatID int64) (map[string]interface{}, error) {
	chat, err := c.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"id":            chat.ID,
		"title":         chat.Title,
		"posts_count":   chat.PostsCount,
		"members_count": len(chat.MemberIDs),
		"admins_count":  len(chat.AdminIDs),
		"is_private":    chat.PM,
		"read_only":     chat.ReadOnly,
		"is_favorite":   chat.IsFavorite,
	}

	return stats, nil
}

// GetTopChatsByMembers returns chats sorted by number of members.
func (c *Client) GetTopChatsByMembers(limit int) (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	// Sort chats by member count
	for i := 0; i < len(chats.Chats)-1; i++ {
		for j := i + 1; j < len(chats.Chats); j++ {
			if len(chats.Chats[j].MemberIDs) > len(chats.Chats[i].MemberIDs) {
				chats.Chats[i], chats.Chats[j] = chats.Chats[j], chats.Chats[i]
			}
		}
	}

	if limit > 0 && limit < len(chats.Chats) {
		return &ChatsResponse{Chats: chats.Chats[:limit]}, nil
	}

	return chats, nil
}

// GetTopChatsByPosts returns chats sorted by number of posts.
func (c *Client) GetTopChatsByPosts(limit int) (*ChatsResponse, error) {
	chats, err := c.GetAllChats()
	if err != nil {
		return nil, err
	}

	// Sort chats by posts count
	for i := 0; i < len(chats.Chats)-1; i++ {
		for j := i + 1; j < len(chats.Chats); j++ {
			if chats.Chats[j].PostsCount > chats.Chats[i].PostsCount {
				chats.Chats[i], chats.Chats[j] = chats.Chats[j], chats.Chats[i]
			}
		}
	}

	if limit > 0 && limit < len(chats.Chats) {
		return &ChatsResponse{Chats: chats.Chats[:limit]}, nil
	}

	return chats, nil
}

// MarshalJSON for Chat to handleomitempty fields properly
func (c Chat) MarshalJSON() ([]byte, error) {
	type Alias Chat
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&c),
	})
}
