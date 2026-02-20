-- Write your migrate up statements here

-- Create activities table
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    issue_id UUID REFERENCES issues(id) ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    old_value JSONB,
    new_value JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_activities_issue_created ON activities(issue_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activities_project_created ON activities(project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activities_actor ON activities(actor_id);

---- create above / drop below ----

-- Drop table
DROP TABLE IF EXISTS activities;
