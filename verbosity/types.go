package verbosity

import "time"

// User represents a Verbosity user.
type User struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	UniqueName    string      `json:"unique_name"`
	Info          string      `json:"info"`
	InfoParsed    interface{} `json:"info_parsed"`
	Deleted       bool        `json:"deleted"`
	Active        bool        `json:"active"`
	TimeUpdated   time.Time   `json:"time_updated"`
	TimeCreated   time.Time   `json:"time_created"`
	IsBot         bool        `json:"is_bot"`
	Organizations []int64     `json:"organizations"`
}

// UsersResponse represents the response for user queries.
type UsersResponse struct {
	Users []User `json:"users"`
}

// Chat represents a Verbosity chat.
type Chat struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	OrganizationID  *int64    `json:"organization_id,omitempty"`
	PostsLiveTime   int       `json:"posts_live_time"`
	TwoStepRequired bool      `json:"two_step_required"`
	HistoryMode     string    `json:"history_mode"`
	OrgVisible      bool      `json:"org_visible"`
	AllowAPI        bool      `json:"allow_api"`
	ReadOnly        bool      `json:"read_only"`
	PostsCount      int       `json:"posts_count"`
	PM              bool      `json:"pm"`
	E2E             bool      `json:"e2e"`
	TimeCreated     time.Time `json:"time_created"`
	TimeUpdated     time.Time `json:"time_updated"`
	TimeEdited      time.Time `json:"time_edited"`
	AuthorID        int64     `json:"author_id"`
	TNew            bool      `json:"tnew"`
	AdmFlag         bool      `json:"adm_flag"`
	CustomTitle     string    `json:"custom_title"`
	IsFavorite      bool      `json:"is_favorite"`
	InviterID       *int64    `json:"inviter_id,omitempty"`
	TShow           bool      `json:"tshow"`
	UserTimeEdited  time.Time `json:"user_time_edited"`
	HistoryStart    *int      `json:"history_start,omitempty"`
	Pinned          []int64   `json:"pinned"`
	MemberIDs       []int64   `json:"member_ids"`
	AdminIDs        []int64   `json:"admin_ids"`
	GroupIDs        []int64   `json:"group_ids"`
	Guests          []int64   `json:"guests"`
	ThreadUsers     []int64   `json:"thread_users"`
	ThreadAdmins    []int64   `json:"thread_admins"`
	ThreadGroups    []int64   `json:"thread_groups"`
	LastMsg         string    `json:"last_msg"`
	LastReadPostNo  int       `json:"last_read_post_no"`
	LastMsgAuthorID int64     `json:"last_msg_author_id"`
	LastMsgAuthor   string    `json:"last_msg_author"`
	LastMsgBotName  string    `json:"last_msg_bot_name"`
	LastMsgText     string    `json:"last_msg_text"`
}

// ChatsResponse represents the response for chat queries.
type ChatsResponse struct {
	Chats []Chat `json:"chats"`
}

// ChatSyncResponse represents the response for chat sync.
type ChatSyncResponse struct {
	Chats []int64 `json:"chats"`
}

// Org represents a Verbosity organization (command).
type Org struct {
	ID                int64       `json:"id"`
	Slug              string      `json:"slug"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	DescriptionParsed interface{} `json:"description_parsed"`
	EmailDomain       string      `json:"email_domain"`
	TimeCreated       time.Time   `json:"time_created"`
	TimeUpdated       time.Time   `json:"time_updated"`
	TwoStepRequired   bool        `json:"two_step_required"`
	DefaultChatID     int64       `json:"default_chat_id"`
	IsMember          bool        `json:"is_member"`
	IsAdmin           bool        `json:"is_admin"`
	State             string      `json:"state"`
	InviterID         *int64      `json:"inviter_id,omitempty"`
	Guests            []int64     `json:"guests"`
	Users             []int64     `json:"users"`
	Admins            []int64     `json:"admins"`
	Groups            []int64     `json:"groups"`
}

// OrgsResponse represents the response for organization queries.
type OrgsResponse struct {
	Orgs []Org `json:"orgs"`
}

// OrgSyncResponse represents the response for organization sync.
type OrgSyncResponse struct {
	IDs []int64 `json:"ids"`
}

// MessageResponse represents the response for message sending.
type MessageResponse struct {
	PostNo int64 `json:"post_no"`
}

// PrivateMessageResponse represents the response for private message sending.
type PrivateMessageResponse struct {
	ChatID int64 `json:"chat_id"`
	PostNo int64 `json:"post_no"`
}

// FileUploadResponse represents the response for file upload.
type FileUploadResponse struct {
	GUID string `json:"guid"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents a validation error.
type ValidationErrorResponse struct {
	TamtamResponseAPI bool              `json:"tamtam_response_api"`
	Codes             map[string]string `json:"codes"`
	FieldErrors       map[string]string `json:"field_errors"`
	Extra             interface{}       `json:"extra,omitempty"`
	Error             string            `json:"error"`
}

// UpdateMessageRequest represents a request to update an existing message.
type UpdateMessageRequest struct {
	Text        string   `json:"text,omitempty"`
	E2E         *bool    `json:"e2e,omitempty"`
	ReplyNo     *int64   `json:"reply_no,omitempty"`
	Quote       *string  `json:"quote,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

// UpdateMessageResponse represents the response for message update.
type UpdateMessageResponse struct {
	UUID    string `json:"uuid"`
	ChatID  int64  `json:"chat_id"`
	PostNo  int64  `json:"post_no"`
	Version *int   `json:"ver,omitempty"`
}
