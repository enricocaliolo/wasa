DROP DATABASE IF EXISTS wasa;

CREATE DATABASE wasa;

USE wasa;

CREATE TABLE `User` (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    icon VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `Conversation` (
    conversation_id INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(64),
    is_group BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `ConversationParticipants` (
    conversation_id INT,
    user_id INT,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES `Conversation`(conversation_id),
    FOREIGN KEY (user_id) REFERENCES `User`(user_id)
);

CREATE TABLE `Message` (
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    content BLOB NOT NULL,
    content_type VARCHAR(32) NOT NULL DEFAULT 'text',
    sent_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_time TIMESTAMP,
    deleted_time TIMESTAMP,
    sender_id INT,
    conversation_id INT,
    replied_to INT,
    forwarded_from INT,
    FOREIGN KEY (sender_id) REFERENCES `User`(user_id),
    FOREIGN KEY (conversation_id) REFERENCES `Conversation`(conversation_id),
    FOREIGN KEY (replied_to) REFERENCES `Message`(message_id),
    FOREIGN KEY (forwarded_from) REFERENCES `Message`(message_id)
);

CREATE TABLE `Reactions` (
    reaction_id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL,
    user_id INT NOT NULL,
    reaction BLOB NOT NULL,
    FOREIGN KEY (message_id) REFERENCES `Message`(message_id),
    FOREIGN KEY (user_id) REFERENCES `User`(user_id)
);

INSERT INTO `User` (username, icon) VALUES
		('enrico', 'default_1.png'),
		('jane_smith', 'profile_2.png'),
		('bob_wilson', 'avatar_3.png'),
		('alice_green', 'user_4.png'),
		('charlie_brown', 'pic_5.png'),
		('diana_prince', 'icon_6.png'),
		('bruce_wayne', 'bat_7.png'),
		('peter_parker', 'spider_8.png'),
		('tony_stark', 'iron_9.png'),
		('steve_rogers', 'cap_10.png');

		INSERT INTO `Conversation` (name, is_group) VALUES
		(NULL, FALSE),
		(NULL, FALSE),
		('Project Team', TRUE),
		('Family Group', TRUE),
		(NULL, FALSE),
		('Gaming Squad', TRUE),
		('Book Club', TRUE),
		(NULL, FALSE),
		('Tech Team', TRUE),
		(NULL, FALSE);

		INSERT INTO `ConversationParticipants` (conversation_id, user_id) VALUES
		(1, 1), (1, 2),
		(2, 3), (2, 4),
		(3, 1), (3, 2), (3, 3), (3, 4),
		(4, 5), (4, 6), (4, 7), (4, 8),
		(5, 5), (5, 6),
		(6, 7), (6, 8), (6, 9), (6, 10),
		(7, 1), (7, 3), (7, 5), (7, 7),
		(8, 8), (8, 9),
		(9, 8), (9, 9), (9, 10),
		(10, 7), (10, 10);

		INSERT INTO `Message` (content, content_type, sender_id, conversation_id, replied_to) VALUES
		('Hey, how are you?', 'text', 1, 1, NULL),
		('Im doing great, thanks for asking!', 'text', 2, 1, 1),
		('Project meeting at 3 PM tomorrow', 'text', 1, 3, NULL),
		('Ill be there', 'text', 2, 3, 3),
		('Can anyone help with the database issue?', 'text', 3, 3, NULL),
		('Sure, whats the problem?', 'text', 4, 3, 5),
		('Dinner at 7?', 'text', 5, 4, NULL),
		('Perfect timing!', 'text', 6, 4, 7),
		('New game release this weekend', 'text', 7, 6, NULL),
		('Count me in for the launch', 'text', 8, 6, 9),
		('Has everyone read chapter 5?', 'text', 1, 7, NULL),
		('Yes, great plot twist!', 'text', 3, 7, 11),
		('Project deadline extended to next week', 'text', 8, 9, NULL),
		('Thanks for the update', 'text', 9, 9, 13),
		('Meeting notes from today', 'text', 10, 9, NULL),
		('New feature deployment tonight', 'text', 7, 10, NULL),
		('All systems ready', 'text', 10, 10, 16),
		('Remember to update the documentation', 'text', 8, 9, NULL),
		('Already on it', 'text', 9, 9, 18),
		('Good morning team!', 'text', 1, 3, NULL);