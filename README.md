# DeepAlert Camera Health Monitoring Service

## Overview
Continuous health monitoring for IP camera RTSP streams. The system reads camera configurations from PostgreSQL, probes each stream on a 60s schedule, and routes status events through Google Cloud Pub/Sub (local emulator) to a consumer service that logs health history.

## Architecture
- **camera-probe (Producer)**: Go microservice that queries a database for camera RTSP URLs, performs health checks (DESCRIBE + frame pull), and publishes events.
- **camera-consumer (Consumer)**: Go microservice that subscribes to health events and persists them to the `camera_health_log` table in PostgreSQL.
- **Pub/Sub Emulator**: Provides local message queueing without GCP credentials.
- **RTSP Simulator (api)**: Simulates 10 IP cameras with various health states (healthy, auth failure, offline).

## Quick Start
1. Ensure Docker and Docker Compose are installed.
2. Run the stack:
   ```bash
   docker-compose up -d --build
   ```
3. To view logs:
   ```bash
   docker-compose logs -f camera-probe camera-consumer
   ```
4. To verify the health log in the database:
   ```bash
   docker-compose exec db psql postgresql://postgres:postgres@localhost:5432/postgres -c "SELECT * FROM camera_health_log ORDER BY checked_at DESC LIMIT 10;"
   ```

## Design Decisions
- **Language**: Go 1.22 for efficient concurrency and low-overhead RTSP probing.
- **Queue**: GCP Pub/Sub for scalability and durability.
- **RTSP Probing**: Uses `gortsplib` for DESCRIBE and SETUP+PLAY sequence to validate both connectivity and actual media delivery.
- **Database**: PostgreSQL with `pgx` for high-performance logging.
- **Infrastructure**: Dockerized environment with platform emulation for cross-architecture compatibility.

## Development Workflow
1. **RTSP Prototyping**: Create a standalone test script using `gortsplib` to verify stream connectivity and frame delivery logic outside the microservice.
2. **Core Logic**: Implement the health check probe based on the working prototype, focusing on the DESCRIBE + SETUP + PLAY sequence.
3. **Containerization**: Use multi-stage Docker builds to keep images small and include only necessary CA certificates and timezone data.
4. **Integration**: Connect services via the Pub/Sub emulator and verify with `docker-compose`.
5. **Production Ready**: Swap the emulator for real GCP Pub/Sub and Cloud SQL by updating environment variables.

## Troubleshooting
- **RTSP URLs**: Camera RTSP URLs in the database must use the internal service name `api` (e.g., `rtsp://api:8554/cam1`) to communicate within the Docker network.
- **Platform Mismatch**: If running on AMD64, the `platform: linux/arm64` settings in `docker-compose.yml` enable QEMU emulation for the ARM-based images provided.
