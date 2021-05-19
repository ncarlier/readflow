import { useRouter } from 'next/router'
import useTranslation from "next-translate/useTranslation"

import Layout from '@/components/Layout'

const Result = () => {
  const { t } = useTranslation("message")
  const router = useRouter()
  const { variant, reason } = router.query

  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t(`${variant}-title`)}</h1>
            <hr />
          </header>
          <article>
            <p>{t(variant)}</p>
            { reason && <pre>{reason}</pre> }
          </article>
        </section>
      </section>
    </Layout>
  )
}

export default Result
