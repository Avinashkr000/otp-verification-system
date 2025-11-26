import { useState, useRef, useEffect } from 'react';
import axios from 'axios';

const VerifyOTP = ({ otpData, onSuccess, onBack }) => {
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [resending, setResending] = useState(false);
  const inputRefs = useRef([]);

  useEffect(() => {
    // Focus first input on mount
    inputRefs.current[0]?.focus();
  }, []);

  const handleChange = (index, value) => {
    // Only allow numbers
    if (value && !/^[0-9]$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);
    setError('');

    // Auto-focus next input
    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (index, e) => {
    if (e.key === 'Backspace') {
      if (!otp[index] && index > 0) {
        // Focus previous input on backspace if current is empty
        inputRefs.current[index - 1]?.focus();
      }
    } else if (e.key === 'ArrowLeft' && index > 0) {
      inputRefs.current[index - 1]?.focus();
    } else if (e.key === 'ArrowRight' && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handlePaste = (e) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').slice(0, 6);
    
    if (/^[0-9]+$/.test(pastedData)) {
      const newOtp = pastedData.split('');
      while (newOtp.length < 6) newOtp.push('');
      setOtp(newOtp);
      
      // Focus last filled input or last input
      const lastIndex = Math.min(pastedData.length - 1, 5);
      inputRefs.current[lastIndex]?.focus();
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const otpCode = otp.join('');
    if (otpCode.length !== 6) {
      setError('Please enter all 6 digits');
      return;
    }

    setError('');
    setLoading(true);

    try {
      const response = await axios.post('http://localhost:8080/api/otp/verify', {
        otp_id: otpData.otp_id,
        otp_code: otpCode,
      });

      if (response.data.success) {
        onSuccess(response.data.data);
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Invalid OTP. Please try again.');
      // Clear OTP on error
      setOtp(['', '', '', '', '', '']);
      inputRefs.current[0]?.focus();
    } finally {
      setLoading(false);
    }
  };

  const handleResend = async () => {
    setResending(true);
    setError('');

    try {
      const response = await axios.post('http://localhost:8080/api/otp/resend', {
        otp_id: otpData.otp_id,
      });

      if (response.data.success) {
        // Update OTP data with new OTP ID
        otpData.otp_id = response.data.data.otp_id;
        otpData.expires_at = response.data.data.expires_at;
        
        // Show success message
        alert('OTP resent successfully!');
        
        // Clear current OTP
        setOtp(['', '', '', '', '', '']);
        inputRefs.current[0]?.focus();
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to resend OTP.');
    } finally {
      setResending(false);
    }
  };

  return (
    <div>
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6 transition-colors"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 19l-7-7 7-7"
          />
        </svg>
        Back
      </button>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-4 text-center">
            Enter 6-digit OTP
          </label>
          
          {/* OTP Input Fields */}
          <div className="flex gap-2 justify-center mb-2">
            {otp.map((digit, index) => (
              <input
                key={index}
                ref={(el) => (inputRefs.current[index] = el)}
                type="text"
                inputMode="numeric"
                maxLength={1}
                value={digit}
                onChange={(e) => handleChange(index, e.target.value)}
                onKeyDown={(e) => handleKeyDown(index, e)}
                onPaste={handlePaste}
                className="w-12 h-14 text-center text-2xl font-bold border-2 border-gray-300 rounded-lg focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-200 transition-all"
              />
            ))}
          </div>

          {/* Development OTP Display */}
          {otpData.otp_code && (
            <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-3">
              <p className="text-xs text-yellow-800 text-center">
                <span className="font-semibold">Development Mode:</span> OTP is{' '}
                <span className="font-mono font-bold text-lg">{otpData.otp_code}</span>
              </p>
            </div>
          )}
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            <p className="text-sm">{error}</p>
          </div>
        )}

        <button
          type="submit"
          disabled={loading || otp.join('').length !== 6}
          className="btn btn-primary w-full py-3 text-lg"
        >
          {loading ? (
            <span className="flex items-center justify-center gap-2">
              <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Verifying...
            </span>
          ) : (
            'Verify OTP'
          )}
        </button>

        <div className="text-center">
          <button
            type="button"
            onClick={handleResend}
            disabled={resending}
            className="text-primary-600 hover:text-primary-700 font-medium text-sm"
          >
            {resending ? 'Resending...' : "Didn't receive OTP? Resend"}
          </button>
        </div>
      </form>

      {/* Info */}
      <div className="mt-6 bg-gray-50 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <svg
            className="w-5 h-5 text-gray-500 flex-shrink-0 mt-0.5"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
              clipRule="evenodd"
            />
          </svg>
          <div className="text-sm text-gray-600">
            <p>OTP is valid for 5 minutes</p>
            <p className="mt-1">You have 3 attempts to enter the correct code</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VerifyOTP;
