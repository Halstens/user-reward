CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    reward INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tasks (
    user_id INTEGER REFERENCES users(id),
    task_id INTEGER REFERENCES tasks(id),
    completed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, task_id)
);

CREATE TABLE IF NOT EXISTS completed_tasks (
    user_id INTEGER REFERENCES users(id),
    task_type VARCHAR(50) NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, task_type)
);

CREATE TABLE IF NOT EXISTS referrals (
    referrer_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    referee_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (referee_id), -- Каждый пользователь может иметь только одного реферера
    CHECK (referrer_id != referee_id) -- Нельзя быть реферером самому себе
);

CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);

INSERT INTO tasks (name, reward) VALUES 
('subscribe_telegram', 50),
('follow_twitter', 30),
('referral_signup', 100);

INSERT INTO users (name, balance) VALUES 
('Alice Johnson', 150),
('Bob Smith', 300),
('Charlie Brown', 75)
ON CONFLICT DO NOTHING;

ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

UPDATE users SET password_hash = 'password1' WHERE id = 1;
UPDATE users SET password_hash = 'password1' WHERE id = 2;
UPDATE users SET password_hash = 'password1' WHERE id = 3;  