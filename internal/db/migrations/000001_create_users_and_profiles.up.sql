-- 000001_create_users_and_profiles.up.sql
BEGIN;

CREATE TABLE IF NOT EXISTS md_profile_tags_or_interests (
    tag_interests_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(10) NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATETIME NOT NULL,
    phone_number VARCHAR(20),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS profiles (
    profile_id VARCHAR(10) NOT NULL PRIMARY KEY,
    user_id VARCHAR(10) NOT NULL,
    profile_image_url VARCHAR(255),
    bio VARCHAR(200),
    instagram_url VARCHAR(100),
    facebook_url VARCHAR(100),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS profile_tags_or_interests (
    profile_id VARCHAR(10) NOT NULL,
    tag_interests_id INT NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles(profile_id),
    FOREIGN KEY (tag_interests_id) REFERENCES md_profile_tags_or_interests(tag_interests_id)
);

CREATE TABLE IF NOT EXISTS friendships (
    user_id VARCHAR(10) NOT NULL,
    friend_id VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL, -- e.g., pending, accepted, rejected
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (friend_id) REFERENCES users(user_id)
);

COMMIT;
