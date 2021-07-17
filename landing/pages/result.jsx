import { useContext } from 'react'
import { AuthContext } from 'oidc-react'
import { useRouter } from 'next/router'
import useTranslation from 'next-translate/useTranslation'

import Layout from '@/components/Layout'

const BackLink = () => {
  const { t } = useTranslation('common')
  const authContext = useContext(AuthContext)
  if (!authContext || !authContext.userData) {
    return null
  }
  const { userData } = authContext
  const title = userData.profile.preferred_username
  return (
    <a href="https://readflow.app" title={ title }>
      { t('back-to-my-readflow') }
    </a>
  )
}

const Result = () => {
  const { t } = useTranslation('message')
  const router = useRouter()
  const { variant, reason } = router.query
  let className = ''
  if (variant) {
    className = variant.split('-')[1]
  }

  return (
    <Layout authenticated>
      <section>
        <section>
          <header>
            <h1>{t(`${variant}-title`)}</h1>
            <hr />
          </header>
          <article className={className}>
            <p>{t(variant)}</p>
            { reason && <pre>{reason}</pre> }
            <BackLink />
          </article>
        </section>
      </section>
    </Layout>
  )
}

export default Result
