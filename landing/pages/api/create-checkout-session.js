import { base_url } from '@/config/url'
import { pricing } from '@/config/stripe'
import { stripe } from '@/helpers/stripe-server'
import { decodeToken } from '@/helpers/token'
import { getOrRegisterUser, updateUser } from '@/helpers/readflow'

/**
 * Create checkout session.
 * @param {import("next").NextApiRequest} req 
 * @param {import("next").NextApiResponse} res 
 */
const createCheckoutSession = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    return res.status(405).end('Method Not Allowed')
  }
  const { plan, token } = req.body
  try {
    const price = pricing[plan]
    if (!price) {
      throw `price not found the ${plan} plan`
    }
    const { sub, email, preferred_username: username } = decodeToken(token)
    // Retrieve user account
    const user = await getOrRegisterUser(username)
    let customer = user.customer_id
    if (!customer) {
      console.debug('promoting user as a customer...', user)
      // Create stripe customer
      const stripeCustomer = await stripe.customers.create({
        email,
        metadata: {
          subject: sub,
          uid: user.id
        }
      })
      customer = stripeCustomer.id
      // Link user account with customer account
      user = await updateUser(user.id, {customer_id: customer})
      console.info('customer created', user)
    }

    // Create checkout session
    console.debug('creating checkout sessions:', customer)
    const session = await stripe.checkout.sessions.create({
      mode: 'subscription',
      payment_method_types: ['card'],
      billing_address_collection: 'required',
      customer,
      line_items: [
        {
          price,
          quantity: 1
        }
      ],
      success_url: `${base_url}/result?variant=subscription-success&session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${base_url}/pricing`
    })
    return res.status(200).json({ sessionId: session.id })
  } catch (err) {
    console.error('❌ error while creating checkout session', err)
    res.status(500).json({
      error: { statusCode: 500, message: err.message }
    })
  }
}

export default createCheckoutSession
