const nextTranslate = require('next-translate')

module.exports = nextTranslate({
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
