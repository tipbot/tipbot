
CREATE DATABASE tipbot;

CREATE TABLE tip_users (user_id SERIAL, github_name varchar(20), account_id varchar(60), secret_key varchar(60));
// Keeps track of the last comment we processed in any thread
CREATE TABLE processed (thread_id int, owner varchar(20), repo varchar(20), issue_number int, since varchar(30));


// LATER
CREATE TABLE presets (user_id int, preset varchar(20), amount float, asset_code varchar(12));
CREATE TABLE tips_sent (TipID, Owner, Repo, CommentID, From, To, Asset,Amount, When);