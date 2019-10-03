import React from 'react'

import Center from '../components/Center'
import { VERSION } from '../constants'
import Page from '../layout/Page'
import styles from './AboutPage.module.css'

export default () => (
  <Page title="About">
    <Center className={styles.about}>
      <h1>
        <img src={process.env.PUBLIC_URL + '/logo.svg'} />
      </h1>
      <span>({VERSION})</span>
      <p>Read your Internet article flow in one place with complete peace of mind and freedom.</p>
      <ul>
        <li>
          <a href="https://github.com/ncarlier/readflow" rel="noreferrer noopener" target="_blank">
            Sources
          </a>
        </li>
        <li>
          <a href="https://github.com/ncarlier/readflow/issues" rel="noreferrer noopener" target="_blank">
            Bug or feature request
          </a>
        </li>
        <li>
          <a href="https://www.paypal.me/nunux" rel="noreferrer noopener" target="_blank">
            Support this project
          </a>
        </li>
      </ul>
    </Center>
  </Page>
)
