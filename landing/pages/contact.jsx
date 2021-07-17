import useTranslation from 'next-translate/useTranslation'

import Layout from '@/components/Layout'

const Contact = () => {
  const { t } = useTranslation('common')

  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t('contact')}</h1>
            <hr />
          </header>
          <section>
            <form method="post" action="/api/contact">
              <div className="form-group ">
                <label htmlFor="subject">{t('contact:subject')}</label>
                <input id="subject" name="subject" type="text" required />
              </div>
              <div className="form-group ">
                <label htmlFor="body">{t('contact:content')}</label>
                <textarea id="body" name="body" required></textarea>
              </div>
              <div className="form-group ">
                <label htmlFor="from">{t('contact:email')}</label>
                <input id="from" name="from" type="email" required />
              </div>
              <div className="form-group" style={{visibility: 'hidden'}}>
                <label htmlFor="name">{t('contact:name')}</label>
                <input id="name" name="name" type="text" />
              </div>
              <div className="actions">
                <button type="submit">{t('contact:send')}</button>
              </div>
            </form>
          </section>
        </section>
      </section>
    </Layout>
  )
}

export default Contact
