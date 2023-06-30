CREATE TABLE IF NOT EXISTS settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key TEXT NOT NULL UNIQUE,
  value TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE,
  admin BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS public_key (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  public_key TEXT NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL,
  UNIQUE (user_id, public_key),
  CONSTRAINT user_id_fk
    FOREIGN KEY(user_id) REFERENCES user(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS repo (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  project_name TEXT NOT NULL,
  description TEXT NOT NULL,
  private BOOLEAN NOT NULL,
  mirror BOOLEAN NOT NULL,
  hidden BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS collab (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  repo_id INTEGER NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL,
  UNIQUE (user_id, repo_id),
  CONSTRAINT user_id_fk
    FOREIGN KEY(user_id) REFERENCES user(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT repo_id_fk
    FOREIGN KEY(repo_id) REFERENCES repo(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
