const path = require('path');
const withTM = require('next-transpile-modules')();

module.exports = withTM({
  reactStrictMode: true,
  output: 'standalone',
  env: {
    NEXT_PUBLIC_API_ENDPOINT: 'https://api.dinnerdonebetter.dev',
    NEXT_PUBLIC_SEGMENT_API_TOKEN: process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN,
    NEXT_PUBLIC_ENVIRONMENT: process.env.NEXT_PUBLIC_ENVIRONMENT,
  },
  experimental: {
    outputFileTracingRoot: path.join(__dirname, '../../'),
  },
});
