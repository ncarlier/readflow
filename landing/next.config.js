const nextTranslate = require('next-translate')

module.exports = nextTranslate({
  target: 'serverless',
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
