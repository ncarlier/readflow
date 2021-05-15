import { decodeToken } from '@/helpers/token'
import { base_url } from '@/config/url'
import { getUser } from '@/helpers/keycloak'
import { stripe } from '@/helpers/stripe-server'

const createPortalLink = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    return res.status(405).end('Method Not Allowed')
  }
  const { token } = req.body
  try {
    const decoded = decodeToken(token)
    const { sub } = decoded
    let { customer } = decoded
    if (!customer) {
      const user = await getUser(sub)
      if (user.attributes && user.attributes.customer_id) {
        customer = user.attributes.customer_id[0]
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
