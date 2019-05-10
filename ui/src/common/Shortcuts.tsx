import React from 'react'

import styles from './Shortcuts.module.css'

export default () => (
  <section className={styles.shortcuts}>
    <section>
      <h1>Article list</h1>
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
        <li>
          Next article
          <kbd>right</kbd>
          <kbd>k</kbd>
          <kbd>n</kbd>
        </li>
        <li>
          Previous article
          <kbd>left</kbd>
          <kbd>j</kbd>
          <kbd>p</kbd>
        </li>
        <li>
          Open active article<kbd>enter</kbd>
        </li>
      </ul>
    </section>
    <section>
      <h1>Article</h1>
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
        <li>
          Back to list<kbd>backspace</kbd>
        </li>
      </ul>
    </section>
  </section>
)
