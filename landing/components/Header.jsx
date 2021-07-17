import { useContext } from 'react'
import Link from 'next/link'
import useTranslation from 'next-translate/useTranslation'
import { AuthContext } from 'oidc-react'

import styles from './Header.module.css'
import Wip from './Wip'

const LoginLink = () => {
  const { t } = useTranslation('common')
  const authContext = useContext(AuthContext)
  let title = t('login')
  let text = title
  if (authContext && authContext.userData) {
    const { userData } = authContext
    title = userData.profile.preferred_username
    text = t('my-readflow')
  }
  return (
    <a href="https://readflow.app/login" title={ title }>
      { text }
    </a>
  )
}

const Header = () => {
  const { t } = useTranslation('common')
  return (
      <header className={styles.header}>
        <nav>
          <Link href="/" passHref><img alt="readflow" src="./img/logo.svg" /></Link>
          <ul>
            <li><Link href="/#features">{t('features')}</Link></li>
            <Wip><li><Link href="/pricing">{t('pricing')}</Link></li></Wip>
            <li><a href="https://docs.readflow.app">{t('docs')}</a></li>
            <li><a href="https://www.github.com/ncarlier/readflow/" target="_blank" rel="noreferrer">{t('sources')}</a></li>
            <li><LoginLink /></li>
          </ul>
        </nav>
      </header>
  )
}

export default Header
