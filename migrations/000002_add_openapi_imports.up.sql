CREATE TABLE openapi_imports (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    filename    VARCHAR(255) NOT NULL,
    spec_url    TEXT,
    content     JSONB NOT NULL,
    version     VARCHAR(50),
    imported_by UUID REFERENCES users(id),
    imported_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_openapi_imports_project ON openapi_imports(project_id);
