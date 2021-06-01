export const pricing = {
  standard: 'price_1IIBJzFmuse653uBwFBBAnHJ',
  premium: 'price_1IIBKUFmuse653uB580Sulj4'
}

/**
 * Get plan regarding a Stripe pricing ID
 * @param {string} price Stripe pricing ID
 * @returns {string} plan
 */
export const getPlan = (price) => Object.keys(pricing).find(key => pricing[key] === price)
