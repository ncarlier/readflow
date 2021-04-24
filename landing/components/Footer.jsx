import Link from 'next/link'
import useTranslation from 'next-translate/useTranslation'

import styles from './Footer.module.css'
import Icon from './Icon'

const Footer = () => {
  const { t } = useTranslation('common')
  return (
    <footer className={styles.footer}>
      <ul className={styles.columns}>
        <li>
          <img className={styles.logo} alt="readflow" src="./img/logo-white.svg" />
        </li>
        <li>
          <header>
            <h3>{t('product')}</h3>
          </header>
          <ul>
            <li>
              <a href="https://www.paypal.me/nunux">{t('support')}</a>
            </li>
            <li>
              <a href="https://github.com/ncarlier/readflow/issue">
                {t('bug')}
              </a>
            </li>
            <li>
              <a href="https://github.com/ncarlier/readflow/issue">
                {t('rfq')}
              </a>
            </li>
            <li>
              <a href="/contact">
                {t('contact')}
              </a>
            </li>
          </ul>
        </li>
        <li>
          <header>
            <h3>{t('documentation')}</h3>
          </header>
          <ul>
            <li>
              <a href="https://about.readflow.app/docs">{t('user-guide')}</a>
            </li>
            <li>
              <Link href="/terms">{t('terms-and-conditions')}</Link>
            </li>
            <li>
              <Link href="/legal">{t('legal-mentions')}</Link>
            </li>
            <li>
              <Link href="/privacy">{t('privacy-policy')}</Link>
            </li>
          </ul>
        </li>
        <li>
          <header>
            <h3>{t('follow-us')}</h3>
          </header>
          <ul>
            <li>
              <a href="https://github.com/ncarlier/readflow">
                <Icon name="github"/>Github
              </a>
            </li>
            <li>
              <a href="https://twitter.com/ncarlier">
                <Icon name="twitter"/>Twitter
              </a>
            </li>
          </ul>
        </li>
      </ul>
    </footer>
  )
}

export default Footer
