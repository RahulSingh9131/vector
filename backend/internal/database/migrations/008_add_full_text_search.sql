-- Write your migrate up statements here

-- Add search_vector columns
ALTER TABLE issues ADD COLUMN IF NOT EXISTS search_vector TSVECTOR;
ALTER TABLE comments ADD COLUMN IF NOT EXISTS search_vector TSVECTOR;

-- Create GIN indexes
CREATE INDEX IF NOT EXISTS idx_issues_search ON issues USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_comments_search ON comments USING GIN(search_vector);

-- Create trigger functions
CREATE OR REPLACE FUNCTION issues_search_trigger() RETURNS trigger AS $$
begin
  new.search_vector :=
    setweight(to_tsvector('english', coalesce(new.title,'')), 'A') ||
    setweight(to_tsvector('english', coalesce(new.description,'')), 'B');
  return new;
end
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comments_search_trigger() RETURNS trigger AS $$
begin
  new.search_vector :=
    setweight(to_tsvector('english', coalesce(new.body,'')), 'A');
  return new;
end
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS trg_issues_search ON issues;
CREATE TRIGGER trg_issues_search BEFORE INSERT OR UPDATE ON issues
    FOR EACH ROW EXECUTE FUNCTION issues_search_trigger();

DROP TRIGGER IF EXISTS trg_comments_search ON comments;
CREATE TRIGGER trg_comments_search BEFORE INSERT OR UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION comments_search_trigger();

-- Initialize existing data
UPDATE issues SET search_vector = 
    setweight(to_tsvector('english', coalesce(title,'')), 'A') ||
    setweight(to_tsvector('english', coalesce(description,'')), 'B')
WHERE search_vector IS NULL;

UPDATE comments SET search_vector = 
    setweight(to_tsvector('english', coalesce(body,'')), 'A')
WHERE search_vector IS NULL;

---- create above / drop below ----

-- Write your migrate down statements here

DROP TRIGGER IF EXISTS trg_comments_search ON comments;
DROP TRIGGER IF EXISTS trg_issues_search ON issues;

DROP FUNCTION IF EXISTS comments_search_trigger();
DROP FUNCTION IF EXISTS issues_search_trigger();

DROP INDEX IF EXISTS idx_comments_search;
DROP INDEX IF EXISTS idx_issues_search;

ALTER TABLE comments DROP COLUMN IF EXISTS search_vector;
ALTER TABLE issues DROP COLUMN IF EXISTS search_vector;
