import ReactMarkdown from 'react-markdown'
import useTranslation from "next-translate/useTranslation"

import Layout from '@/components/Layout'

const Legal = ({content}) => {
  const { t } = useTranslation("common")
  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t("legal-mentions")}</h1>
            <hr />
          </header>
          <article>
            <ReactMarkdown children={content} />
          </article>
        </section>
      </section>
    </Layout>
  )
}

export async function getStaticProps(context) {
  const { locale } = context
  const content = await import(`../policies/legal_${locale}.md`)
  return {
    props: { content: content.default },
  }
}

export default Legal
