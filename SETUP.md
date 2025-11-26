# Setup Guide

Detailed setup instructions for the OTP Verification System.

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software

1. **Go** (version 1.21 or higher)
   ```bash
   # Check Go version
   go version
   
   # Install Go from https://golang.org/dl/
   ```

2. **Node.js** (version 18 or higher) and npm
   ```bash
   # Check Node.js version
   node --version
   npm --version
   
   # Install from https://nodejs.org/
   ```

3. **MySQL** (version 8.0 or higher)
   ```bash
   # Check MySQL version
   mysql --version
   
   # Install from https://dev.mysql.com/downloads/
   ```

4. **Docker** and **Docker Compose** (optional, for containerized setup)
   ```bash
   # Check Docker version
   docker --version
   docker-compose --version
   
   # Install from https://docs.docker.com/get-docker/
   ```

## Setup Methods

Choose one of the following setup methods:

### Method 1: Docker Setup (Recommended)

Easiest way to get started. Everything runs in containers.

#### Steps:

1. **Clone the repository**
   ```bash
   git clone https://github.com/Avinashkr000/otp-verification-system.git
   cd otp-verification-system
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Check if services are running**
   ```bash
   docker-compose ps
   ```

4. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - MySQL: localhost:3306

5. **View logs** (if needed)
   ```bash
   # All services
   docker-compose logs -f
   
   # Specific service
   docker-compose logs -f backend
   ```

6. **Stop services**
   ```bash
   docker-compose down
   ```

### Method 2: Manual Setup

For development or if you prefer running services locally.

#### Step 1: Database Setup

1. **Start MySQL service**
   ```bash
   # On macOS
   brew services start mysql
   
   # On Linux
   sudo systemctl start mysql
   
   # On Windows
   # Start MySQL from Services or MySQL Workbench
   ```

2. **Create database and tables**
   ```bash
   # Login to MySQL
   mysql -u root -p
   
   # Run the schema
   source database/schema.sql
   
   # Or run directly
   mysql -u root -p < database/schema.sql
   ```

3. **Verify database creation**
   ```bash
   mysql -u root -p
   USE otp_system;
   SHOW TABLES;
   ```

#### Step 2: Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Configure database connection**
   
   Edit `backend/config/database.go` and update the following:
   ```go
   config := Config{
       Host:     "localhost",
       Port:     "3306",
       User:     "root",
       Password: "your_mysql_password", // Change this
       DBName:   "otp_system",
   }
   ```

4. **Run the backend server**
   ```bash
   go run main.go
   ```

   You should see:
   ```
   Database connected successfully!
   Database migration completed!
   Server starting on http://localhost:8080
   ```

5. **Test the API**
   ```bash
   curl http://localhost:8080/health
   ```

#### Step 3: Frontend Setup

1. **Open a new terminal and navigate to frontend**
   ```bash
   cd frontend
   ```

2. **Install npm dependencies**
   ```bash
   npm install
   ```

3. **Start the development server**
   ```bash
   npm run dev
   ```

   You should see:
   ```
   VITE v5.4.1  ready in 500 ms
   ➜  Local:   http://localhost:5173/
   ```

4. **Open in browser**
   
   Visit http://localhost:5173

## Testing the Application

### 1. Request OTP

1. Open http://localhost:5173 in your browser
2. Choose Email or Phone
3. Enter a valid email (e.g., `test@example.com`) or phone (e.g., `+919876543210`)
4. Click "Send OTP"

### 2. Check Backend Console

The OTP will be displayed in the backend console:

```
=== OTP Generated ===
OTP ID: 123e4567-e89b-12d3-a456-426614174000
OTP Code: 123456
Email: test@example.com
Expires At: 2024-11-26T13:00:00Z
===================
```

### 3. Verify OTP

1. Enter the 6-digit OTP shown in the console
2. Click "Verify OTP"
3. You should see the success page

## Common Issues and Solutions

### Issue 1: MySQL Connection Failed

**Error:** `Failed to connect to database`

**Solutions:**
- Check if MySQL is running: `mysql -u root -p`
- Verify credentials in `backend/config/database.go`
- Ensure database `otp_system` exists
- Check MySQL port (default: 3306)

### Issue 2: Backend Port Already in Use

**Error:** `bind: address already in use`

**Solutions:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in backend/main.go
```

### Issue 3: Frontend Can't Connect to Backend

**Error:** `Network Error` or `CORS Error`

**Solutions:**
- Ensure backend is running on port 8080
- Check CORS configuration in `backend/main.go`
- Verify proxy configuration in `frontend/vite.config.js`

### Issue 4: npm Install Fails

**Error:** Various npm errors

**Solutions:**
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and package-lock.json
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

### Issue 5: Docker Compose Fails

**Error:** Various Docker errors

**Solutions:**
```bash
# Stop and remove all containers
docker-compose down -v

# Rebuild images
docker-compose build --no-cache

# Start again
docker-compose up -d
```

## Environment Variables

### Backend (.env)

Create `backend/.env` file:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=otp_system
SERVER_PORT=8080
OTP_EXPIRY_MINUTES=5
```

### Frontend (.env)

Create `frontend/.env` file:

```env
VITE_API_URL=http://localhost:8080
```

## Production Deployment

### Backend

1. **Build the binary**
   ```bash
   cd backend
   go build -o bin/otp-backend main.go
   ```

2. **Run the binary**
   ```bash
   ./bin/otp-backend
   ```

### Frontend

1. **Build for production**
   ```bash
   cd frontend
   npm run build
   ```

2. **Serve the dist folder**
   
   Deploy the `dist` folder to:
   - Vercel
   - Netlify
   - AWS S3 + CloudFront
   - Any static hosting service

### Database

For production, use managed database services:
- AWS RDS
- DigitalOcean Managed Databases
- Google Cloud SQL
- Azure Database for MySQL

## Next Steps

1. **Integrate Email Service**
   - Add SendGrid, Mailgun, or AWS SES
   - Update OTP generation to send emails

2. **Integrate SMS Service**
   - Add Twilio, AWS SNS, or similar
   - Update OTP generation to send SMS

3. **Add Authentication**
   - Implement JWT tokens
   - Add user sessions
   - Protect routes

4. **Add Rate Limiting**
   - Implement Redis for rate limiting
   - Add IP-based restrictions

5. **Monitoring & Logging**
   - Add structured logging
   - Set up monitoring (Prometheus, Grafana)
   - Error tracking (Sentry)

## Support

If you encounter any issues:

1. Check this setup guide
2. Review the README.md
3. Open an issue on GitHub
4. Contact: [GitHub](https://github.com/Avinashkr000)

---

Good luck with your OTP verification system! 🚀
