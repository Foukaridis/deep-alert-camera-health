Changes Made
1. RTSP Probe Refactoring
The RTSP health check logic in 
rtsp.go
 was rewritten to use the gortsp library. This significantly simplified the code, replacing manual RTSP handshake and packet reading logic with a single call to rtsp.HealthCheck.

2. Dependency Update
Updated 
go.mod
 to include github.com/Foukaridis/gortsp v0.1.0 and removed the now-unused gortsplib and pion/rtp dependencies.

3. Integration Testing
Created a new test suite 
health_results_test.go
 that verifies the health status of 10 specific camera scenarios:

Cameras 1-4: Verified as HEALTHY when correct credentials (admin:secret) are provided.
Camera 5: Identified as flaky (toggles between HEALTHY and OFFLINE).
Cameras 6-8: Identified as UNAUTHENTICATED due to incorrect credentials.
Cameras 9-10: Identified as OFFLINE due to incorrect configuration details (missing paths/auth).
Verification Results
The tests were executed against the live api (camera-simulator) container.

bash
cd services/camera-probe
go test -v ./probe/
Highlights:
Status Mapping: The library correctly maps RTSP statuses to our internal model error categories:
rtsp.Unauthenticated -> model.ErrAuthFailed
rtsp.Offline -> model.ErrConnectionRefused
rtsp.Unhealthy -> model.ErrNoFrames
Host Mapping: Confirmed that the API_HOST environment variable correctly redirects probes to the Docker service name when running in a containerized environment.
NOTE

During local testing from the host, Ensure localhost:8554 is reachable. In some environments, switching to 127.0.0.1:8554 in 
health_results_test.go
 resolves IPv6/IPv4 preference issues.