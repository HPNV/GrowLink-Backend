CREATE TABLE IF NOT EXISTS users (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) CHECK (role IN ('student', 'business')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS businesses (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE,
    company_name VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS students (
    uuid UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE,
    university VARCHAR(100)
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
    status VARCHAR(20) CHECK (status IN ('open', 'in_progress', 'completed')) DEFAULT 'open',
    created_by UUID REFERENCES businesses(uuid) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS student_skills (
    student_uuid UUID REFERENCES students(uuid) ON DELETE CASCADE,
    skill_uuid UUID REFERENCES skills(uuid) ON DELETE CASCADE,
    PRIMARY KEY (student_uuid, skill_uuid)
);

CREATE TABLE IF NOT EXISTS project_skills (
    project_uuid UUID REFERENCES projects(uuid) ON DELETE CASCADE,
    skill_uuid UUID REFERENCES skills(uuid) ON DELETE CASCADE,
    PRIMARY KEY (project_uuid, skill_uuid)
);

CREATE TABLE IF NOT EXISTS student_projects (
    student_uuid UUID REFERENCES students(uuid) ON DELETE CASCADE,
    project_uuid UUID REFERENCES projects(uuid) ON DELETE CASCADE,
    PRIMARY KEY (student_uuid, project_uuid)
);



