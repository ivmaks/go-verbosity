package verbosity

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:            12345,
		Name:          "Test User",
		UniqueName:    "testuser",
		Info:          "Test user info",
		InfoParsed:    "parsed info",
		Deleted:       false,
		Active:        true,
		TimeUpdated:   now,
		TimeCreated:   now,
		IsBot:         false,
		Organizations: []int64{1, 2, 3},
	}

	if user.ID != 12345 {
		t.Errorf("Expected ID to be 12345, got %d", user.ID)
	}

	if user.Name != "Test User" {
		t.Errorf("Expected Name to be 'Test User', got '%s'", user.Name)
	}

	if user.UniqueName != "testuser" {
		t.Errorf("Expected UniqueName to be 'testuser', got '%s'", user.UniqueName)
	}

	if user.Info != "Test user info" {
		t.Errorf("Expected Info to be 'Test user info', got '%s'", user.Info)
	}

	if user.Deleted != false {
		t.Errorf("Expected Deleted to be false, got %v", user.Deleted)
	}

	if user.Active != true {
		t.Errorf("Expected Active to be true, got %v", user.Active)
	}

	if user.IsBot != false {
		t.Errorf("Expected IsBot to be false, got %v", user.IsBot)
	}

	if len(user.Organizations) != 3 {
		t.Errorf("Expected Organizations to have 3 elements, got %d", len(user.Organizations))
	}

	if user.Organizations[0] != 1 {
		t.Errorf("Expected first Organization to be 1, got %d", user.Organizations[0])
	}
}

func TestChat(t *testing.T) {
	now := time.Now()
	chat := &Chat{
		ID:              67890,
		Title:           "Test Chat",
		Description:     "Test chat description",
		PostsLiveTime:   86400,
		TwoStepRequired: false,
		HistoryMode:     "full",
		OrgVisible:      true,
		AllowAPI:        true,
		ReadOnly:        false,
		PostsCount:      100,
		PM:              false,
		E2E:             false,
		TimeCreated:     now,
		TimeUpdated:     now,
		TimeEdited:      now,
		AuthorID:        12345,
		TNew:            false,
		AdmFlag:         false,
		CustomTitle:     "Custom Title",
		IsFavorite:      false,
		TShow:           true,
		UserTimeEdited:  now,
		MemberIDs:       []int64{1, 2, 3},
		AdminIDs:        []int64{1},
		GroupIDs:        []int64{10},
		Guests:          []int64{},
		ThreadUsers:     []int64{1, 2},
		ThreadAdmins:    []int64{1},
		ThreadGroups:    []int64{10},
		LastMsg:         "Last message",
		LastReadPostNo:  50,
		LastMsgAuthorID: 12345,
		LastMsgAuthor:   "Test Author",
		LastMsgBotName:  "",
		LastMsgText:     "Last message text",
	}

	if chat.ID != 67890 {
		t.Errorf("Expected ID to be 67890, got %d", chat.ID)
	}

	if chat.Title != "Test Chat" {
		t.Errorf("Expected Title to be 'Test Chat', got '%s'", chat.Title)
	}

	if chat.Description != "Test chat description" {
		t.Errorf("Expected Description to be 'Test chat description', got '%s'", chat.Description)
	}

	if chat.PostsLiveTime != 86400 {
		t.Errorf("Expected PostsLiveTime to be 86400, got %d", chat.PostsLiveTime)
	}

	if len(chat.MemberIDs) != 3 {
		t.Errorf("Expected MemberIDs to have 3 elements, got %d", len(chat.MemberIDs))
	}

	if chat.MemberIDs[0] != 1 {
		t.Errorf("Expected first MemberID to be 1, got %d", chat.MemberIDs[0])
	}

	if len(chat.AdminIDs) != 1 {
		t.Errorf("Expected AdminIDs to have 1 element, got %d", len(chat.AdminIDs))
	}

	if chat.AdminIDs[0] != 1 {
		t.Errorf("Expected first AdminID to be 1, got %d", chat.AdminIDs[0])
	}

	if chat.LastMsg != "Last message" {
		t.Errorf("Expected LastMsg to be 'Last message', got '%s'", chat.LastMsg)
	}

	if chat.LastReadPostNo != 50 {
		t.Errorf("Expected LastReadPostNo to be 50, got %d", chat.LastReadPostNo)
	}
}

func TestOrg(t *testing.T) {
	now := time.Now()
	org := &Org{
		ID:                11111,
		Slug:              "test-org",
		Title:             "Test Organization",
		Description:       "Test organization description",
		DescriptionParsed: "parsed description",
		EmailDomain:       "test.com",
		TimeCreated:       now,
		TimeUpdated:       now,
		TwoStepRequired:   false,
		DefaultChatID:     12345,
		IsMember:          true,
		IsAdmin:           false,
		State:             "active",
		Users:             []int64{1, 2, 3, 4, 5},
		Admins:            []int64{1, 2},
		Groups:            []int64{10},
		Guests:            []int64{},
	}

	if org.ID != 11111 {
		t.Errorf("Expected ID to be 11111, got %d", org.ID)
	}

	if org.Slug != "test-org" {
		t.Errorf("Expected Slug to be 'test-org', got '%s'", org.Slug)
	}

	if org.Title != "Test Organization" {
		t.Errorf("Expected Title to be 'Test Organization', got '%s'", org.Title)
	}

	if org.Description != "Test organization description" {
		t.Errorf("Expected Description to be 'Test organization description', got '%s'", org.Description)
	}

	if org.EmailDomain != "test.com" {
		t.Errorf("Expected EmailDomain to be 'test.com', got '%s'", org.EmailDomain)
	}

	if org.IsMember != true {
		t.Errorf("Expected IsMember to be true, got %v", org.IsMember)
	}

	if org.IsAdmin != false {
		t.Errorf("Expected IsAdmin to be false, got %v", org.IsAdmin)
	}

	if len(org.Users) != 5 {
		t.Errorf("Expected Users to have 5 elements, got %d", len(org.Users))
	}

	if len(org.Admins) != 2 {
		t.Errorf("Expected Admins to have 2 elements, got %d", len(org.Admins))
	}

	if org.Admins[0] != 1 || org.Admins[1] != 2 {
		t.Errorf("Expected Admins to be [1, 2], got %v", org.Admins)
	}
}

