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
	GetAllUsers() ([]models.User, error)
	GetUser(username string) (models.User, error)
	CreateUser(username string) (models.User, error)
	// UpdateProfile(user models.User) bool
	UpdateUsername(user models.User) (bool, error)
	UpdatePhoto(userId int, imageData []byte) (bool, error)

	// conversation operations
	GetAllConversations(id int) ([]models.Conversation, error)
	GetMessagesFromConversation(id int) ([]models.Message, error)
	IsUserInConversation(conversation_id int, user_id int) (bool, error)
	SendMessage(models.Message) (*models.Message, error)
	ReplyToMessage(models.Message) (*models.Message, error)
	ForwardMessage(models.Message) (*models.Message, error)
	RemoveUserFromConversation(conversation_id int, user_id int) (bool, error)
	DeleteConversation(conversation_id int) (bool, error)
	CountParticipants(conversation_id int) (int, error)
	IsMessageFromUser(message_id int, user_id int) (bool, error)

	// message operations
	GetMessage(message_id int, conversation_id int) (models.Message, error)
	DeleteMessage(message_id int) (models.Message, error)
	CommentMessage(user_id int, message_id int, reaction string) (models.Reaction, error)
	UncommentMessage(reaction_id int) (models.Reaction, error)
	IsReactionFromUser(user_id int, reaction_id int) (bool, error)

	UpdateGroupName(conversation_id int, name string) (bool, error)
	UpdateGroupPhoto(conversation_id int, photo []byte) (bool, error)
	IsGroup(conversation_id int) (bool, error)
	CreateConversation(members []int, name string) (models.Conversation, error)
	AddGroupMembers(conversation_id int, members []int) error

	ConversationExists(conversation_id int) (bool, error)

	MarkMessagesSeen(userID int, messageIDs []int) error
	GetMessageSeenStatus(messageIDs []int) (map[int][]int, error)
	GetConversation(conversation_id int) (*models.Conversation, error)
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

func (db *appdbimpl) GetAllUsers() ([]models.User, error) {
	return userDB.GetAllUsers(db.c)
}

func (db *appdbimpl) GetUser(username string) (models.User, error) {
	return userDB.GetUser(db.c, username)
}

func (db *appdbimpl) CreateUser(username string) (models.User, error) {
	return userDB.CreateUser(db.c, username)
}

func (db *appdbimpl) UpdateUsername(user models.User) (bool, error) {
	return userDB.UpdateUsername(db.c, user)
}

func (db *appdbimpl) UpdatePhoto(userId int, imageData []byte) (bool, error) {
	return userDB.UpdatePhoto(db.c, userId, imageData)
}

func (db *appdbimpl) ValidateUser(id int) bool {
	return userDB.ValidateUser(db.c, id)
}

func (db *appdbimpl) GetAllConversations(id int) ([]models.Conversation, error) {
	return conversationDB.GetAllConversations(db.c, id)
}

func (db *appdbimpl) GetMessagesFromConversation(id int) ([]models.Message, error) {
	return conversationDB.GetMessagesFromConversation(db.c, id)
}

func (db *appdbimpl) IsUserInConversation(user_id int, conversation_id int) (bool, error) {
	return conversationDB.IsUserInConversation(db.c, user_id, conversation_id)
}

func (db *appdbimpl) IsMessageFromUser(message_id int, user_id int) (bool, error) {
	return conversationDB.IsMessageFromUser(db.c, message_id, user_id)
}

func (db *appdbimpl) SendMessage(message models.Message) (*models.Message, error) {
	return conversationDB.SendMessage(db.c, message)
}

func (db *appdbimpl) ReplyToMessage(message models.Message) (*models.Message, error) {
	return conversationDB.ReplyToMessage(db.c, message)
}

func (db *appdbimpl) ForwardMessage(message models.Message) (*models.Message, error) {
	return conversationDB.ForwardMessage(db.c, message)
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

func (db *appdbimpl) DeleteMessage(message_id int) (models.Message, error) {
	return conversationDB.DeleteMessage(db.c, message_id)
}

func (db *appdbimpl) CountParticipants(conversation_id int) (int, error) {
	return conversationDB.CountParticipants(db.c, conversation_id)
}

func (db *appdbimpl) CommentMessage(user_id int, message_id int, reaction string) (models.Reaction, error) {
	return messagesdb.CommentMessage(db.c, user_id, message_id, reaction)
}

func (db *appdbimpl) UncommentMessage(reaction_id int) (models.Reaction, error) {
	return messagesdb.UncommentMessage(db.c, reaction_id)
}

func (db *appdbimpl) IsReactionFromUser(user_id int, reaction_id int) (bool, error) {
	return messagesdb.IsReactionFromUser(db.c, reaction_id, user_id)
}

func (db *appdbimpl) UpdateGroupName(conversation_id int, name string) (bool, error) {
	return conversationDB.UpdateGroupName(db.c, conversation_id, name)
}

func (db *appdbimpl) UpdateGroupPhoto(conversation_id int, photo []byte) (bool, error) {
	return conversationDB.UpdateGroupPhoto(db.c, conversation_id, photo)
}
func (db *appdbimpl) IsGroup(conversation_id int) (bool, error) {
	return conversationDB.IsGroup(db.c, conversation_id)
}
func (db *appdbimpl) CreateConversation(members []int, name string) (models.Conversation, error) {
	return conversationDB.CreateConversation(db.c, members, name)
}
func (db *appdbimpl) AddGroupMembers(conversation_id int, members []int) error {
	return conversationDB.AddGroupMembers(db.c, conversation_id, members)
}
func (db *appdbimpl) ConversationExists(conversation_id int) (bool, error) {
	return conversationDB.ConversationExists(db.c, conversation_id)
}
func (db *appdbimpl) MarkMessagesSeen(userID int, messageIDs []int) error {
	return messagesdb.MarkMessagesSeen(db.c, userID, messageIDs)
}
func (db *appdbimpl) GetMessageSeenStatus(messageIDs []int) (map[int][]int, error) {
	return messagesdb.GetMessageSeenStatus(db.c, messageIDs)
}
func (db *appdbimpl) GetConversation(conversation_id int) (*models.Conversation, error) {
	return conversationDB.GetConversation(db.c, conversation_id)
}
