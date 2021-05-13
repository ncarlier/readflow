import jwt_decode from "jwt-decode"
import Stripe from 'stripe'

const stripe = new Stripe(
  process.env.STRIPE_SECRET_KEY_LIVE ?? process.env.STRIPE_SECRET_KEY,
  {
    // https://github.com/stripe/stripe-node#configuration
    apiVersion: '2020-08-27',
    // Register this as an official Stripe plugin.
    // https://stripe.com/docs/building-plugins#setappinfo
    appInfo: {
      name: 'Next.js Subscription Starter',
      version: '0.1.0'
    }
  }
)

const base_url = process.env.NODE_ENV === "development" ? "http://localhost:3000" : "https://about.readflow.app"

const createCheckoutSession = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    res.status(405).end('Method Not Allowed')
    return
  }
  const { price, token } = req.body
    try {
      const decoded = jwt_decode(token)
      const {sub, email} = decoded
      // Create strip customer
      const customer = await stripe.customers.create({
        email,
        metadata: {
          subject: sub
        }
      })

      // Create checkout session
      const session = await stripe.checkout.sessions.create({
        mode: 'subscription',
        payment_method_types: ['card'],
        billing_address_collection: 'required',
        customer: customer.id,
        line_items: [
          {
            price,
            quantity: 1
          }
        ],
        success_url: `${base_url}/account?session_id={CHECKOUT_SESSION_ID}`,
        cancel_url: `${base_url}/`
      })
      return res.status(200).json({ sessionId: session.id })
    } catch (err) {
      console.log(err);
      res.status(500).json({
        error: { statusCode: 500, message: err.message }
      })
    }
}

export default createCheckoutSession
