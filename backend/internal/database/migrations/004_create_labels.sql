-- Write your migrate up statements here

-- Create labels table (project-scoped)
CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE(project_id, name)
);

-- Create issue_labels junction table (many-to-many)
CREATE TABLE IF NOT EXISTS issue_labels (
    issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    label_id UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    PRIMARY KEY (issue_id, label_id)
);

-- Indexes for labels
CREATE INDEX IF NOT EXISTS idx_labels_project ON labels(project_id);

-- Indexes for issue_labels
CREATE INDEX IF NOT EXISTS idx_issue_labels_issue ON issue_labels(issue_id);
CREATE INDEX IF NOT EXISTS idx_issue_labels_label ON issue_labels(label_id);

-- Trigger for updated_at
CREATE TRIGGER update_labels_updated_at BEFORE UPDATE ON labels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

---- create above / drop below ----

-- Drop trigger
DROP TRIGGER IF EXISTS update_labels_updated_at ON labels;

-- Drop tables
DROP TABLE IF EXISTS issue_labels;
DROP TABLE IF EXISTS labels;
