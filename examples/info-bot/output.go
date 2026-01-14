package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ivmaks/go-verbosity/verbosity"
)

func printUser(user *verbosity.User, mode string) {
	if user == nil {
		return
	}

	data := map[string]interface{}{
		"id":            user.ID,
		"name":          user.Name,
		"unique_name":   user.UniqueName,
		"info":          user.Info,
		"is_bot":        user.IsBot,
		"is_active":     user.Active,
		"is_deleted":    user.Deleted,
		"organizations": user.Organizations,
		"time_created":  user.TimeCreated.Format(time.RFC3339),
		"time_updated":  user.TimeUpdated.Format(time.RFC3339),
	}

	switch mode {
	case "json":
		printJSONCompact(data)
	case "json-pretty":
		printJSONPretty(data)
	default:
		fmt.Printf("User Information:\n")
		fmt.Printf("  ID: %d\n", user.ID)
		fmt.Printf("  Name: %s\n", user.Name)
		fmt.Printf("  Unique Name: %s\n", user.UniqueName)
		if user.Info != "" {
			fmt.Printf("  Info: %s\n", user.Info)
		}
		fmt.Printf("  Is Bot: %t\n", user.IsBot)
		fmt.Printf("  Active: %t\n", user.Active)
		if len(user.Organizations) > 0 {
			fmt.Printf("  Organizations: %v\n", user.Organizations)
		}
	}
	fmt.Println()
}

func printChat(chat *verbosity.Chat, mode string) {
	if chat == nil {
		return
	}

	data := map[string]interface{}{
		"id":            chat.ID,
		"title":         chat.Title,
		"description":   chat.Description,
		"is_private":    chat.PM,
		"is_read_only":  chat.ReadOnly,
		"is_favorite":   chat.IsFavorite,
		"posts_count":   chat.PostsCount,
		"members_count": len(chat.MemberIDs),
		"admins_count":  len(chat.AdminIDs),
		"member_ids":    chat.MemberIDs,
		"admin_ids":     chat.AdminIDs,
		"time_created":  chat.TimeCreated.Format(time.RFC3339),
		"time_updated":  chat.TimeUpdated.Format(time.RFC3339),
	}

	if chat.OrganizationID != nil {
		data["organization_id"] = *chat.OrganizationID
	}

	switch mode {
	case "json":
		printJSONCompact(data)
	case "json-pretty":
		printJSONPretty(data)
	default:
		fmt.Printf("Chat Information:\n")
		fmt.Printf("  ID: %d\n", chat.ID)
		fmt.Printf("  Title: %s\n", chat.Title)
		if chat.Description != "" {
			fmt.Printf("  Description: %s\n", chat.Description)
		}
		fmt.Printf("  Private: %t\n", chat.PM)
		fmt.Printf("  Read Only: %t\n", chat.ReadOnly)
		fmt.Printf("  Favorite: %t\n", chat.IsFavorite)
		fmt.Printf("  Posts: %d\n", chat.PostsCount)
		fmt.Printf("  Members: %d\n", len(chat.MemberIDs))
		fmt.Printf("  Admins: %d\n", len(chat.AdminIDs))
		if chat.OrganizationID != nil {
			fmt.Printf("  Organization ID: %d\n", *chat.OrganizationID)
		}
	}
	fmt.Println()
}

func printOrg(org *verbosity.Org, mode string) {
	if org == nil {
		return
	}

	data := map[string]interface{}{
		"id":              org.ID,
		"slug":            org.Slug,
		"title":           org.Title,
		"description":     org.Description,
		"email_domain":    org.EmailDomain,
		"users_count":     len(org.Users),
		"admins_count":    len(org.Admins),
		"groups_count":    len(org.Groups),
		"guests_count":    len(org.Guests),
		"is_member":       org.IsMember,
		"is_admin":        org.IsAdmin,
		"state":           org.State,
		"default_chat_id": org.DefaultChatID,
		"time_created":    org.TimeCreated.Format(time.RFC3339),
		"time_updated":    org.TimeUpdated.Format(time.RFC3339),
	}

	switch mode {
	case "json":
		printJSONCompact(data)
	case "json-pretty":
		printJSONPretty(data)
	default:
		fmt.Printf("Organization Information:\n")
		fmt.Printf("  ID: %d\n", org.ID)
		fmt.Printf("  Slug: %s\n", org.Slug)
		fmt.Printf("  Title: %s\n", org.Title)
		if org.Description != "" {
			fmt.Printf("  Description: %s\n", org.Description)
		}
		fmt.Printf("  Email Domain: %s\n", org.EmailDomain)
		fmt.Printf("  Users: %d\n", len(org.Users))
		fmt.Printf("  Admins: %d\n", len(org.Admins))
		fmt.Printf("  Groups: %d\n", len(org.Groups))
		fmt.Printf("  Guests: %d\n", len(org.Guests))
		fmt.Printf("  Is Member: %t\n", org.IsMember)
		fmt.Printf("  Is Admin: %t\n", org.IsAdmin)
		fmt.Printf("  State: %s\n", org.State)
		fmt.Printf("  Default Chat ID: %d\n", org.DefaultChatID)
	}
	fmt.Println()
}

