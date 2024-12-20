/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"wasa/service/database/conversationDB"
	messagesdb "wasa/service/database/messagesDB"
	"wasa/service/database/schema"
	"wasa/service/database/userDB"
	"wasa/service/shared/models"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error

	// user operations
	ValidateUser(id int) bool
	GetAllUsers() []models.User
	GetUser(username string) (int, error)
	CreateUser(username string) (int, error)
	// UpdateProfile(user models.User) bool
	UpdateUsername(user models.User) bool
	UpdatePhoto(user models.User) bool

	// conversation operations
	GetAllConversations(id int) []models.Conversation
	GetMessagesFromConversation(id int) []models.Message
	IsUserInConversation(conversation_id int, user_id int) (bool, error)
	SendMessage(models.Message) (*models.Message, error)
	RemoveUserFromConversation(conversation_id int, user_id int) (bool, error)
	DeleteConversation(conversation_id int) (bool, error)
	CountParticipants(conversation_id int) (int, error)
	IsMessageFromUser(message_id int, user_id int) (bool, error)

	// message operations
	GetMessage(message_id int, conversation_id int) (models.Message, error)
	DeleteMessage(message_id int) (bool, error)
	CommentMessage(user_id int, message_id int, reaction []byte) (bool, error)
	UncommentMessage(reaction_id int) (bool, error)
	IsReactionFromUser(user_id int, reaction_id int) (bool, error)

	UpdateGroupName(conversation_id int, name string) (bool, error)
	UpdateGroupPhoto(conversation_id int, photo string) (bool, error)
	IsGroup(conversation_id int) (bool, error)
	CreateConversation(creator_id int, members []int) (int, error)
	AddGroupMembers(conversation_id int, members []int) error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	schema.CreateTables(db)
	schema.InsertMockValues(db)

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) GetAllUsers() []models.User {
	return userDB.GetAllUsers(db.c)
}

func (db *appdbimpl) GetUser(username string) (int, error) {
	return userDB.GetUser(db.c, username)
}

func (db *appdbimpl) CreateUser(username string) (int, error) {
	return userDB.CreateUser(db.c, username)
}

func (db *appdbimpl) UpdateUsername(user models.User) bool {
	return userDB.UpdateUsername(db.c, user)
}

func (db *appdbimpl) UpdatePhoto(user models.User) bool {
	return userDB.UpdatePhoto(db.c, user)
}

func (db *appdbimpl) ValidateUser(id int) bool {
	return userDB.ValidateUser(db.c, id)
}

func (db *appdbimpl) GetAllConversations(id int) []models.Conversation {
	return conversationDB.GetAllConversations(db.c, id)
}

func (db *appdbimpl) GetMessagesFromConversation(id int) []models.Message {
	return conversationDB.GetMessagesFromConversation(db.c, id)
}

func (db *appdbimpl) IsUserInConversation(conversation_id int, user_id int) (bool, error) {
	return conversationDB.IsUserInConversation(db.c, conversation_id, user_id)
}

func (db *appdbimpl) IsMessageFromUser(message_id int, user_id int) (bool, error) {
	return conversationDB.IsMessageFromUser(db.c, message_id, user_id)
}

func (db *appdbimpl) SendMessage(message models.Message) (*models.Message, error) {
	return conversationDB.SendMessage(db.c, message)
}

func (db *appdbimpl) GetMessage(message_id int, conversation_id int) (models.Message, error) {
	return messagesdb.GetMessage(db.c, message_id, conversation_id)
}

func (db *appdbimpl) RemoveUserFromConversation(conversation_id int, user_id int) (bool, error) {
	return conversationDB.RemoveUserFromConversation(db.c, conversation_id, user_id)
}

func (db *appdbimpl) DeleteConversation(conversation_id int) (bool, error) {
	return conversationDB.DeleteConversation(db.c, conversation_id)
}

func (db *appdbimpl) DeleteMessage(message_id int) (bool, error) {
	return conversationDB.DeleteMessage(db.c, message_id)
}

func (db *appdbimpl) CountParticipants(conversation_id int) (int, error) {
	return conversationDB.CountParticipants(db.c, conversation_id)
}

func (db *appdbimpl) CommentMessage(user_id int, message_id int, reaction []byte) (bool, error) {
	return messagesdb.CommentMessage(db.c, message_id, user_id, reaction)
}

func (db *appdbimpl) UncommentMessage(reaction_id int) (bool, error) {
	return messagesdb.UncommentMessage(db.c, reaction_id)
}

func (db *appdbimpl) IsReactionFromUser(user_id int, reaction_id int) (bool, error) {
	return messagesdb.IsReactionFromUser(db.c, reaction_id, user_id)
}

func (db *appdbimpl) UpdateGroupName(conversation_id int, name string) (bool, error) {
	return conversationDB.UpdateGroupName(db.c, conversation_id, name)
}

func (db *appdbimpl) UpdateGroupPhoto(conversation_id int, photo string) (bool, error) {
	return conversationDB.UpdateGroupPhoto(db.c, conversation_id, photo)
}
func (db *appdbimpl) IsGroup(conversation_id int) (bool, error) {
	return conversationDB.IsGroup(db.c, conversation_id)
}
func (db *appdbimpl) CreateConversation(creator_id int, members []int) (int, error) {
	return conversationDB.CreateConversation(db.c, creator_id, members)
}
func (db *appdbimpl) AddGroupMembers(conversation_id int, members []int) error {
	return conversationDB.AddGroupMembers(db.c, conversation_id, members)
}
