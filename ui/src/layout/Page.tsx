import React, { ReactNode } from 'react'

import { classNames } from '../helpers'
import { usePageTitle } from '../hooks'
import Appbar from './Appbar'
import Content from './Content'
import styles from './Page.module.css'

interface Props {
  title?: string
  subtitle?: string
  className?: string
  children: ReactNode
  header?: ReactNode
  actions?: ReactNode
}

export default (props: Props) => {
  const { title, subtitle, className, children, actions } = props
  let { header = <Appbar title={title} actions={actions} /> } = props

  usePageTitle(title, subtitle)

  return (
    <section className={classNames(styles.page, className)}>
      <header>{header}</header>
      <Content>{children}</Content>
    </section>
  )
}
