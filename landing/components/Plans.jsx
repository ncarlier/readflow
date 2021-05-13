import useTranslation from "next-translate/useTranslation"

import styles from "./Plans.module.css"

const Plans = ({ onChoosePlan }) => {
  const { t } = useTranslation("pricing")
  return (
    <section>
      <header>
        <h2>{t("pricing")}</h2>
        <h1>{t("plans")}</h1>
        <hr />
      </header>
      <ul className={styles.plans}>
        <li>
          <section className={styles.plan}>
            <header>
              <h1>Free to play</h1>
              <p>{t("free-desc")}</p>
            </header>
            <ul>
              <li>{t("up-to-art", { count: 200 })}</li>
              <li>{t("up-to-cat", { count: 5 })}</li>
            </ul>
            <footer>
              <h2>
                €0<sub>/{t("month")}</sub>
              </h2>
              <button onClick={() => onChoosePlan('free')}>{t("get-started")}</button>
            </footer>
          </section>
        </li>
        <li>
          <section className={styles.plan}>
            <header>
              <h1>Standard</h1>
              <p>{t("standard-desc")}</p>
            </header>
            <ul>
              <li>{t("up-to-art", { count: 2000 })}</li>
              <li>{t("up-to-cat", { count: 20 })}</li>
            </ul>
            <footer>
              <h2>
                €3<sub>/{t("month")}</sub>
              </h2>
              <button className="primary" onClick={() => onChoosePlan('standard')}>{t("subscribe")}</button>
            </footer>
          </section>
        </li>
        <li>
          <section className={styles.plan}>
            <header>
              <h1>Premium</h1>
              <p>{t("premium-desc")}</p>
            </header>
            <ul>
              <li>{t("up-to-art", { count: 10000 })}</li>
              <li>{t("up-to-cat", { count: 50 })}</li>
              <li>
                RSS feeds<span>with a dedicated Feedpushr instance</span>
              </li>
            </ul>
            <footer>
              <h2>
                €5<sub>/{t("month")}</sub>
              </h2>
              <button className="primary" onClick={() => onChoosePlan('premium')}>{t("subscribe")}</button>
            </footer>
          </section>
        </li>
      </ul>
    </section>
  )
}

export default Plans
