import { stripe } from '@/helpers/stripe-server'
import { getPlan } from '@/config/stripe'
import { updateUser } from '@/helpers/readflow'

// Stripe requires the raw body to construct the event.
export const config = {
  api: {
    bodyParser: false
  }
}

async function buffer(readable) {
  const chunks = []
  for await (const chunk of readable) {
    chunks.push(
      typeof chunk === "string" ? Buffer.from(chunk) : chunk
    )
  }
  return Buffer.concat(chunks)
}

const relevantEvents = new Set([
  'customer.subscription.created',
  'customer.subscription.updated',
  'customer.subscription.deleted'
])

const webhookHandler = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    return res.status(405).end('Method Not Allowed')
  }
  
  const buf = await buffer(req)
  const sig = req.headers['stripe-signature']
  const webhookSecret =
    process.env.STRIPE_WEBHOOK_SECRET_LIVE ??
    process.env.STRIPE_WEBHOOK_SECRET
  let event

  try {
    event = stripe.webhooks.constructEvent(buf, sig, webhookSecret)
  } catch (err) {
    console.error(`❌ unable to construct Stripe event: ${err.message}`)
    return res.status(400).send(`Webhook Error: ${err.message}`)
  }

  if (relevantEvents.has(event.type)) {
    console.debug(`stripe event received: ${event.type}`)
    try {
      switch (event.type) {
        case 'customer.subscription.created':
        case 'customer.subscription.updated':
        case 'customer.subscription.deleted':
          const subscription = event.data.object
          await manageSubscriptionStatusChange(
            subscription.id,
            subscription.customer
          )
          break
        default:
          throw new Error('Unhandled relevant event!')
      }
    } catch (err) {
      console.error(`❌ error while processing webhook event: ${err.message}`)
      return res.json({ error: 'Webhook handler failed. View logs.' })
    }
  }
  res.json({ received: true })
}

const manageSubscriptionStatusChange = async (
  subscriptionId,
  customerId,
) => {
  const subscription = await stripe.subscriptions.retrieve(subscriptionId)
  // console.debug(subscription)
  const customer = await stripe.customers.retrieve(customerId)
  // console.debug(customer)
  const { uid } = customer.metadata
  if (!uid) {
    throw new Error(`customer ${customerId} not bound with an user account`)
  }
  const { price } = subscription.items.data[0]
  const plan = subscription.status ===  'active' ? getPlan(price.id) : 'default'
  if (plan === null) {
    throw new Error(`no plan defined for this price id: ${proce.id}`)
  }
  await updateUser(uid, { plan })
  // TODO send notification to activate third party services
}

export default webhookHandler
