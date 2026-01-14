package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ivmaks/go-verbosity/verbosity"
)

// Version программы
const Version = "1.0.0"

const (
	envAPIURL     = "VERBOSITY_API_URL"
	envFileURL    = "VERBOSITY_FILE_URL"
	envAPIToken   = "VERBOSITY_API_TOKEN"
	envOutputMode = "VERBOSITY_OUTPUT_MODE"
)

func main() {
	// Определяем флаги
	apiURL := flag.String("api-url", getEnvDefault(envAPIURL, "https://api.verbosity.io"), "API URL")
	fileURL := flag.String("file-url", getEnvDefault(envFileURL, "https://file.verbosity.io"), "File upload URL")
	token := flag.String("token", os.Getenv(envAPIToken), "API token (or VERBOSITY_API_TOKEN env)")
	outputMode := flag.String("output", getEnvDefault(envOutputMode, "text"), "Output mode: text, json, json-pretty")
	showHelp := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")
	userID := flag.Int64("user-id", 0, "Get info about specific user by ID")
	userName := flag.String("user-name", "", "Get info about specific user by unique name")
	chatID := flag.Int64("chat-id", 0, "Get info about specific chat by ID")
	chatTitle := flag.String("chat-title", "", "Find chat by title")
	orgID := flag.Int64("org-id", 0, "Get info about specific organization by ID")
	orgTitle := flag.String("org-title", "", "Find organization by title")
	listChats := flag.Bool("list-chats", false, "List all available chats")
	listOrgs := flag.Bool("list-orgs", false, "List all organizations")
	listUsers := flag.Bool("list-users", false, "List users (requires user IDs)")
	getChatMembers := flag.Int64("chat-members", 0, "Get member list for a chat")
	getChatAdmins := flag.Int64("chat-admins", 0, "Get admin list for a chat")
	getOrgMembers := flag.Int64("org-members", 0, "Get member list for an organization")
	getOrgAdmins := flag.Int64("org-admins", 0, "Get admin list for an organization")
	topChatsByMembers := flag.Int("top-chats-members", 0, "Show top N chats by members")
	topChatsByPosts := flag.Int("top-chats-posts", 0, "Show top N chats by posts")
	topOrgsByUsers := flag.Int("top-orgs-users", 0, "Show top N organizations by users")
	myChats := flag.Bool("my-chats", false, "Show chats where bot is a member")
	favoriteChats := flag.Bool("favorite-chats", false, "Show favorite chats")
	publicChats := flag.Bool("public-chats", false, "Show public chats")
	privateChats := flag.Bool("private-chats", false, "Show private chats")
	myOrgs := flag.Bool("my-orgs", false, "Show organizations where bot is member")
	adminOrgs := flag.Bool("admin-orgs", false, "Show organizations where bot is admin")
	chatStats := flag.Int64("chat-stats", 0, "Show statistics for a chat")
	orgStats := flag.Int64("org-stats", 0, "Show statistics for an organization")

	// Message sending flags
	sendPrivateID := flag.Int64("send-private", 0, "Send private message to user by ID")
	sendPublicID := flag.Int64("send-public", 0, "Send message to chat by ID")
	messageText := flag.String("message", "", "Message text to send")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Info-Bot for Verbosity API v%s

USAGE:
    %s [OPTIONS]

OPTIONS:
`, Version, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
ENVIRONMENT VARIABLES:
    VERBOSITY_API_URL      API URL (default: https://api.verbosity.io)
    VERBOSITY_FILE_URL     File upload URL (default: https://file.verbosity.io)
    VERBOSITY_API_TOKEN    Bot API token (required)
    VERBOSITY_OUTPUT_MODE  Output mode: text, json, json-pretty (default: text)

EXAMPLES:
    # Get help
    %s -help

    # Show version
    %s -version

    # List all chats
    %s -list-chats -token YOUR_TOKEN

    # List all organizations
    %s -list-orgs -token YOUR_TOKEN

    # Get user by ID
    %s -user-id 123 -token YOUR_TOKEN

    # Get chat by ID
    %s -chat-id 456 -token YOUR_TOKEN

    # Get organization by ID
    %s -org-id 789 -token YOUR_TOKEN

    # Find chat by title
    %s -chat-title "General" -token YOUR_TOKEN

    # Show top 5 chats by members
    %s -top-chats-members 5 -token YOUR_TOKEN

    # Show statistics for a chat
    %s -chat-stats 456 -token YOUR_TOKEN

    # Use custom API URL
    %s -api-url https://custom-api.example.com -token YOUR_TOKEN
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	}

	flag.Parse()

	// Показываем справку если запрошено или нет токена
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("Info-Bot for Verbosity API v%s\n", Version)
		os.Exit(0)
	}

	if *token == "" {
		fmt.Fprintf(os.Stderr, "Error: API token is required. Set VERBOSITY_API_TOKEN environment variable or use -token flag.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Создаем конфигурацию
	config := &verbosity.Config{
		APIURL:   strings.TrimRight(*apiURL, "/"),
		FileURL:  strings.TrimRight(*fileURL, "/"),
		APIToken: *token,
	}

	// Создаем клиент
	client := verbosity.NewClient(config)

	// Проверяем подключение - пытаемся получить список чатов
	_, err := client.GetChatIDs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to connect to Verbosity API: %v\n", err)
		os.Exit(1)
	}

	// Выполняем запрошенные операции
	operations := countOperations(
		*userID, *userName, *chatID, *chatTitle, *orgID, *orgTitle,
		*listChats, *listOrgs, *listUsers, *getChatMembers, *getChatAdmins,
		*getOrgMembers, *getOrgAdmins, *topChatsByMembers, *topChatsByPosts,
		*topOrgsByUsers, *myChats, *favoriteChats, *publicChats, *privateChats,
		*myOrgs, *adminOrgs, *chatStats, *orgStats,
		*sendPrivateID, *sendPublicID, *messageText,
	)

	if operations == 0 {
		fmt.Printf("Info-Bot for Verbosity API v%s\n\n", Version)
		fmt.Println("Available Commands (use -<command> to execute):")
		fmt.Println()
		fmt.Println("User Information:")
		fmt.Println("  -user-id <id>      Get info about specific user by ID")
		fmt.Println("  -user-name <name>  Get info about specific user by unique name")
		fmt.Println("  -list-users        List users (requires user IDs)")
		fmt.Println()
		fmt.Println("Chat Information:")
		fmt.Println("  -list-chats              List all available chats")
		fmt.Println("  -chat-id <id>            Get info about specific chat by ID")
		fmt.Println("  -chat-title <title>      Find chat by title")
		fmt.Println("  -chat-members <id>       Get member list for a chat")
		fmt.Println("  -chat-admins <id>        Get admin list for a chat")
		fmt.Println("  -chat-stats <id>         Show statistics for a chat")
		fmt.Println("  -my-chats                Show chats where bot is a member")
		fmt.Println("  -favorite-chats          Show favorite chats")
		fmt.Println("  -public-chats            Show public (non-private) chats")
		fmt.Println("  -private-chats           Show private chats")
		fmt.Println("  -top-chats-members <n>   Show top N chats by members")
		fmt.Println("  -top-chats-posts <n>     Show top N chats by posts")
		fmt.Println()
		fmt.Println("Organization Information:")
		fmt.Println("  -list-orgs          List all organizations")
		fmt.Println("  -org-id <id>        Get info about specific organization by ID")
		fmt.Println("  -org-title <title>  Find organization by title")
		fmt.Println("  -org-members <id>   Get member list for an organization")
		fmt.Println("  -org-admins <id>    Get admin list for an organization")
		fmt.Println("  -org-stats <id>     Show statistics for an organization")
		fmt.Println("  -my-orgs            Show organizations where bot is member")
		fmt.Println("  -admin-orgs         Show organizations where bot is admin")
		fmt.Println("  -top-orgs-users <n> Show top N organizations by users")
		fmt.Println()
		fmt.Println("Message Sending:")
		fmt.Println("  -send-private <id>  Send private message to user by ID")
		fmt.Println("  -send-public <id>   Send message to chat by ID")
		fmt.Println("  -message <text>     Message text (use with -send-private or -send-public)")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -api-url <url>   API URL (default: https://api.verbosity.io)")
		fmt.Println("  -file-url <url>  File URL (default: https://file.verbosity.io)")
		fmt.Println("  -output <mode>   Output mode: text, json, json-pretty")
		fmt.Println()
		fmt.Println("Environment Variables:")
		fmt.Println("  VERBOSITY_API_URL, VERBOSITY_FILE_URL")
		fmt.Println("  VERBOSITY_API_TOKEN, VERBOSITY_OUTPUT_MODE")
		fmt.Println()
		fmt.Println("For more information, see: https://verbosity.ru/pages/docs/bots.html")
		os.Exit(0)
	}

	// Выполняем операции и собираем результаты
	startTime := time.Now()

	// Выполняем все запрошенные операции
	processOperations(client, config, &OperationsConfig{
		UserID:          *userID,
		UserName:        *userName,
		ChatID:          *chatID,
		ChatTitle:       *chatTitle,
		OrgID:           *orgID,
		OrgTitle:        *orgTitle,
		ListChats:       *listChats,
		ListOrgs:        *listOrgs,
		ListUsers:       *listUsers,
		GetChatMembers:  *getChatMembers,
		GetChatAdmins:   *getChatAdmins,
		GetOrgMembers:   *getOrgMembers,
		GetOrgAdmins:    *getOrgAdmins,
		TopChatsMembers: *topChatsByMembers,
		TopChatsPosts:   *topChatsByPosts,
		TopOrgsUsers:    *topOrgsByUsers,
		MyChats:         *myChats,
		FavoriteChats:   *favoriteChats,
		PublicChats:     *publicChats,
		PrivateChats:    *privateChats,
		MyOrgs:          *myOrgs,
		AdminOrgs:       *adminOrgs,
		ChatStats:       *chatStats,
		OrgStats:        *orgStats,
		SendPrivateID:   *sendPrivateID,
		SendPublicID:    *sendPublicID,
		MessageText:     *messageText,
		OutputMode:      *outputMode,
	})

	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "\nExecuted in %v\n", elapsed)
}

// OperationsConfig holds all operation flags
type OperationsConfig struct {
	UserID          int64
	UserName        string
	ChatID          int64
	ChatTitle       string
	OrgID           int64
	OrgTitle        string
	ListChats       bool
	ListOrgs        bool
	ListUsers       bool
	GetChatMembers  int64
	GetChatAdmins   int64
	GetOrgMembers   int64
	GetOrgAdmins    int64
	TopChatsMembers int
	TopChatsPosts   int
	TopOrgsUsers    int
	MyChats         bool
	FavoriteChats   bool
	PublicChats     bool
	PrivateChats    bool
	MyOrgs          bool
	AdminOrgs       bool
	ChatStats       int64
	OrgStats        int64
	SendPrivateID   int64
	SendPublicID    int64
	MessageText     string
	OutputMode      string
}

func processOperations(client *verbosity.Client, config *verbosity.Config, ops *OperationsConfig) {
	// User operations
	if ops.UserID != 0 {
		user, err := client.GetUserByID(ops.UserID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting user by ID: %v\n", err)
		} else {
			printUser(user, ops.OutputMode)
		}
	}

	if ops.UserName != "" {
		user, err := client.GetUserByUniqueName(ops.UserName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting user by name: %v\n", err)
		} else {
			printUser(user, ops.OutputMode)
		}
	}

	// Chat operations
	if ops.ChatID != 0 {
		chat, err := client.GetChatByID(ops.ChatID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting chat by ID: %v\n", err)
		} else {
			printChat(chat, ops.OutputMode)
		}
	}

	if ops.ChatTitle != "" {
		chat, err := client.FindChatByTitle(ops.ChatTitle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding chat by title: %v\n", err)
		} else {
			printChat(chat, ops.OutputMode)
		}
	}

	if ops.ListChats {
		chats, err := client.GetAllChats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing chats: %v\n", err)
		} else {
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.GetChatMembers != 0 {
		members, err := client.ChatMemberIDs(ops.GetChatMembers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting chat members: %v\n", err)
		} else {
			printMemberIDs("Chat Members", members, ops.OutputMode)
		}
	}

	if ops.GetChatAdmins != 0 {
		admins, err := client.ChatAdminIDs(ops.GetChatAdmins)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting chat admins: %v\n", err)
		} else {
			printMemberIDs("Chat Admins", admins, ops.OutputMode)
		}
	}

	if ops.TopChatsMembers > 0 {
		chats, err := client.GetTopChatsByMembers(ops.TopChatsMembers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting top chats by members: %v\n", err)
		} else {
			fmt.Printf("Top %d chats by members:\n", ops.TopChatsMembers)
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.TopChatsPosts > 0 {
		chats, err := client.GetTopChatsByPosts(ops.TopChatsPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting top chats by posts: %v\n", err)
		} else {
			fmt.Printf("Top %d chats by posts:\n", ops.TopChatsPosts)
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.MyChats {
		chats, err := client.GetMyChats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting my chats: %v\n", err)
		} else {
			fmt.Println("Chats where bot is a member:")
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.FavoriteChats {
		chats, err := client.GetFavoriteChats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting favorite chats: %v\n", err)
		} else {
			fmt.Println("Favorite chats:")
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.PublicChats {
		chats, err := client.GetPublicChats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting public chats: %v\n", err)
		} else {
			fmt.Println("Public chats:")
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.PrivateChats {
		chats, err := client.GetPrivateChats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting private chats: %v\n", err)
		} else {
			fmt.Println("Private chats:")
			printChats(chats, ops.OutputMode)
		}
	}

	if ops.ChatStats != 0 {
		stats, err := client.GetChatStats(ops.ChatStats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting chat stats: %v\n", err)
		} else {
			fmt.Printf("Statistics for chat %d:\n", ops.ChatStats)
			printStats(stats, ops.OutputMode)
		}
	}

	// Organization operations
	if ops.OrgID != 0 {
		org, err := client.GetOrganizationByID(ops.OrgID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting organization by ID: %v\n", err)
		} else {
			printOrg(org, ops.OutputMode)
		}
	}

	if ops.OrgTitle != "" {
		org, err := client.FindOrganizationByTitle(ops.OrgTitle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding organization by title: %v\n", err)
		} else {
			printOrg(org, ops.OutputMode)
		}
	}

	if ops.ListOrgs {
		orgs, err := client.GetAllOrganizations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing organizations: %v\n", err)
		} else {
			printOrgs(orgs, ops.OutputMode)
		}
	}

	if ops.GetOrgMembers != 0 {
		members, err := client.OrganizationMembers(ops.GetOrgMembers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting organization members: %v\n", err)
		} else {
			printMemberIDs("Organization Members", members, ops.OutputMode)
		}
	}

	if ops.GetOrgAdmins != 0 {
		admins, err := client.OrganizationAdmins(ops.GetOrgAdmins)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting organization admins: %v\n", err)
		} else {
			printMemberIDs("Organization Admins", admins, ops.OutputMode)
		}
	}

	if ops.TopOrgsUsers > 0 {
		orgs, err := client.GetTopOrgsByUsers(ops.TopOrgsUsers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting top organizations by users: %v\n", err)
		} else {
			fmt.Printf("Top %d organizations by users:\n", ops.TopOrgsUsers)
			printOrgs(orgs, ops.OutputMode)
		}
	}

	if ops.MyOrgs {
		orgs, err := client.GetMyOrganizations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting my organizations: %v\n", err)
		} else {
			fmt.Println("Organizations where bot is a member:")
			printOrgs(orgs, ops.OutputMode)
		}
	}

	if ops.AdminOrgs {
		orgs, err := client.GetAdminOrganizations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting admin organizations: %v\n", err)
		} else {
			fmt.Println("Organizations where bot is an admin:")
			printOrgs(orgs, ops.OutputMode)
		}
	}

	if ops.OrgStats != 0 {
		stats, err := client.GetOrganizationStats(ops.OrgStats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting organization stats: %v\n", err)
		} else {
			fmt.Printf("Statistics for organization %d:\n", ops.OrgStats)
			printStats(stats, ops.OutputMode)
		}
	}

	// Message sending operations
	if ops.SendPrivateID != 0 && ops.MessageText != "" {
		response, err := client.SendPrivateMessageByID(ops.SendPrivateID, ops.MessageText, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending private message: %v\n", err)
		} else {
			fmt.Printf("Private message sent to user %d:\n", ops.SendPrivateID)
			printPrivateMessageResponse(response, ops.OutputMode)
		}
	}

	if ops.SendPublicID != 0 && ops.MessageText != "" {
		response, err := client.SendMessage(ops.SendPublicID, ops.MessageText, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message to chat: %v\n", err)
		} else {
			fmt.Printf("Message sent to chat %d:\n", ops.SendPublicID)
			printMessageResponse(response, ops.OutputMode)
		}
	}
}

// Вспомогательные функции для вывода

func printMessageResponse(resp *verbosity.MessageResponse, mode string) {
	data := map[string]interface{}{
		"post_no": resp.PostNo,
	}
	printData(data, mode)
}

func printPrivateMessageResponse(resp *verbosity.PrivateMessageResponse, mode string) {
	data := map[string]interface{}{
		"chat_id": resp.ChatID,
		"post_no": resp.PostNo,
	}
	printData(data, mode)
}

func printData(data map[string]interface{}, mode string) {
	switch mode {
	case "json":
		b, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(b))
	case "json-pretty":
		b, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(b))
	default:
		for key, value := range data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

func getEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func countOperations(args ...interface{}) int {
	count := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case int64:
			if v != 0 {
				count++
			}
		case int:
			if v != 0 {
				count++
			}
		case string:
			if v != "" {
				count++
			}
		case bool:
			if v {
				count++
			}
		}
	}
	return count
}
