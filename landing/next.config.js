const nextTranslate = require('next-translate')

module.exports = nextTranslate({
  target: 'serverless',
  eslint: {
    // Warning: due to Netifly error
    ignoreDuringBuilds: true,
  },
  webpack: (config) => {
    config.module.rules.push(
      {
        test: /\.md$/,
        use: 'raw-loader'
      }
    )
    return config
  },
})
