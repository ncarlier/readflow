import dynamic from 'next/dynamic'
import { useAuth } from 'oidc-react'

import Layout from '@/components/Layout'
import Plans from '@/components/Plans'
import Wip from '@/components/Wip'
import { postData } from '@/helpers/http'
import { getStripe } from '@/helpers/stripe-client'

const PricingSection = () => {
  const auth = useAuth()

  const handlePlanCheckout = async (plan) => {
    console.log(`choosed plan: ${plan}`)
    if (plan === 'free') {
      return document.location = 'https://readflow.app/login'
    }
    if (auth.userData === null) {
      return auth.signIn({redirect_uri: document.location})
    }
    try {
      const { id_token: token } = auth.userData
      const { sessionId } = await postData('/api/create-checkout-session', { plan, token })
      const stripe = await getStripe()
      stripe.redirectToCheckout({ sessionId })
    } catch (error) {
      return alert(error.message)
    }
  }
  return (
    <section>
      <Wip placeholder={<p>WORK IN PROGRESS...</p>}>
        <Plans onChoosePlan={handlePlanCheckout} />
      </Wip>
    </section>
  )
}

const Pricing = () => (
  <Layout authenticated>
    <PricingSection />
  </Layout>
)

export default Pricing
