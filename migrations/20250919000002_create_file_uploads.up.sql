CREATE TABLE file_uploads (
    id VARCHAR(27) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    module VARCHAR(50),
    reference_id VARCHAR(27),
    uploaded_by VARCHAR(27) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for common queries
CREATE INDEX idx_file_uploads_module ON file_uploads(module) WHERE module IS NOT NULL;
CREATE INDEX idx_file_uploads_reference ON file_uploads(reference_id) WHERE reference_id IS NOT NULL;
CREATE INDEX idx_file_uploads_uploaded_by ON file_uploads(uploaded_by);
CREATE INDEX idx_file_uploads_public ON file_uploads(is_public) WHERE is_public = true;
