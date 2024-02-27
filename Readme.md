# **Attendance Management System (AMS) - GoLang**
This is an **Attendance Management System (AMS)** written in **GoLang**. It is a backend system responsible for managing attendance records. This project utilizes various third-party modules for different functionalities.

#### **Modules Used**
- github.com/caarlos0/env
- github.com/go-pg/pg
- github.com/golang-jwt/jwt
- github.com/google/subcommands
- github.com/google/uuid
- github.com/google/wire
- github.com/gorilla/mux
- github.com/jinzhu/inflection
- github.com/joho/godotenv
- github.com/pmezard/go-difflib
- github.com/rs/cors
- go.uber.org/multierr
- go.uber.org/zap
- golang.org/x/crypto
- golang.org/x/mod
- golang.org/x/tools
- mellium.im/sasl

### **Installation and Setup**
Clone the repository:

```bash
git clone https://github.com/raunit-verma/Attendance-Go.git
```
Navigate to the project directory:
```bash
cd  Attendance-Go  
```
Set up environment variables:

```bash
DB_USER=postgres
DB_ADDRESS=localhost:5432
DB_PASSWORD=1234
DB_DATABASE=attendance

TYPE=Development # or TYPE=Production

DB_USER_PRODUCTION=
DB_ADDRESS_PRODUCTION=
DB_PASSWORD_PRODUCTION=
DB_DATABASE_PRODUCTION=

PORT=1025
URL=http://localhost:3000

DOMAIN=localhost
JWT_KEY=Devtron@Raunit

PRINCIPAL_PASSWORD= # Do not use plain password, use hased password using bcrypt
```
Download all modules:
```bash
go mod download
```

Run wire to generate wire gen.go:

```bash
go run -mod=mod github.com/google/wire/cmd/wire
```
Build the server:
```bash
go build -o server .
```
Run the server:

```bash
./server
```
### **Docker**
You can also run the application using Docker. A Dockerfile is provided in the repository. In case of running docker image don't forget to add PSQL in same image or another image.

Build the Docker image:

```bash
docker build -t attendance-management-system .
```
Run the Docker container:

```bash
docker run -p 1025:1025 attendance-management-system
```
#### Frontend
The frontend code for this project can be found [here](https://github.com/raunit-verma/Attendance-React "here"). You can integrate it with this backend to create a complete Attendance Management System.

#### Docker Hub Image
An image of this backend is available on Docker Hub [here](https://hub.docker.com/repository/docker/raunitverma/attendance-project-go/general "here"). You can pull this image and run it as a Docker container.

##### Contribution
Contributions are welcome. Feel free to open issues and pull requests.

##### License
This project is licensed under the MIT License.