import Link from 'next/link'
import Head from 'next/head'

import Header from './Header'
import Footer from './Footer'

const Layout = ({
  children,
  title = 'Welcome to readflow!',
}) => (
  <>
    <Head>
      <title>{title}</title>
      <meta charSet="utf-8" />
      <meta name="viewport" content="initial-scale=1.0, width=device-width" />
      <meta property="og:title" content="readflow" />
      <meta property="og:description" content="Read your Internet article flow in one place with complete peace of mind and freedom" />
      <meta property="og:type" content="website" />
      <meta property="og:url" content="https://about.readflow.app" />
      <meta property="og:image" content="https://about.readflow.app/images/readflow.png" />
      <link rel="icon" type="image/png" href="/favicon.png"></link>
    </Head>
    <Header />
    <main>
      {children}
    </main>
    <Footer />
  </>
)

export default Layout
