import '../styles/globals.css'

import dynamic from 'next/dynamic'

const Auth = dynamic(() => import('../components/Auth'), {
  ssr: false,
})

const App = ({ Component, pageProps }) => (
  <Auth>
    <Component {...pageProps} />
  </Auth>
)

export default App
