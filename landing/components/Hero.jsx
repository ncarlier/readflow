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
        <img src="./img/relax.png" alt="Reading with peaceful mindset" />
      </div>
    </section>
  )
}

export default Hero
