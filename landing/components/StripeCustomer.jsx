import { useState } from 'react'

import { postData } from '@/helpers/http'

const StripeCustomer = ({ userData }) => {
  const [loading, setLoading] = useState(false)
  const redirectToCustomerPortal = async () => {
    setLoading(true)
    const { id_token: token } = userData
    const { url, error } = await postData('/api/create-portal-link', { token })
    if (error) return alert(error.message)
    if (url != null) {
      window.location.assign(url)
      setLoading(false)
    }
  }
  return (
    <section>
      <dl>
        <dt>Manage your subscription on Stripe</dt>
        <dd>
          <button disabled={loading} onClick={redirectToCustomerPortal}>
            Open customer portal
          </button>
        </dd>
      </dl>
    </section>
  )
}

export default StripeCustomer
