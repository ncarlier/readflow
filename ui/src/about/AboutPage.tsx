import React from 'react'

import Page from '../common/Page'
import Center from '../common/Center'

import styles from './AboutPage.module.css'

export default () => (
  <Page title="About" >
    <Center className={styles.about}>
      <h1>
        <img src={process.env.PUBLIC_URL + '/logo.svg'} />
      </h1>
      <p>Read your Internet article flow in one place with complete peace of mind and freedom.</p>
      <ul>
        <li>Sources</li>
        <li>Bug or feature request</li>
        <li>Support this project</li>
      </ul>
    </Center>
  </Page>
)
