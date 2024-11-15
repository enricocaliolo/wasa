DROP DATABASE IF EXISTS wasa;

CREATE DATABASE wasa;

USE wasa;

CREATE TABLE User (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    icon VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Conversation (
    conversation_id INT AUTO_INCREMENT PRIMARY KEY,
    user1 INT,
    user2 INT,
    FOREIGN KEY (user1) REFERENCES User(user_id),
    FOREIGN KEY (user2) REFERENCES User(user_id)
);

CREATE TABLE `Group` (
    group_id INT AUTO_INCREMENT PRIMARY KEY,
    group_name VARCHAR(64) NOT NULL,
    icon BLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_user_id INT,
    FOREIGN KEY (created_by_user_id) REFERENCES User (user_id)
);

CREATE TABLE GroupParticipants (
    group_id INT,
    user_id INT,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES `Group`(group_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id)
);

CREATE TABLE Message (
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    content BLOB NOT NULL,
    sent_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_time TIMESTAMP,
    deleted_time TIMESTAMP,
    sender_id INT,
    conversation_id INT,
    group_id INT,
    replied_to INT,
    forwarded_from INT,
    FOREIGN KEY (sender_id) REFERENCES User(user_id),
    FOREIGN KEY (conversation_id) REFERENCES Conversation(conversation_id),
    FOREIGN KEY (group_id) REFERENCES `Group`(group_id),
    FOREIGN KEY (replied_to) REFERENCES `Message`(message_id),
    FOREIGN KEY (forwarded_from) REFERENCES `Group`(message_id),
);

-- users

INSERT INTO User (user_id, username)
VALUES  (1, 'enrico'),
        (2, 'manu'),
        (3, 'teste');

INSERT INTO Conversation (conversation_id, user1, user2)
VALUES (1, 1, 2);

INSERT INTO `GROUP` (group_id, group_name)
VALUES (1, 'grupo');

INSERT INTO GroupParticipants (group_id, user_id)
VALUES  (1, 1),
        (1, 2),
        (1, 3);

INSERT INTO Message (id_message, content, sender_id, conversation_id)
VALUES  (1, 'oi, tudo bem?', 1, 1),
        (2, 'tudo bem, e vc?', 2, 1),
        (3, 'tudo bem tbm', 1, 1);