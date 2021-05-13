import { useAuth } from "oidc-react"

import Layout from '@/components/Layout'
import Plans from '@/components/Plans'
import { createCheckoutSession } from '@/helpers/checkout'
import { getStripe } from '@/helpers/strip'
import { pricing } from '@/config/strip'

const Pricing = () => {
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
      const { id_token } = auth.userData
      const price = pricing[plan]
      if (!price) {
        throw `price not found the ${plan} plan`
      }
      const { sessionId } = await createCheckoutSession(price, id_token)
      const stripe = await getStripe()
      stripe.redirectToCheckout({ sessionId })
    } catch (error) {
      return alert(error.message)
    }
  }
  return (
    <Layout>
      <section>
        <Plans onChoosePlan={handlePlanCheckout} />
      </section>
    </Layout>
  )
}

export default Pricing
