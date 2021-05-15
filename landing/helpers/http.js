
export const postData = async (url, body) => {
  const res = await fetch(url, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json'
    }),
    credentials: 'same-origin',
    body: JSON.stringify(body)
  })

  if (res.error) {
    throw error
  }
  return res.json()
}
