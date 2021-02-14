import styles from "./Plans.module.css"

const Plans = () => (
  <section>
    <header>
      <h2>pricing</h2>
      <h1>Plans for everybody</h1>
      <hr />
    </header>
    <ul className={styles.plans}>
      <li>
        <section className={styles.plan}>
          <header>
            <h1>Free to play</h1>
            <p>Discover readflow for free.</p>
          </header>
          <ul>
            <li>up to 200 articles</li>
            <li>up to 5 categories</li>
          </ul>
          <footer>
            <h2>
              €0<sub>/month</sub>
            </h2>
            <button>Get started</button>
          </footer>
        </section>
      </li>
      <li>
        <section className={styles.plan}>
          <header>
            <h1>Standard</h1>
            <p>Ideal for common needs.</p>
          </header>
          <ul>
            <li>up to 2000 articles</li>
            <li>up to 20 categories</li>
          </ul>
          <footer>
            <h2>
              €3<sub>/month</sub>
            </h2>
            <button className="primary">Buy now</button>
          </footer>
        </section>
      </li>
      <li>
        <section className={styles.plan}>
          <header>
            <h1>Premium</h1>
            <p>Battery included.</p>
          </header>
          <ul>
            <li>up to 10000 articles</li>
            <li>up to 50 categories</li>
            <li>
              RSS feeds<span>with a dedicated Feedpushr instance</span>
            </li>
          </ul>
          <footer>
            <h2>
              €5<sub>/month</sub>
            </h2>
            <button className="primary">Buy now</button>
          </footer>
        </section>
      </li>
    </ul>
  </section>
)

export default Plans
