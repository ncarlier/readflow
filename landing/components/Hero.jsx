import useTranslation from 'next-translate/useTranslation'

import styles from './Hero.module.css'

const Hero = () => {
  const { t } = useTranslation('home')
  return (
    <section className={styles.hero}>
      <div>
        <h1>{t('read')}</h1>
        <h2>{t('hero')}</h2>
        <a href="https://readflow.app/login" className="primary btn">{t('get-started')}</a>
      </div>
      <div>
        <picture>
          <source srcset="./img/relax.webp" type="image/webp" />
          <img src="./img/relax.png" alt="Reading with peaceful mindset" />
        </picture>
      </div>
    </section>
  )
}

export default Hero
