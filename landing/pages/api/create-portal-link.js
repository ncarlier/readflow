import { decodeToken } from '@/helpers/token'
import { base_url } from '@/config/url'
import { getOrRegisterUser } from '@/helpers/readflow'
import { stripe } from '@/helpers/stripe-server'

/**
 * Get Stripe customer portal link.
 * @param {import("next").NextApiRequest} req 
 * @param {import("next").NextApiResponse} res 
 */
const createPortalLink = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    return res.status(405).end('Method Not Allowed')
  }
  const { token } = req.body
  try {
    const decoded = decodeToken(token)
    const { preferred_username: username } = decoded
    let { customer } = decoded
    if (!customer) {
      const user = await getOrRegisterUser(username)
      if (user.customer_id) {
        customer = user.customer_id
      } else {
        return res.status(200).json({url: null})
      }
    }
    const { url } = await stripe.billingPortal.sessions.create({
      customer,
      return_url: `${base_url}/account`
    })
    return res.status(200).json({ url })
  } catch (err) {
    console.log(err);
    res
      .status(500)
      .json({ error: { statusCode: 500, message: err.message } });
  }
}

export default createPortalLink
