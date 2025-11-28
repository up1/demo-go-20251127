-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Insert initial data
INSERT INTO messages (content) VALUES ('Hello, World!');
INSERT INTO messages (content) VALUES ('Welcome to MySQL!');