const Success = ({ userData, onReset }) => {
  return (
    <div className="text-center">
      {/* Success Icon */}
      <div className="inline-flex items-center justify-center w-20 h-20 bg-green-100 rounded-full mb-6">
        <svg
          className="w-10 h-10 text-green-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M5 13l4 4L19 7"
          />
        </svg>
      </div>

      <h2 className="text-2xl font-bold text-gray-900 mb-2">
        Verification Successful!
      </h2>
      <p className="text-gray-600 mb-8">
        Your {userData.email ? 'email' : 'phone number'} has been verified successfully.
      </p>

      {/* User Info Card */}
      <div className="bg-gradient-to-br from-primary-50 to-blue-50 rounded-xl p-6 mb-8">
        <div className="space-y-3">
          {userData.email && (
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Email</span>
              <span className="font-medium text-gray-900">{userData.email}</span>
            </div>
          )}
          {userData.phone && (
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Phone</span>
              <span className="font-medium text-gray-900">{userData.phone}</span>
            </div>
          )}
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">User ID</span>
            <span className="font-mono text-xs text-gray-900">
              {userData.user_id.slice(0, 8)}...
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Verified At</span>
            <span className="text-xs text-gray-900">
              {new Date(userData.timestamp).toLocaleString()}
            </span>
          </div>
        </div>
      </div>

      {/* Success Features */}
      <div className="grid grid-cols-2 gap-4 mb-8">
        <div className="bg-white border-2 border-gray-100 rounded-lg p-4">
          <div className="text-3xl mb-2">🔒</div>
          <p className="text-sm font-medium text-gray-900">Secure</p>
          <p className="text-xs text-gray-600">End-to-end encryption</p>
        </div>
        <div className="bg-white border-2 border-gray-100 rounded-lg p-4">
          <div className="text-3xl mb-2">⚡</div>
          <p className="text-sm font-medium text-gray-900">Fast</p>
          <p className="text-xs text-gray-600">Instant verification</p>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="space-y-3">
        <button
          onClick={onReset}
          className="btn btn-primary w-full py-3 text-lg"
        >
          Verify Another
        </button>
        
        <a
          href="https://github.com/Avinashkr000/otp-verification-system"
          target="_blank"
          rel="noopener noreferrer"
          className="btn btn-secondary w-full py-3 text-lg inline-block"
        >
          <span className="flex items-center justify-center gap-2">
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            View on GitHub
          </span>
        </a>
      </div>

      {/* Stats */}
      <div className="mt-8 pt-6 border-t border-gray-200">
        <p className="text-xs text-gray-500">
          This verification was completed in less than 30 seconds
        </p>
      </div>
    </div>
  );
};

export default Success;
