CREATE TABLE
    IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS tokens (
        token TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS templates (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        repo_name TEXT NOT NULL,
        dockerfile TEXT NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS containers (
        id TEXT PRIMARY KEY,
        docker_id TEXT NOT NULL,
        image_name TEXT NOT NULL,
        container_name TEXT NOT NULL,
        git_repo TEXT,
        user_id TEXT NOT NULL,
        env_vars TEXT,
        ports TEXT,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS ports (
        port INTEGER PRIMARY KEY,
        in_use BOOLEAN NOT NULL,
        container_id TEXT,
        FOREIGN KEY (container_id) REFERENCES containers (id)
    );