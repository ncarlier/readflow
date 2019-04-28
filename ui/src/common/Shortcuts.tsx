import React from 'react'

import styles from './Shortcuts.module.css'

export default () => (
  <section className={styles.shortcuts}>
    <ul>
      <li>Articles list</li>
      <ul>
        <li>
          Refresh list<kbd>shift+r</kbd>
        </li>
        <li>
          Invert list order<kbd>shift+o</kbd>
        </li>
        <li>
          Mark all as read (if category or main list)<kbd>shift+m</kbd>
        </li>
        <li>
          Toggle history list (if category)<kbd>shift+h</kbd>
        </li>
      </ul>
    </ul>
    <ul>
      <li>Article view</li>
      <ul>
        <li>
          Toggle read status<kbd>m</kbd>
        </li>
        <li>
          Put article offline<kbd>o</kbd>
        </li>
        <li>
          Remove article offline<kbd>r</kbd>
        </li>
        <li>
          Save to default archive service<kbd>s</kbd>
        </li>
      </ul>
    </ul>
  </section>
)
