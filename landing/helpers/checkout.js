
export const createCheckoutSession = async (price, token) => {
  const res = await fetch('/api/create-checkout-session', {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json'
    }),
    credentials: 'same-origin',
    body: JSON.stringify({
      price,
      token
    })
  })

  if (res.error) {
    throw error
  }
  return res.json()
}
