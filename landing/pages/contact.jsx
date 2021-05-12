import useTranslation from "next-translate/useTranslation"

import Layout from '../components/Layout'

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
            FORM
          </section>
        </section>
      </section>
    </Layout>
  )
}

export default Contact
