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

-- users

INSERT INTO User (user_id, username)
VALUES  (1, 'enrico'),
        (2, 'manu'),
        (3, 'teste');

INSERT INTO `Conversation` (conversation_id)
VALUES  (1),
        (2);


INSERT INTO `ConversationParticipants` (conversation_id, user_id)
VALUES  (1, 1),
        (1, 2),
        (1, 3),
        (2, 1),
        (2, 3);

INSERT INTO `Message` (message_id, content, sender_id, conversation_id)
VALUES  (1, 'oi, tudo bem?', 1, 1),
        (2, 'tudo bem, e vc?', 2, 1),
        (3, 'tudo bem tbm', 1, 1);