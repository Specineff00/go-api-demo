-- Insert or ignore if already exists
-- Without IGNORE if it tries again to insert the same data it CAN crash
-- Only if the schema for users has a UNIQUE attribute assigned to it (which it does)
-- REMEMBER ' not "
INSERT OR IGNORE INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Johnson', 'bob@example.com');