func TestMessageResponse(t *testing.T) {
	response := &MessageResponse{
		PostNo: 67890,
	}

	if response.PostNo != 67890 {
		t.Errorf("Expected PostNo to be 67890, got %d", response.PostNo)
	}
}

func TestPrivateMessageResponse(t *testing.T) {
	response := &PrivateMessageResponse{
		ChatID: 12345,
		PostNo: 67890,
	}

	if response.ChatID != 12345 {
		t.Errorf("Expected ChatID to be 12345, got %d", response.ChatID)
	}

	if response.PostNo != 67890 {
		t.Errorf("Expected PostNo to be 67890, got %d", response.PostNo)
	}
}

func TestFileUploadResponse(t *testing.T) {
	response := &FileUploadResponse{
		GUID: "test-guid-12345",
	}

	if response.GUID != "test-guid-12345" {
		t.Errorf("Expected GUID to be 'test-guid-12345', got '%s'", response.GUID)
	}
}

func TestErrorResponse(t *testing.T) {
	response := &ErrorResponse{
		Code:    "ERROR_CODE",
		Message: "Error message",
	}

	if response.Code != "ERROR_CODE" {
		t.Errorf("Expected Code to be 'ERROR_CODE', got '%s'", response.Code)
	}

	if response.Message != "Error message" {
		t.Errorf("Expected Message to be 'Error message', got '%s'", response.Message)
	}
}

func TestValidationErrorResponse(t *testing.T) {
	response := &ValidationErrorResponse{
		TamtamResponseAPI: true,
		Codes:             map[string]string{"field1": "error1"},
		FieldErrors:       map[string]string{"field1": "field error"},
		Extra:             "extra data",
		Error:             "validation error",
	}

	if response.TamtamResponseAPI != true {
		t.Errorf("Expected TamtamResponseAPI to be true, got %v", response.TamtamResponseAPI)
	}

	if len(response.Codes) != 1 {
		t.Errorf("Expected Codes to have 1 element, got %d", len(response.Codes))
	}

	if len(response.FieldErrors) != 1 {
		t.Errorf("Expected FieldErrors to have 1 element, got %d", len(response.FieldErrors))
	}

	if response.Error != "validation error" {
		t.Errorf("Expected Error to be 'validation error', got '%s'", response.Error)
	}
}

func TestChatsResponse(t *testing.T) {
	chats := &ChatsResponse{
		Chats: []Chat{
			{ID: 1, Title: "Chat 1"},
			{ID: 2, Title: "Chat 2"},
		},
	}

	if len(chats.Chats) != 2 {
		t.Errorf("Expected Chats to have 2 elements, got %d", len(chats.Chats))
	}

	if chats.Chats[0].ID != 1 {
		t.Errorf("Expected first ChatID to be 1, got %d", chats.Chats[0].ID)
	}

	if chats.Chats[1].ID != 2 {
		t.Errorf("Expected second ChatID to be 2, got %d", chats.Chats[1].ID)
	}
}

func TestOrgsResponse(t *testing.T) {
	orgs := &OrgsResponse{
		Orgs: []Org{
			{ID: 1, Title: "Org 1"},
			{ID: 2, Title: "Org 2"},
		},
	}

	if len(orgs.Orgs) != 2 {
		t.Errorf("Expected Orgs to have 2 elements, got %d", len(orgs.Orgs))
	}

	if orgs.Orgs[0].ID != 1 {
		t.Errorf("Expected first OrgID to be 1, got %d", orgs.Orgs[0].ID)
	}

	if orgs.Orgs[1].ID != 2 {
		t.Errorf("Expected second OrgID to be 2, got %d", orgs.Orgs[1].ID)
	}
}

func TestChatSyncResponse(t *testing.T) {
	response := &ChatSyncResponse{
		Chats: []int64{1, 2, 3},
	}

	if len(response.Chats) != 3 {
		t.Errorf("Expected Chats to have 3 elements, got %d", len(response.Chats))
	}

	if response.Chats[0] != 1 {
		t.Errorf("Expected first Chat to be 1, got %d", response.Chats[0])
	}
}

func TestOrgSyncResponse(t *testing.T) {
	response := &OrgSyncResponse{
		IDs: []int64{1, 2, 3},
	}

	if len(response.IDs) != 3 {
		t.Errorf("Expected IDs to have 3 elements, got %d", len(response.IDs))
	}

	if response.IDs[0] != 1 {
		t.Errorf("Expected first ID to be 1, got %d", response.IDs[0])
	}
}

func TestUsersResponse(t *testing.T) {
	response := &UsersResponse{
		Users: []User{
			{ID: 1, Name: "User 1"},
			{ID: 2, Name: "User 2"},
		},
	}

	if len(response.Users) != 2 {
		t.Errorf("Expected Users to have 2 elements, got %d", len(response.Users))
	}

	if response.Users[0].ID != 1 {
		t.Errorf("Expected first User ID to be 1, got %d", response.Users[0].ID)
	}

	if response.Users[1].Name != "User 2" {
		t.Errorf("Expected second User Name to be 'User 2', got '%s'", response.Users[1].Name)
	}
}
