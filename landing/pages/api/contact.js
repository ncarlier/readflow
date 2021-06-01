import validate from 'deep-email-validator'

import { sendmail_url } from '@/config/url'

/**
 * Post contact form to HTTP endpoint.
 * @param {import("next").NextApiRequest} req 
 * @param {import("next").NextApiResponse} res 
 */
const contactForm = async (req, res) => {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST')
    return res.status(405).end('Method Not Allowed')
  }
  const { subject, from, body, name } = req.body

  if (name) {
    return res.redirect('/result?variant=contact-success&honey')
  }

  const { valid, reason, validators } = await validate({
    email: from,
    validateRegex: true,
    validateMx: true,
    validateTypo: true,
    validateDisposable: true,
    validateSMTP: false,
  })
  if (!valid) {
    const message = validators[reason].reason
    console.error('❌ unable to send contact form: invalid email', message)
    return res.redirect(`/result?variant=contact-error&reason=${message}`)
  }

  const url = new URL(sendmail_url)
  url.searchParams.set('subject', "[readflow-contact] " + subject)
  url.searchParams.set('from', from)
  const resp = await fetch(url, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json'
    }),
    body
  })

  if (resp.error) {
    console.error('❌ unable to send contact form', resp.error)
    return res.redirect(`/result?variant=contact-error&reason=${resp.error}`)
  }

  console.info('contact form sent', from, subject, resp.headers.get('x-hook-id'))

  res.redirect('/result?variant=contact-success')
}

export default contactForm
