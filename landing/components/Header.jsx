import Link from 'next/link'
import useTranslation from 'next-translate/useTranslation'

import styles from './Header.module.css'

const Header = () => {
  const { t } = useTranslation('common')
  return (
    <header className={styles.header}>
      <nav>
        <Link href="/"><img alt="readflow" src="./img/logo.svg" /></Link>
        <ul>
          <li><Link href="/#features">{t('features')}</Link></li>
          <li><a href="https://docs.readflow.app">{t('docs')}</a></li>
          <li><a href="https://www.github.com/ncarlier/readflow/" target="_blank">{t('sources')}</a></li>
          <li><a href="https://readflow.app/login">{t('login')}</a></li>
        </ul>
      </nav>
    </header>
  )
}

export default Header
