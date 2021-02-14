import styles from './Header.module.css'

const Header = () => (
  <header className={styles.header}>
    <nav>
      <a href="/"><img alt="readflow" src="./img/logo.svg" /></a>
      <ul>
        <li><a href="/#features">Features</a></li>
        <li><a href="/pricing">Pricing</a></li>
        <li><a href="https://about.readflow.app/docs">Docs</a></li>
        <li><a href="https://www.github.com/ncarlier/readflow/" target="_blank">Sources</a></li>
        <li><a href="https://readflow.app/login">Login</a></li>
      </ul>
    </nav>
  </header>
)

export default Header