func printChats(chats *verbosity.ChatsResponse, mode string) {
	if chats == nil || len(chats.Chats) == 0 {
		fmt.Println("No chats found")
		return
	}

	switch mode {
	case "json":
		data := map[string]interface{}{
			"total": len(chats.Chats),
			"chats": chats.Chats,
		}
		printJSONCompact(data)
	case "json-pretty":
		data := map[string]interface{}{
			"total": len(chats.Chats),
			"chats": chats.Chats,
		}
		printJSONPretty(data)
	default:
		fmt.Printf("Total chats: %d\n\n", len(chats.Chats))
		for i, chat := range chats.Chats {
			if i > 0 {
				fmt.Println("---")
			}
			fmt.Printf("Chat #%d:\n", i+1)
			fmt.Printf("  ID: %d\n", chat.ID)
			fmt.Printf("  Title: %s\n", chat.Title)
			if chat.Description != "" {
				fmt.Printf("  Description: %s\n", chat.Description)
			}
			fmt.Printf("  Members: %d\n", len(chat.MemberIDs))
			fmt.Printf("  Posts: %d\n", chat.PostsCount)
			fmt.Printf("  Private: %t\n", chat.PM)
		}
	}
	fmt.Println()
}

func printOrgs(orgs *verbosity.OrgsResponse, mode string) {
	if orgs == nil || len(orgs.Orgs) == 0 {
		fmt.Println("No organizations found")
		return
	}

	switch mode {
	case "json":
		data := map[string]interface{}{
			"total": len(orgs.Orgs),
			"orgs":  orgs.Orgs,
		}
		printJSONCompact(data)
	case "json-pretty":
		data := map[string]interface{}{
			"total": len(orgs.Orgs),
			"orgs":  orgs.Orgs,
		}
		printJSONPretty(data)
	default:
		fmt.Printf("Total organizations: %d\n\n", len(orgs.Orgs))
		for i, org := range orgs.Orgs {
			if i > 0 {
				fmt.Println("---")
			}
			fmt.Printf("Organization #%d:\n", i+1)
			fmt.Printf("  ID: %d\n", org.ID)
			fmt.Printf("  Slug: %s\n", org.Slug)
			fmt.Printf("  Title: %s\n", org.Title)
			if org.Description != "" {
				fmt.Printf("  Description: %s\n", org.Description)
			}
			fmt.Printf("  Users: %d\n", len(org.Users))
			fmt.Printf("  Is Member: %t\n", org.IsMember)
			fmt.Printf("  Is Admin: %t\n", org.IsAdmin)
		}
	}
	fmt.Println()
}

func printMemberIDs(title string, ids []int64, mode string) {
	if len(ids) == 0 {
		fmt.Printf("%s: (none)\n", title)
		return
	}

	data := map[string]interface{}{
		"title": title,
		"total": len(ids),
		"ids":   ids,
	}

	switch mode {
	case "json":
		printJSONCompact(data)
	case "json-pretty":
		printJSONPretty(data)
	default:
		fmt.Printf("%s (%d):\n", title, len(ids))
		for i, id := range ids {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%d", id)
		}
		fmt.Println()
	}
	fmt.Println()
}

func printStats(stats map[string]interface{}, mode string) {
	switch mode {
	case "json":
		printJSONCompact(stats)
	case "json-pretty":
		printJSONPretty(stats)
	default:
		for key, value := range stats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	fmt.Println()
}

func printJSONCompact(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}

func printJSONPretty(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
}
