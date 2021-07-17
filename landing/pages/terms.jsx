import ReactMarkdown from 'react-markdown'
import useTranslation from 'next-translate/useTranslation'

import Layout from '@/components/Layout'

const Terms = ({content}) => {
  const { t } = useTranslation('common')
  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t('terms-and-conditions')}</h1>
            <hr />
          </header>
          <article>
            <ReactMarkdown>
              {content}
            </ReactMarkdown>
          </article>
        </section>
      </section>
    </Layout>
  )
}

export const getStaticProps = async ({ locale }) => {
  const { default: content } = await import(`../policies/terms_${locale}.md`)
  return { props: { content } }
}

export default Terms
