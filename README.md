# URL Shortener Project

This project implements a URL shortener API using Go and MySQL. It allows you to shorten long URLs and retrieve the original URL using the generated short URL.

## Prerequisites

Before running the project, ensure that you have the following prerequisites installed on your system:

- Docker: [Installation Guide](https://docs.docker.com/get-docker/)
- Docker Compose: [Installation Guide](https://docs.docker.com/compose/install/)

## Getting Started

Follow the steps below to set up and run the project.

### Step 1: Clone the Repository

Clone the project repository to your local machine:


git clone https://github.com/LoluPapi/altschcapstone

Step 2: Navigate to the Project Directory

Navigate to the project directory:

Step 3: Configure Environment Variables
Create a .env file in the project directory and provide the required environment variables. The file should have the following content:
DB_USER=user
DB_PASSWORD=password
DB_HOST=db
DB_PORT=3306
DB_NAME=short_urls
Replace the values with your desired username, password, and database name.

Step 4: Build and Run the Docker Containers
Build and run the Docker containers using Docker Compose:
docker-compose up --build

Docker Compose will build the necessary containers and start the project services.

### Step 5: Access the Application

Once the Docker containers are up and running, you can access the application using the following URL:
http://localhost:8080


The application provides the following endpoints:

- `/api/shorten` - Shorten a long URL by sending a POST request with a JSON body containing the `longUrl` field.
- `/health` - Check the health status of the application and database connection.

curl -iX POST -H "Content-Type: application/json" -d '{"longUrl": "https://docs.imperva.se-reference-guide/page/abnormally_long_url.htm"}' http://localhost:8080/api/shorten
### Step 6: Clean Up

To stop and remove the Docker containers, use the following command:

docker-compose down


This will stop and remove the containers, networks, and volumes associated with the project.

## Conclusion

You have successfully set up and run the URL shortener project. Follow the provided steps to access the application and explore its functionalities.
