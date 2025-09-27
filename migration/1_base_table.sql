CREATE TABLE IF NOT EXISTS developers (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skills (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS developer_skills (
    developer_uuid UUID REFERENCES developer(uuid) ON DELETE CASCADE,
    skill_uuid UUID REFERENCES skills(uuid) ON DELETE CASCADE,
    PRIMARY KEY (developer_uuid, skill_uuid)
);

CREATE TABLE IF NOT EXISTS project_skills (
    project_uuid UUID REFERENCES projects(uuid) ON DELETE CASCADE,
    skill_uuid UUID REFERENCES skills(uuid) ON DELETE CASCADE,
    PRIMARY KEY (project_uuid, skill_uuid)
);

CREATE TABLE IF NOT EXISTS developer_projects (
    developer_uuid UUID REFERENCES developer(uuid) ON DELETE CASCADE,
    project_uuid UUID REFERENCES projects(uuid) ON DELETE CASCADE,
    PRIMARY KEY (developer_uuid, project_uuid)
);

