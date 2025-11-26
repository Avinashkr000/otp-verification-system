-- Create database
CREATE DATABASE IF NOT EXISTS otp_system;
USE otp_system;

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS otps;
DROP TABLE IF EXISTS users;

-- Create OTPs table
CREATE TABLE otps (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) DEFAULT NULL,
    phone VARCHAR(20) DEFAULT NULL,
    otp_code VARCHAR(6) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    attempt_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP NULL,
    
    -- Indexes for better query performance
    INDEX idx_email (email),
    INDEX idx_phone (phone),
    INDEX idx_created_at (created_at),
    INDEX idx_is_verified (is_verified)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE DEFAULT NULL,
    phone VARCHAR(20) UNIQUE DEFAULT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_phone_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_email (email),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert some sample data for testing (optional)
-- INSERT INTO users (id, email, phone, is_email_verified, is_phone_verified) 
-- VALUES 
--     (UUID(), 'test@example.com', '+919876543210', TRUE, FALSE),
--     (UUID(), 'demo@example.com', '+919876543211', FALSE, TRUE);

SELECT 'Database schema created successfully!' AS message;
