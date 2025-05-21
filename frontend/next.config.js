/** @type {import('next').NextConfig} */
const nextConfig = {
  // Expose environment variables to the frontend
  env: {
    NEXT_PUBLIC_GRPC_HOST: process.env.NEXT_PUBLIC_GRPC_HOST,
    NEXT_PUBLIC_MAX_RESULTS: process.env.NEXT_PUBLIC_MAX_RESULTS,
    NEXT_PUBLIC_DEFAULT_SUBREDDIT: process.env.NEXT_PUBLIC_DEFAULT_SUBREDDIT,
  },
  
  // Configure webpack to handle gRPC-web properly
  webpack: (config, { isServer }) => {
    // Add the .proto extension to be resolved by webpack
    config.resolve.extensions.push('.proto');
    
    if (!isServer) {
      // This ensures the proper browser polyfills are available for gRPC-web in the browser
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
      };
    }
    
    return config;
  },
  
  // Enable strict mode for React
  reactStrictMode: true,
  
  // Custom async rewrites for development proxy if needed
  async rewrites() {
    return [
      // If your backend requires a proxy for CORS or other reasons
      // You could add it here like this:
      /*
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/:path*', // Proxy to Backend
      }
      */
    ];
  },
};

module.exports = nextConfig;

