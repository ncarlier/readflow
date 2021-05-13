import { useAuth } from "oidc-react"
import useTranslation from "next-translate/useTranslation"

import Layout from '@/components/Layout'
import UserInfo from '@/components/UserInfo'
import { useEffect } from 'react'

const Account = () => {
  const { t } = useTranslation("common")
  const auth = useAuth()

  return (
    <Layout>
      <section>
        <section>
          <header>
            <h1>{t("account")}</h1>
            <hr />
          </header>
          <section>
            <UserInfo userData={auth.userData} />
          </section>
        </section>
      </section>
    </Layout>
  )
}

export default Account
