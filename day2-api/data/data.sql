INSERT INTO accounts (id, balance) VALUES (1, 100), (2, 200) ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob') ON CONFLICT (id) DO NOTHING;