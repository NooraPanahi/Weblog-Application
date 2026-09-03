CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE weblogs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image VARCHAR(500),
    author_id BIGINT NOT NULL,
    privacy VARCHAR(10) NOT NULL DEFAULT 'public',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_weblog_author
        FOREIGN KEY (author_id) REFERENCES users(id),

    CONSTRAINT check_weblog_privacy
        CHECK (privacy IN ('public', 'private'))
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    weblog_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_comment_blog
        FOREIGN KEY (weblog_id) REFERENCES weblogs(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comment_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE weblog_shares (
    id BIGSERIAL PRIMARY KEY,
    weblog_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_share_weblog
        FOREIGN KEY (weblog_id) REFERENCES weblogs(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_share_user 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE,

    CONSTRAINT unique_weblog_share
        UNIQUE (weblog_id, user_id)
)