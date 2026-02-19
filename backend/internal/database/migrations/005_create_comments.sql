-- Write your migrate up statements here

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    body TEXT NOT NULL,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE SET NULL,
    is_edited BOOLEAN DEFAULT false NOT NULL,
    edited_at TIMESTAMP,
    is_deleted BOOLEAN DEFAULT false NOT NULL,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_comments_issue_created ON comments(issue_id, created_at);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_comments_author ON comments(author_id);

-- Trigger for updated_at
CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

---- create above / drop below ----

-- Drop trigger
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;

-- Drop table
DROP TABLE IF EXISTS comments;
