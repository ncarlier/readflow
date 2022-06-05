import validate from 'deep-email-validator'

import { serviceURL, contactMail } from '@/config/sendmail'

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
    return res.redirect(302, '/result?variant=contact-success&honey')
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
    return res.redirect(302, `/result?variant=contact-error&reason=${message}`)
  }

  const formData = new FormData()
  formData.append('subject', '[readflow-contact] ' + subject)
  formData.append('to', contactMail)
  formData.append('from', from)
  formData.append('text', body)
  const resp = await fetch(serviceURL, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/x-www-form-urlencoded',
    }),
    body: formData
  })

  if (resp.error) {
    console.error('❌ unable to send contact form', resp.error)
    return res.redirect(302, `/result?variant=contact-error&reason=${resp.error}`)
  }

  console.info('contact form sent', from, subject, resp.headers.get('x-hook-id'))

  res.redirect(302, '/result?variant=contact-success')
}

export default contactForm
