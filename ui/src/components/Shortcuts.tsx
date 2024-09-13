import React from 'react'

import styles from './Shortcuts.module.css'

export const Shortcuts = () => (
  <section className={styles.shortcuts}>
    <section>
      <h1>Article list</h1>
      <ul>
        <li>
          Refresh list<kbd>shift+r</kbd>
        </li>
        <li>
          Toggle sort order<kbd>shift+o</kbd>
        </li>
        <li>
          Toggle sort by (if starred)<kbd>shift+b</kbd>
        </li>
        <li>
          Mark all as read (if category or main list)<kbd>shift+del</kbd>
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
        <li>
          Add new article<kbd>+</kbd>
        </li>
      </ul>
    </section>
    <section>
      <h1>Article</h1>
      <ul>
        <li>
          Mark article as read<kbd>del</kbd>
        </li>
        <li>
          Mark article as unread<kbd>ins</kbd>
        </li>
        <li>
          Read article later<kbd>r</kbd>
        </li>
        <li>
          Star an article<kbd>s</kbd>
        </li>
        <li>
          Edit an article<kbd>e</kbd>
        </li>
        <li>
          Put article offline<kbd>o</kbd>
        </li>
        <li>
          Delete article offline<kbd>d</kbd>
        </li>
        <li>
          Save to default archive service<kbd>shift+s</kbd>
        </li>
        <li>
          Open article menu<kbd>m</kbd>
        </li>
        <li>
          Back to list<kbd>backspace</kbd>
        </li>
      </ul>
    </section>
  </section>
)
