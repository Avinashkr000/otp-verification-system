# OTP Verification System

A complete full-stack OTP (One-Time Password) verification system built with **Go** backend, **React** frontend, and **MySQL** database.

## 🚀 Features

- **Secure OTP Generation**: 6-digit random OTP with expiry time (5 minutes)
- **Email/Phone Number Verification**: Support for both email and phone verification
- **Rate Limiting**: Prevents spam by limiting OTP requests
- **Modern UI**: Clean and responsive React interface
- **RESTful API**: Well-structured Go backend with proper error handling
- **MySQL Database**: Persistent storage for OTPs and user data
- **Docker Support**: Easy deployment with Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- MySQL 8.0+
- Docker and Docker Compose (optional)

## 🛠️ Tech Stack

### Backend
- **Go** (Golang) - Backend programming language
- **Gin** - Web framework for Go
- **GORM** - ORM library for database operations
- **MySQL** - Relational database
- **JWT** - Token-based authentication (optional)

### Frontend
- **React** - UI library
- **Vite** - Build tool
- **Axios** - HTTP client
- **Tailwind CSS** - Styling
- **React OTP Input** - Custom OTP input component

## 📁 Project Structure

```
otp-verification-system/
├── backend/
│   ├── config/          # Configuration files
│   ├── controllers/     # Request handlers
│   ├── models/          # Database models
│   ├── routes/          # API routes
│   ├── utils/           # Helper functions
│   ├── main.go          # Entry point
│   └── go.mod           # Go dependencies
├── frontend/
│   ├── public/          # Static files
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   ├── App.jsx      # Main app component
│   │   └── main.jsx     # Entry point
│   ├── package.json     # npm dependencies
│   └── vite.config.js   # Vite configuration
├── database/
│   └── schema.sql       # MySQL database schema
├── docker-compose.yml   # Docker compose configuration
└── README.md
```

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
```bash
git clone https://github.com/Avinashkr000/otp-verification-system.git
cd otp-verification-system
```

2. **Start all services**
```bash
docker-compose up -d
```

3. **Access the application**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- MySQL: localhost:3306

### Manual Setup

#### 1. Database Setup

```bash
# Create MySQL database
mysql -u root -p
CREATE DATABASE otp_system;
USE otp_system;
SOURCE database/schema.sql;
```

#### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Update config/config.go with your MySQL credentials
# Run the server
go run main.go
```

Backend will start on `http://localhost:8080`

#### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will start on `http://localhost:5173`

## 📡 API Endpoints

### 1. Generate OTP
```http
POST /api/otp/generate
Content-Type: application/json

{
  "email": "user@example.com",
  "phone": "+919876543210"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "otp_id": "uuid-here",
    "expires_at": "2024-11-26T12:55:00Z"
  }
}
```

### 2. Verify OTP
```http
POST /api/otp/verify
Content-Type: application/json

{
  "otp_id": "uuid-here",
  "otp_code": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP verified successfully",
  "data": {
    "verified": true
  }
}
```

### 3. Resend OTP
```http
POST /api/otp/resend
Content-Type: application/json

{
  "otp_id": "uuid-here"
}
```

## 🔒 Security Features

1. **OTP Expiry**: OTPs expire after 5 minutes
2. **Rate Limiting**: Maximum 3 attempts per email/phone per hour
3. **Single Use**: OTPs can only be used once
4. **Secure Generation**: Cryptographically secure random number generation
5. **Input Validation**: All inputs are validated and sanitized
6. **CORS Protection**: Configured CORS for frontend-backend communication

## 🎨 Frontend Features

- **Responsive Design**: Works on all device sizes
- **Modern UI**: Clean and intuitive interface using Tailwind CSS
- **Custom OTP Input**: Auto-focus and navigation between input fields
- **Real-time Validation**: Instant feedback on OTP verification
- **Loading States**: Visual feedback during API calls
- **Error Handling**: User-friendly error messages

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 📦 Environment Variables

### Backend (.env)
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
```env
VITE_API_URL=http://localhost:8080
```

## 🐳 Docker Configuration

The project includes a complete Docker setup:

- **MySQL Container**: Database with persistent volume
- **Backend Container**: Go application
- **Frontend Container**: React application with Nginx

All services are orchestrated with Docker Compose for easy deployment.

## 📝 Database Schema

```sql
CREATE TABLE otps (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255),
    phone VARCHAR(20),
    otp_code VARCHAR(6) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    attempt_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP NULL
);

CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_phone_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 🚀 Deployment

### Production Deployment

1. **Backend**: Deploy to services like AWS EC2, DigitalOcean, or Heroku
2. **Frontend**: Deploy to Vercel, Netlify, or AWS S3 + CloudFront
3. **Database**: Use managed MySQL services like AWS RDS or DigitalOcean Managed Databases

### Environment-specific Configuration

- Set proper CORS origins in production
- Use environment variables for sensitive data
- Enable HTTPS for all endpoints
- Set up proper logging and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Avinash Kumar**
- GitHub: [@Avinashkr000](https://github.com/Avinashkr000)
- LinkedIn: [Avinash Kumar](https://linkedin.com/in/avinashkr000)

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [React](https://reactjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Vite](https://vitejs.dev/)

## 📞 Support

If you have any questions or need help, please open an issue in the GitHub repository.

---

⭐ **If you found this project helpful, please give it a star!**