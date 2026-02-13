-- Drop tables in reverse order (respecting foreign key dependencies)
DROP TABLE IF EXISTS subscriptions;
DROP TYPE IF EXISTS subscription_status;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS techniques;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS fighting_books;
DROP TABLE IF EXISTS sword_masters;
