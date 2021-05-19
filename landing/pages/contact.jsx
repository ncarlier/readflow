import useTranslation from "next-translate/useTranslation"

import Layout from '@/components/Layout'

const Contact = () => {
  const { t } = useTranslation("common")

  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t("contact")}</h1>
            <hr />
          </header>
          <section>
            <form method="post" action="/api/contact">
              <div className="form-group ">
                <label htmlFor="subject">Subject</label>
                <input id="subject" name="subject" type="text" required />
              </div>
              <div className="form-group ">
                <label htmlFor="body">Content</label>
                <textarea id="body" name="body" required></textarea>
              </div>
              <div className="form-group ">
                <label htmlFor="from">Email</label>
                <input id="from" name="from" type="email" required />
              </div>
              <div className="form-group" style={{visibility: 'hidden'}}>
                <label htmlFor="name">Name</label>
                <input id="name" name="name" type="text" />
              </div>
              <div className="actions">
                <button type="submit">Send</button>
              </div>
            </form>
          </section>
        </section>
      </section>
    </Layout>
  )
}

export default Contact
