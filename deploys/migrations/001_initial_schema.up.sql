CREATE TYPE deploy_status AS ENUM ('in_progress', 'succeeded', 'failed');

CREATE TABLE deploys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    service_name VARCHAR(255) NOT NULL,
    status deploy_status NOT NULL
);
