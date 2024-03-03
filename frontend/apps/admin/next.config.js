const path = require('path');
const withTM = require('next-transpile-modules')([
  '@dinnerdonebetter/models',
  '@dinnerdonebetter/utils',
  '@dinnerdonebetter/api-client',
  '@dinnerdonebetter/logger',
  '@dinnerdonebetter/server-timing',
]);

module.exports = withTM({
  reactStrictMode: true,
  output: 'standalone',
  env: {
    NEXT_PUBLIC_API_ENDPOINT: 'https://api.dinnerdonebetter.dev',
  },
  experimental: {
    outputFileTracingRoot: path.join(__dirname, '../../'),
  },
});
