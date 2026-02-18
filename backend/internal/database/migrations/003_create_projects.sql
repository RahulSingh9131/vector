-- Write your migrate up statements here

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active' NOT NULL,
    identifier VARCHAR(10) NOT NULL,
    issue_counter INT DEFAULT 0 NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE(organization_id, slug),
    UNIQUE(organization_id, identifier)
);

-- Create project_members junction table
CREATE TABLE IF NOT EXISTS project_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member' NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE(project_id, user_id)
);

-- Create issues table
CREATE TABLE IF NOT EXISTS issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    issue_key VARCHAR(50) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'backlog' NOT NULL,
    priority VARCHAR(50) DEFAULT 'none' NOT NULL,
    type VARCHAR(50) DEFAULT 'task' NOT NULL,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    sort_order INT DEFAULT 0 NOT NULL,
    parent_issue_id UUID REFERENCES issues(id) ON DELETE SET NULL,
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE(project_id, issue_key)
);

-- Create indexes for projects
CREATE INDEX IF NOT EXISTS idx_projects_org ON projects(organization_id);
CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(organization_id, slug);
CREATE INDEX IF NOT EXISTS idx_projects_identifier ON projects(organization_id, identifier);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);

-- Create indexes for project_members
CREATE INDEX IF NOT EXISTS idx_project_members_project ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user ON project_members(user_id);
CREATE INDEX IF NOT EXISTS idx_project_members_role ON project_members(role);

-- Create indexes for issues
CREATE INDEX IF NOT EXISTS idx_issues_project ON issues(project_id);
CREATE INDEX IF NOT EXISTS idx_issues_issue_key ON issues(project_id, issue_key);
CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_priority ON issues(priority);
CREATE INDEX IF NOT EXISTS idx_issues_type ON issues(type);
CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee_id);
CREATE INDEX IF NOT EXISTS idx_issues_reporter ON issues(reporter_id);
CREATE INDEX IF NOT EXISTS idx_issues_parent ON issues(parent_issue_id);
CREATE INDEX IF NOT EXISTS idx_issues_sort_order ON issues(project_id, sort_order);

-- Create triggers for updated_at
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_issues_updated_at BEFORE UPDATE ON issues
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

---- create above / drop below ----

-- Write your migrate down statements here

-- Drop triggers
DROP TRIGGER IF EXISTS update_issues_updated_at ON issues;
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;

-- Drop indexes for issues
DROP INDEX IF EXISTS idx_issues_sort_order;
DROP INDEX IF EXISTS idx_issues_parent;
DROP INDEX IF EXISTS idx_issues_reporter;
DROP INDEX IF EXISTS idx_issues_assignee;
DROP INDEX IF EXISTS idx_issues_type;
DROP INDEX IF EXISTS idx_issues_priority;
DROP INDEX IF EXISTS idx_issues_status;
DROP INDEX IF EXISTS idx_issues_issue_key;
DROP INDEX IF EXISTS idx_issues_project;

-- Drop indexes for project_members
DROP INDEX IF EXISTS idx_project_members_role;
DROP INDEX IF EXISTS idx_project_members_user;
DROP INDEX IF EXISTS idx_project_members_project;

-- Drop indexes for projects
DROP INDEX IF EXISTS idx_projects_created_by;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_identifier;
DROP INDEX IF EXISTS idx_projects_slug;
DROP INDEX IF EXISTS idx_projects_org;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
