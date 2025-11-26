import { useState } from 'react';
import RequestOTP from './pages/RequestOTP';
import VerifyOTP from './pages/VerifyOTP';
import Success from './pages/Success';

function App() {
  const [currentStep, setCurrentStep] = useState('request'); // 'request', 'verify', 'success'
  const [otpData, setOtpData] = useState(null);
  const [userData, setUserData] = useState(null);

  const handleOTPGenerated = (data) => {
    setOtpData(data);
    setCurrentStep('verify');
  };

  const handleOTPVerified = (data) => {
    setUserData(data);
    setCurrentStep('success');
  };

  const handleReset = () => {
    setCurrentStep('request');
    setOtpData(null);
    setUserData(null);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 via-blue-50 to-purple-50 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-primary-600 rounded-2xl mb-4">
            <svg
              className="w-8 h-8 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            OTP Verification
          </h1>
          <p className="text-gray-600">
            {currentStep === 'request' && 'Enter your details to receive OTP'}
            {currentStep === 'verify' && 'Enter the OTP sent to you'}
            {currentStep === 'success' && 'Verification successful!'}
          </p>
        </div>

        {/* Main Content */}
        <div className="card">
          {currentStep === 'request' && (
            <RequestOTP onSuccess={handleOTPGenerated} />
          )}
          {currentStep === 'verify' && (
            <VerifyOTP
              otpData={otpData}
              onSuccess={handleOTPVerified}
              onBack={handleReset}
            />
          )}
          {currentStep === 'success' && (
            <Success userData={userData} onReset={handleReset} />
          )}
        </div>

        {/* Footer */}
        <div className="text-center mt-6 text-sm text-gray-600">
          <p>Secure OTP verification powered by Go & React</p>
          <p className="mt-2">
            Built by{' '}
            <a
              href="https://github.com/Avinashkr000"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary-600 hover:text-primary-700 font-medium"
            >
              Avinash Kumar
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;
