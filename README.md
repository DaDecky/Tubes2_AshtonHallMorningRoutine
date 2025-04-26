---

## Prerequisites
- [Node.js](https://nodejs.org) and [npm](https://npm.io/) (for local frontend development)
- [Golang](https://go.dev/) (for local backend development)
- [Docker](https://docker.com) & [Docker Compose](https://docs.docker.com/compose/) (for containerized development)

---

## Running the Application

### 1. Local Development
#### **Frontend Only**
1. Navigate to the frontend directory:
   ```bash
   cd src/frontend
   ```
2. Install dependencies using `npm`:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```
4. Access the frontend at [http://localhost:8080](http://localhost:8080).

#### **Backend Only**
1. Navigate to the backend directory:
   ```bash
   cd src/backend
   ```
2. Install dependencies:
   ```bash
   go build ./...
   ```
3. Start the FastAPI application:
   ```bash
   go run .
   ```
4. Access the backend at [http://localhost:8081](http://localhost:8081).

#### **Frontend and Backend**
Run both services simultaneously using the steps outlined above in different terminals.

---