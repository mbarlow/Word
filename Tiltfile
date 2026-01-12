# Word - Local Development with Tilt

# Build the word-service
docker_build(
    'word-service',
    '.',
    dockerfile='deploy/Dockerfile',
    live_update=[
        sync('./cmd', '/app/cmd'),
        sync('./internal', '/app/internal'),
        run('go build -o /app/word-service ./cmd/server', trigger=['./cmd', './internal']),
    ],
)

# Deploy using docker-compose
docker_compose('./docker-compose.yml')

# Resource configuration
dc_resource('word-service', labels=['api'])

# Local resource for downloading sources (manual trigger)
local_resource(
    'download-sources',
    cmd='./scripts/download-sources.sh',
    labels=['pipeline'],
    auto_init=False,
)

# Local resource for running the pipeline (manual trigger)
local_resource(
    'run-pipeline',
    cmd='./scripts/pipeline.sh',
    labels=['pipeline'],
    auto_init=False,
    resource_deps=['download-sources'],
)

# Port forwards
# API: http://localhost:8080
# Metrics: http://localhost:8080/metrics
