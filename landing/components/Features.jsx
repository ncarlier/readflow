import styles from './Features.module.css'

const Features = () => (
  <section className={styles.features}>
    <header>
      <h1>Features</h1>
      <h2>main features at glance</h2>
      <hr />
    </header>
    <ul className={styles.cards}>
      <li>
        <header>
          <h2>Read everything in one place</h2>
        </header>
        <picture>
          <img src="./img/read.svg" />
        </picture>
        <p>Centralize your news monitoring without frills</p>
        <a href="https://about.readflow.app/docs/en/read-flow/read/" target="_blank" className="btn">Read more</a>
      </li>
      <li>
        <header>
          <h2>Progressive Web App</h2>
        </header>
        <picture>
          <img src="./img/pwa.svg" />
        </picture>
        <p>Great user experiences on all your devices</p>
        <a href="https://about.readflow.app/docs/en/read-flow/mobile/" target="_blank" className="btn">Read more</a>
      </li>
      <li>
        <header>
          <h2>Cloud integration</h2>
        </header>
        <picture>
          <img src="./img/cloud.svg" />
        </picture>
        <p>Easy integration with many other services</p>
        <a href="https://about.readflow.app/docs/en/integrations/" target="_blank" className="btn">Read more</a>
      </li>
    </ul>
  </section>
)

export default Features
