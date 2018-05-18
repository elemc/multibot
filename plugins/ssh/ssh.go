package main

type UserCommand struct {
	ID                      int64
	ChatID                  int64
	Host                    string
	Port                    string
	User                    string
	Pass                    string
	Command                 string
	NotificationIfErrorOnly bool
}
