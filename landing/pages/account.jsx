import { useAuth } from "oidc-react"

import Layout from '@/components/Layout'
import { useEffect } from 'react'

const Account = () => {
  const { userData } = useAuth()

  useEffect(async () => {
    if (userData === null) {
      return auth.signIn({redirect_uri: document.location})
    }
    const { id_token: token } = userData
    const { url, error } = await postData('/api/create-portal-link', { token })
    if (error) return alert(error.message)
    if (url != null) {
      window.location.assign(url)
    }
  }, [userData])

  return (
    <Layout>
      <section>
        <section>
          <p>redirecting to the billing portal...</p>
        </section>
      </section>
    </Layout>
  )
}

export default Account
