INSERT INTO users (email, password_hash) VALUES ('alice@test.com', 'alice_hash'), ('bob@test.com', 'bob_hash');
INSERT INTO versions (release_date) VALUES (CURRENT_TIMESTAMP);
INSERT INTO boitiers (user_id, version_id) VALUES (1, 1), (2, 1);
INSERT INTO quantities (unit) VALUES ('temperature');

