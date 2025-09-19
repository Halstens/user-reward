CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    reward INTEGER NOT NULL
);

CREATE TABLE user_tasks (
    user_id INTEGER REFERENCES users(id),
    task_id INTEGER REFERENCES tasks(id),
    completed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, task_id)
);

CREATE TABLE referrals (
    referrer_id INTEGER REFERENCES users(id),
    referee_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (referee_id)
);

CREATE TABLE completed_tasks (
    user_id INTEGER REFERENCES users(id),
    task_type VARCHAR(50) NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, task_type)
);

CREATE TABLE referrals (
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