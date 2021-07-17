import { useEffect } from 'react'
import { useAuth } from 'oidc-react'

import Layout from '@/components/Layout'
import { postData } from '@/helpers/http'

const PortalRedirection = () => {
  const auth = useAuth()

  useEffect(() => {
    async function getPortalLink() {
      if (auth.userData === null) {
        console.log('user not connected, redirecting to the login portal...')
        return auth.signIn({redirect_uri: document.location})
      }
      const { id_token: token } = auth.userData
      const { url, error } = await postData('/api/create-portal-link', { token })
      if (error) return alert(error.message)
      if (url != null) {
        window.location.assign(url)
      }
    }
    if (!auth.isLoading) {
      getPortalLink()
    }
  }, [ auth ])

  return <p>redirecting to the billing portal...</p>
}

const Account = () => ( 
  <Layout authenticated>
    <section>
      <section>
        <PortalRedirection />
      </section>
    </section>
  </Layout>
)

export default Account
