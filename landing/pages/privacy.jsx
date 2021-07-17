import ReactMarkdown from 'react-markdown'
import useTranslation from 'next-translate/useTranslation'

import Layout from '@/components/Layout'

const Privacy = ({content}) => {
  const { t } = useTranslation('common')
  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t('privacy-policy')}</h1>
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
  const { default: content } = await import(`../policies/privacy_${locale}.md`)
  return { props: { content } }
}

export default Privacy
