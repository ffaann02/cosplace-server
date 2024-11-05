-- 000001_create_users_and_profiles.down.sql
BEGIN;

DROP TABLE IF EXISTS friendships;

DROP TABLE IF EXISTS profile_tags_or_interests;

DROP TABLE IF EXISTS profiles;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS md_profile_tags_or_interests;

COMMIT;