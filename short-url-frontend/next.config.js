const {
  PHASE_DEVELOPMENT_SERVER,
  PHASE_PRODUCTION_SERVER,
} = require('next/constants')

module.exports = (phase) => {
  const isDev = phase === PHASE_DEVELOPMENT_SERVER
  const isProd = phase === PHASE_PRODUCTION_SERVER
  var publicRuntimeConfig = {};

  if (isDev) {
    publicRuntimeConfig = {
      apiURL: process.env.apiURL || 'http://localhost/',
      cloudFunctionUrl: process.env.cloudFunctionUrl || 'http://localhost/'
    }
  }
  if (isProd) {
    publicRuntimeConfig = {
      apiURL: process.env.apiURL,
      cloudFunctionUrl: process.env.cloudFunctionUrl
    }
  }

  return {
    publicRuntimeConfig: publicRuntimeConfig,
  }
}