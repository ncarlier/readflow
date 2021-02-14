import styles from './Footer.module.css'

const Footer = () => (
  <footer className={styles.footer}>
    <ul className={styles.columns}>
      <li>
        <img className={styles.logo} alt="readflow" src="./img/logo-white.svg" />
      </li>
      <li>
        <header>
          <h3>product</h3>
        </header>
        <ul>
          <li>
            <a href="https://www.paypal.me/nunux">Support this project</a>
          </li>
          <li>
            <a href="https://github.com/ncarlier/readflow/issue">
              Report a bug
            </a>
          </li>
          <li>
            <a href="https://github.com/ncarlier/readflow/issue">
              Request features
            </a>
          </li>
        </ul>
      </li>
      <li>
        <header>
          <h3>docs</h3>
        </header>
        <ul>
          <li>
            <a href="/docs">Get started</a>
          </li>
          <li>
            <a href="/docs">User guides</a>
          </li>
        </ul>
      </li>
      <li>
        <header>
          <h3>follow us</h3>
        </header>
        <ul>
          <li>
            <a href="https://github.com/ncarlier/readflow">Github</a>
          </li>
          <li>
            <a href="https://twitter.com/ncarlier">Twitter</a>
          </li>
        </ul>
      </li>
    </ul>
  </footer>
)

export default Footer
