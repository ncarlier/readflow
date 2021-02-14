import styles from './Screenshot.module.css'

const Screenshot = () => (
  <section className={styles.screenshot}>
    <header>
      <h1>Screenshot</h1>
      <hr />
    </header>
    <figure>
      <img src="./img/screenshot.png" alt="screenshot" />
    </figure>
    <header className={styles.try}>
      <h1>Give a try?</h1>
      <h2>our hosted service is waiting for you!</h2>
      <a href="https://readflow.app/login" className="primary btn">Get started</a>
    </header>
  </section>
)

export default Screenshot
