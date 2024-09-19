# APOD API Service

This service fetches data from the NASA Astronomy Picture of the Day (APOD) API and stores it in a PostgreSQL database.

## Environment Variables

### Database Configuration

- **DB_HOST**: The hostname of the PostgreSQL database. Default is `postgres`.
- **DB_PORT**: The port on which the PostgreSQL database is running. Default is `5432`.
- **DB_USERNAME**: The username for connecting to the PostgreSQL database. Default is `admin`.
- **DB_PASSWORD**: The password for connecting to the PostgreSQL database. Default is `root`.
- **DB_NAME**: The name of the PostgreSQL database. Default is `test_db`.
- **DB_RECONN_RETRY**: Number of times to retry connecting to the database upon failure. Default is `3`.
- **DB_TIME_WAIT_PER_TRY**: Time to wait before retrying a failed database connection. Default is `5s`.

### Server Configuration

- **SERVER_HOST**: The host address on which the server runs. Default is `0.0.0.0`.
- **SERVER_PORT**: The port on which the server listens. Default is `8080`.

### NASA API Key

- **NASA_API_KEY**: The API key for accessing the NASA APOD API. Default is `DEMO_KEY`.

### Worker Configuration

- **WORKER_RUN_TIME**: The time at which the worker fetches new APOD data daily. Example: `19:00`.
- **RUN_FETCHING_ON_START**: Whether to fetch APOD data immediately on service start. Default is `false`.
- **NASA_API_URL**: The URL for the NASA APOD API. Default is `https://api.nasa.gov/planetary/apod`.

## Commands

1. **Build the Application**: Use `make build` to compile the application.
2. **Run the Application**: Use `make run` to start the service.
3. **Run Tests**: Use `make test` to execute tests.
4. **Docker**: Use `make docker-up` and `make docker-down` to manage Docker containers.

## Additional Information

- adjust environment variables as needed for your specific setup using .env file for local launch
- `-env` flag can be used to set the path to .env file, but it's not necessary to provide .env file at all