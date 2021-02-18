import styles from './Icon.module.css'

const Icon = ({name}) => (
  <i className={styles.icon + ' ' + styles[name]} />
)

export default Icon
