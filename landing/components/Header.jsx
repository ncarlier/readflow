import Link from 'next/link'
import useTranslation from 'next-translate/useTranslation'

import styles from './Header.module.css'
import Wip from './Wip'
import { useAuth } from 'oidc-react'

const Header = () => {
  const { t } = useTranslation('common')
  const { userData } = useAuth()
  return (
    <header className={styles.header}>
      <nav>
        <Link href="/" passHref><img alt="readflow" src="./img/logo.svg" /></Link>
        <ul>
          <li><Link href="/#features">{t('features')}</Link></li>
          <Wip><li><Link href="/pricing">{t('pricing')}</Link></li></Wip>
          <li><a href="https://docs.readflow.app">{t('docs')}</a></li>
          <li><a href="https://www.github.com/ncarlier/readflow/" target="_blank" rel="noreferrer">{t('sources')}</a></li>
          <li>
            <a href="https://readflow.app/login" title={ userData ? userData.profile.preferred_username : t('login')}>
              { userData ? t('my-readflow') : t('login')}
            </a>
          </li>
        </ul>
      </nav>
    </header>
  )
}

export default Header
