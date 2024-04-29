import useTranslation from 'next-translate/useTranslation'

import styles from './Screenshot.module.css'

const Screenshot = () => {
  const { t } = useTranslation('home')
  return (
    <section className={styles.screenshot}>
      <header>
        <h1>{t('screenshot')}</h1>
        <hr />
      </header>
      <figure>
        <picture>
          <source srcset="./img/screenshot.webp" type="image/webp" />
          <img src="./img/screenshot.png" alt="screenshot" />
        </picture>
      </figure>
      <header className={styles.try}>
        <h1>{t('try')}</h1>
        <h2>{t('hosted-service')}</h2>
        <a href="https://readflow.app/login" className="primary btn">{t('get-started')}</a>
      </header>
    </section>
  )
}

export default Screenshot
