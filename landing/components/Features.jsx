import useTranslation from 'next-translate/useTranslation'

import styles from './Features.module.css'

const Features = () => {
  const { t } = useTranslation('home')
  return (
    <section className={styles.features}>
      <header>
        <h1>{t('common:features')}</h1>
        <h2>{t('at-glance')}</h2>
        <hr />
      </header>
      <ul className={styles.cards}>
        <li>
          <header>
            <h2>{t('feat-1-title')}</h2>
          </header>
          <picture>
            <img src="./img/read.svg" />
          </picture>
          <p>{t('feat-1-desc')}</p>
          <a href="https://about.readflow.app/docs/en/read-flow/read/" target="_blank" className="btn">{t('read-more')}</a>
        </li>
        <li>
          <header>
            <h2>{t('feat-2-title')}</h2>
          </header>
          <picture>
            <img src="./img/pwa.svg" />
          </picture>
          <p>{t('feat-2-desc')}</p>
          <a href="https://about.readflow.app/docs/en/read-flow/mobile/" target="_blank" className="btn">{t('read-more')}</a>
        </li>
        <li>
          <header>
            <h2>{t('feat-3-title')}</h2>
          </header>
          <picture>
            <img src="./img/cloud.svg" />
          </picture>
          <p>{t('feat-3-desc')}</p>
          <a href="https://about.readflow.app/docs/en/integrations/" target="_blank" className="btn">{t('read-more')}</a>
        </li>
      </ul>
    </section>
  )
}

export default Features
