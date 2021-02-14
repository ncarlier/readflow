import styles from './Hero.module.css'

const Hero = () => (
  <section className={styles.hero}>
    <div>
      <h1>Read</h1>
      <h2>your Internet article flow in one place with complete peace of mind and freedom</h2>
      <a href="https://readflow.app/login" className="primary btn">Get started</a>
    </div>
    <div>
      <img src="./img/worker.svg" alt="Reading with peaceful mindset" />
    </div>
  </section>
)

export default Hero
