# Twilio SMS Integration Guide

Complete step-by-step guide to integrate Twilio SMS for OTP delivery.

## 📱 Why Twilio?

- **Reliable**: 99.95% uptime SLA
- **Global**: Send SMS to 180+ countries  
- **Free Trial**: $15 credit to get started
- **Easy Integration**: Simple REST API
- **Scalable**: Handles millions of messages

---

## 🚀 Step 1: Create Twilio Account

### 1.1 Sign Up

1. Visit: [https://www.twilio.com/try-twilio](https://www.twilio.com/try-twilio)
2. Fill in your details (Name, Email, Password)
3. Click **Start your free trial**

### 1.2 Verify Your Email

1. Check your email inbox
2. Click verification link from Twilio
3. Complete email verification

### 1.3 Verify Your Phone Number

1. Enter your phone number (with country code)
2. Choose verification method (SMS or Call)
3. Enter the 6-digit verification code
4. Click **Submit**

---

## 🔑 Step 2: Get Your Credentials

### 2.1 Access Twilio Console

1. After verification, visit: [https://console.twilio.com/](https://console.twilio.com/)
2. You'll see your Dashboard with Account Info

### 2.2 Copy Your Credentials

On the Dashboard, find:

- **Account SID**: Starts with `AC` (34 characters)
- **Auth Token**: Click **Show** to reveal (32 characters)

> ⚠️ **Keep these secret!** Don't commit to Git or share publicly.

---

## 📞 Step 3: Get a Phone Number

### 3.1 Get Trial Number

1. In Twilio Console, click **Get a Trial Number**
2. Twilio assigns you a free phone number
3. Click **Choose this Number**
4. Save this number (format: +1234567890)

### 3.2 Trial Limitations

**Free Trial:**
- ✅ Send to verified numbers only
- ✅ Free $15 credit
- ❌ Messages include "trial account" text
- ❌ Can't send to unverified numbers

**To send to any number**: Upgrade account

### 3.3 Verify Test Numbers

1. Go to **Phone Numbers** → **Verified Caller IDs**
2. Click **Add a new Caller ID**
3. Enter phone number
4. Verify via SMS/Call
5. Now you can send OTPs to this number!

---

## ⚙️ Step 4: Configure Application

### 4.1 Create `.env` File

```bash
cd backend
cp .env.example .env
```

### 4.2 Add Twilio Credentials

Edit `backend/.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=otp_system

# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Twilio - Replace with YOUR values from Twilio Console
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_actual_auth_token_here
TWILIO_PHONE_NUMBER=+15551234567
```

**Replace:**
- `TWILIO_ACCOUNT_SID`: Your actual SID from dashboard
- `TWILIO_AUTH_TOKEN`: Your actual token from dashboard  
- `TWILIO_PHONE_NUMBER`: Your Twilio number with `+` and country code

### 4.3 Pull Latest Code

```bash
git pull origin main
```

### 4.4 Restart Backend

```bash
cd backend
go run main.go
```

You should see: `Server starting on http://localhost:8080`

---

## 🧪 Step 5: Test SMS Delivery

### 5.1 Using Frontend

1. **Start Frontend**:
   ```bash
   cd frontend
   npm run dev
   ```

2. **Open**: http://localhost:5173

3. **Click Phone Tab**

4. **Enter Verified Number**: 
   - Must be verified in Twilio Console
   - Include country code: `+919876543210` (India)

5. **Click Send OTP**

6. **Check Phone**: SMS should arrive within seconds!

### 5.2 Using API (cURL)

```bash
curl -X POST http://localhost:8080/api/otp/generate \
  -H "Content-Type: application/json" \
  -d '{"phone": "+919876543210"}'
```

**Response:**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "otp_id": "uuid-here",
    "expires_at": "2024-11-26T13:30:00Z",
    "otp_code": "123456"
  }
}
```

### 5.3 Check Backend Logs

```
=== OTP Generated ===
OTP Code: 123456
Phone: +919876543210
===================

SMS sent successfully! SID: SMxxxx, Status: queued
```

---

## 📊 Step 6: Monitor in Twilio Console

### View Message Logs

1. **Monitor** → **Logs** → **Messaging**
2. See all SMS with status:
   - ✅ **Delivered**: Success!
   - ⏳ **Queued/Sent**: In progress
   - ❌ **Failed/Undelivered**: Check error

### Check Usage

1. **Monitor** → **Usage**
2. View:
   - Messages sent
   - Remaining credit
   - Cost per message

---

## 💰 Step 7: Upgrade (Production)

### When to Upgrade?

- Send to any number
- Remove "trial" message
- Production use
- High volume

### How to Upgrade?

1. **Billing** → **Upgrade**
2. Add payment method
3. Add credit ($20 min recommended)
4. Done!

### Pricing Examples

**India**:
- SMS: ₹0.50-₹1.50/message
- Number: ~₹80/month

**USA**:
- SMS: $0.0079/message
- Number: $1.15/month

---

## 🔒 Security Best Practices

### Protect Credentials

✅ **DO**:
- Store in `.env` file
- Add `.env` to `.gitignore`
- Use environment variables
- Rotate tokens periodically

❌ **DON'T**:
- Commit to Git
- Hardcode in source
- Share publicly
- Expose in frontend

### Rate Limiting

Already implemented:
- Max 3 OTP/hour per phone
- Prevents spam
- Saves credits

---

## 🐛 Troubleshooting

### SMS Not Received?

**Check:**
1. Phone verified in Twilio? (trial accounts)
2. Correct number format? (+91 for India)
3. Backend logs for errors?
4. Twilio Console → Logs for delivery status?
5. Wait 30-60 seconds (carrier delay)

### "Twilio not configured"?

**Fix:**
1. `.env` file exists in `backend/`?
2. All 3 variables set?
3. Restart backend server

### Authentication Error?

**Fix:**
1. Double-check SID and Token
2. No extra spaces in `.env`
3. Copy full token
4. Regenerate token in Twilio

### "Number not verified"?

**Fix:**
- Verify in Twilio Console, OR
- Upgrade account

---

## 📖 Additional Resources

- **Docs**: https://www.twilio.com/docs/sms
- **Best Practices**: https://www.twilio.com/docs/sms/best-practices
- **Pricing**: https://www.twilio.com/sms/pricing
- **Support**: https://support.twilio.com/

---

## ✅ Testing Checklist

- [ ] Twilio account created
- [ ] Credentials in `.env`
- [ ] Phone number obtained
- [ ] Test number verified
- [ ] SMS received successfully
- [ ] OTP verification works
- [ ] Rate limiting tested
- [ ] Logs reviewed

---

## 🎉 Success!

You're now sending OTPs via SMS! 📱✨

**Next Steps:**
- Test with multiple numbers
- Monitor usage
- Consider upgrading
- Add email OTP backup

---

**Need Help?** Open issue: [GitHub Issues](https://github.com/Avinashkr000/otp-verification-system/issues)
