export const pricing = {
  standard: 'price_1IIBJzFmuse653uBwFBBAnHJ',
  premium: 'price_1IIBKUFmuse653uB580Sulj4'
}

export const getPlan = (price) => Object.keys(pricing).find(key => pricing[key] === price)
