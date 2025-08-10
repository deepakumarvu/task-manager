package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const apiBase = "http://localhost:8080"

type Session struct {
	UserID string
}

func main() {
	sess := &Session{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Task CLI. Type 'help' for commands.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		args := strings.Split(line, " ")
		if len(args) == 0 || args[0] == "" {
			continue
		}
		switch args[0] {
		case "help":
			printHelp()
		case "exit", "quit":
			return
		case "create-user":
			createUser()
		case "set-user":
			if len(args) < 2 {
				fmt.Println("Usage: set-user <user_id>")
				continue
			}
			sess.UserID = args[1]
			fmt.Println("Session user set to:", sess.UserID)
		case "get-user":
			getUser(sess.UserID)
		case "list-users":
			listUsers()
		case "delete-user":
			deleteUser(sess.UserID)
		case "create-task":
			if sess.UserID == "" {
				fmt.Println("Set user first with: set-user <user_id>")
				continue
			}
			createTask(sess.UserID)
		case "get-task":
			if len(args) < 2 {
				fmt.Println("Usage: get-task <task_id>")
				continue
			}
			getTask(sess.UserID, args[1])
		case "list-tasks":
			listTasks(sess.UserID)
		case "update-task":
			if len(args) < 2 {
				fmt.Println("Usage: update-task <task_id>")
				continue
			}
			updateTask(sess.UserID, args[1])
		case "delete-task":
			if len(args) < 2 {
				fmt.Println("Usage: delete-task <task_id>")
				continue
			}
			deleteTask(sess.UserID, args[1])
		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}

func printHelp() {
	fmt.Println(`
Commands:
  create-user           - Create a new user (prompts for name/email)
  set-user <user_id>    - Set current session user
  get-user              - Get current user info
  list-users            - List all users
  delete-user           - Delete current user
  create-task           - Create a new task (prompts for details)
  get-task <task_id>    - Get a task by ID
  list-tasks            - List tasks for current user
  update-task <task_id> - Update a task (prompts for details)
  delete-task <task_id> - Delete a task
  help                  - Show this help
  exit/quit             - Exit CLI`)
}

func prompt(fields ...string) map[string]string {
	result := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)
	for _, f := range fields {
		fmt.Printf("%s: ", f)
		val, _ := reader.ReadString('\n')
		result[f] = strings.TrimSpace(val)
	}
	return result
}

func createUser() {
	input := prompt("name", "email")
	body, _ := json.Marshal(input)
	resp, err := http.Post(apiBase+"/users", "application/json", bytes.NewBuffer(body))
	handleResp(resp, err)
}

func getUser(userID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	url := fmt.Sprintf("%s/users/%s", apiBase, userID)
	resp, err := http.Get(url)
	handleResp(resp, err)
}

func listUsers() {
	resp, err := http.Get(apiBase + "/users")
	handleResp(resp, err)
}

func deleteUser(userID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", apiBase, userID), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	handleResp(resp, err)
}

func createTask(userID string) {
	input := prompt("title", "description", "due_date", "status")
	body, _ := json.Marshal(input)
	url := fmt.Sprintf("%s/users/%s/tasks", apiBase, userID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	handleResp(resp, err)
}

func getTask(userID, taskID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	url := fmt.Sprintf("%s/users/%s/tasks/%s", apiBase, userID, taskID)
	resp, err := http.Get(url)
	handleResp(resp, err)
}

func listTasks(userID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	url := fmt.Sprintf("%s/users/%s/tasks", apiBase, userID)
	resp, err := http.Get(url)
	handleResp(resp, err)
}

func updateTask(userID, taskID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	input := prompt("title", "description", "due_date", "status")
	body, _ := json.Marshal(input)
	url := fmt.Sprintf("%s/users/%s/tasks/%s", apiBase, userID, taskID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	handleResp(resp, err)
}

func deleteTask(userID, taskID string) {
	if userID == "" {
		fmt.Println("Set user first with: set-user <user_id>")
		return
	}
	url := fmt.Sprintf("%s/users/%s/tasks/%s", apiBase, userID, taskID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	handleResp(resp, err)
}

func handleResp(resp *http.Response, err error) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	var pretty bytes.Buffer
	if json.Valid(out) {
		if err := json.Indent(&pretty, out, "", "  "); err == nil {
			fmt.Println(pretty.String())
			return
		}
	}
	fmt.Println(string(out))
}
