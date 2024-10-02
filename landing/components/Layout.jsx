import Head from 'next/head'
import dynamic from 'next/dynamic'

import Header from './Header'
import Footer from './Footer'

const AppAuthProvider = dynamic(() => import('../context/AppAuthProvider'), {
  ssr: false,
})

const Page = ({children}) => (
  <>
    <Header />
    <main>
        {children}
    </main>
    <Footer />
  </>
)

const AutheticatedPage = ({children}) => (
  <AppAuthProvider>
    <Page>{children}</Page>
  </AppAuthProvider>
)

const Layout = ({
  children,
  title = 'Welcome to readflow!',
  desc = 'readflow is a news-reading (or read-it-later) solution focused on versatility and simplicity',
  authenticated = false,
}) => (
  <>
    <Head>
      <title>{title}</title>
      <meta charSet="utf-8" />
      <meta name="description" content={desc} />
      <meta name="viewport" content="initial-scale=1.0, width=device-width" />
      <meta property="og:title" content="readflow" />
      <meta property="og:description" content="read your Internet article flow in one place with complete peace of mind and freedom" />
      <meta property="og:type" content="website" />
      <meta property="og:url" content="https://about.readflow.app" />
      <meta property="og:image" content="https://about.readflow.app/img/readflow.png" />
      <link rel="icon" type="image/png" href="/favicon.png"></link>
    </Head>
    { authenticated ? <AutheticatedPage>{children}</AutheticatedPage> : <Page>{children}</Page>}
  </>
)

export default Layout
